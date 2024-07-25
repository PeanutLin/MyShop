package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"productshop/common"
	"productshop/datamodels"
	"productshop/rabbitmq"
	"strconv"
	"sync"
	"time"
)

// 接入的用户的信息
type AccessControl struct {
	// 用户时间
	sourceArray map[int]time.Time
	// 保证 sourceArray 的并发安全
	sync.RWMutex
}

var (
	// 本机地址
	localHost = ""

	// 一致性哈希句柄
	hashConsistent *common.Consistent
	
	// 用户接入控制句柄
	accessControl = &AccessControl{
		sourceArray: make(map[int]time.Time),
	}

	// rabbitmq 句柄
	rabbitMqValidate *rabbitmq.RabbitMQ
)

// 获取接入用户的信息
func (m * AccessControl) GetNewRecord(uid int) time.Time {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	return m.sourceArray[uid]
}

// 设置接入用户的信息
func (m * AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourceArray[uid] = time.Now()
}

// 分布式验证
func (m * AccessControl) GetDistributedRight(r *http.Request) bool {
	// 获取用用户 uid
	uid, err := r.Cookie("uid")
	if err != nil {
		return false
	}

	// 采用一致性 hash 算法，根据用户 ID，判断具体机器
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	// 是否本机
	if hostRequest == localHost {
		// 执行本机数据读取和校验
		return m.GetDataFromMap(uid.Value)
	} else {
		// 代理访问结果
		return m.GetDataFromOtherMap(uid.Value, r)
	}
}

// 获取用户数据（实际作用）
func (m *AccessControl) GetDataFromMap(uid string) bool {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}

	// 一个用户每隔 20s 才能抢购一次
	dataRecord := m.GetNewRecord(uidInt)
	if !dataRecord.IsZero() {
		if dataRecord.Add(time.Duration(common.Interval) * time.Second).After(time.Now()) {
			return false
		}
	}

	m.SetNewRecord(uidInt)
	return true
}

// 获取用户数据（代理获取）
func (m *AccessControl) GetDataFromOtherMap(host string, request *http.Request) bool {
	hostUrl := "http://" + host + ":" + common.ValidatePort + "/checkRight" 
	response, body,err := common.GetCurl(hostUrl, request)
	if err != nil {
		fmt.Println("get curl error")
		return false
	}
	
	// 判断状态
	if response.StatusCode == 200 {
		return string(body) == "true"
	}
	return false
}

// 统一验证拦截器，每个接口都需要提前验证
func Auth(res http.ResponseWriter, req *http.Request) error {
	return CheckUserInfo(req)
}

func CheckUserInfo(req *http.Request) error {
	uidCookie, err := req.Cookie("uid")
	if err != nil {
		return errors.New("user is not logged in")
	}
	signCookie, err := req.Cookie("sign")
	if err != nil {
		return errors.New("failed to obtain the user encryption string")
	}
	deSignByte, err := common.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("the encryption string has been tampered with")
	}
	if CheckIdInfo(uidCookie.Value, string(deSignByte)) {
		return nil
	}
	return errors.New("identity verification failed")
}

// 自定义逻辑判断
func CheckIdInfo(checkStr string, signStr string) bool {
	return checkStr == signStr
}

// 生成访问数量控制接口的 URL
func generateProductURL(productID int64, productNum int64) string {
	tail := "?productName=product-" + strconv.FormatInt(productID, 10) + "&productNum=" + strconv.FormatInt(productNum, 10)
	result := "http://" + common.GetProductIP + ":" + common.GetProductPort + "/getProduct" +  tail
	fmt.Println(result)
	return result
}

func BuyProduct(userID int64, productID int64, r *http.Request) []byte {
	// 请求 getProduct 服务器
	hostUrl := generateProductURL(productID, 1)
	fmt.Printf("one product : %s\n", hostUrl)
	responseValidate, validateBody, err := common.GetCurl(hostUrl, r)
	if err !=nil {
		return []byte("getProduct false")
	}
	
	// 判断数量控制接口请求状态
	if responseValidate.StatusCode == 200 {
		if string(validateBody)=="true" {

			// 1.创建消息体
			message := datamodels.NewMessage(userID, productID, 1)
			// 类型转化
			byteMessage, err :=json.Marshal(message)
			if err !=nil {
				return[]byte("json Marshal false")
			}

			// 2.生产消息
			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err !=nil {
				return []byte("PublishSimple false")
				
			}
			return []byte("true")
		}
	}
	return []byte("not 200 false")
}

// 抢购处理
// http://localhost:8083/onsale?productID=1 cookie
func OnSaleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("running OnSale")
	// 获得 productID
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err !=nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("query false"))
		return
	}

	productString := queryForm["productID"][0]

	// 获取用户 cookie
	userCookie, err := r.Cookie("uid")
	if err !=nil {
		w.Write([]byte("cookie false"))
		return
	}

	// 获取商品ID
	productID,err :=strconv.ParseInt(productString,10,64)
	if err !=nil {
		w.Write([]byte("Parse productID false"))
		return
	}

	// 获取用户ID
	userID,err := strconv.ParseInt(userCookie.Value, 10, 64)
	if err !=nil {
		w.Write([]byte("Parse userID false"))
		return
	}

	// 1.分布式权限验证
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("GetDistributedRight false"))
		return
	}

	// 2.获取数量控制权限，防止秒杀出现超卖现象
	result := BuyProduct(userID, productID, r)
	w.Write(result)
}

// 权限检查处理
func CheckRightHandler(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
}


// 启动 HTTP 服务器
func StartHTTPServer() {
	// 1. 过滤器
	filter := common.NewFilter()
	// 注册拦截器
	filter.RegisterUriFilter("/onsale", Auth)
	filter.RegisterUriFilter("/checkRight", Auth)

	// 2,启动服务
	// 用于验证和访问数量控制服务
	http.HandleFunc("/onsale", filter.Handler(OnSaleHandler))
	// 用于分布式验证
	http.HandleFunc("/checkRight", filter.Handler(CheckRightHandler))
	http.ListenAndServe(common.ValidateHost1 + ":" + common.ValidatePort, nil)
}

func main() {
	// 负载均衡器设置 采用一致性哈希算法
	hashConsistent = common.NewConsistent()
	for _, v := range common.ClusterHostArray {
		hashConsistent.Add(v)
	}

	// 获取本地 IP
	localIp, err := common.GetIntranceIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Printf("localHost: %s\n", localHost)

	// rabbitmq
	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("productshop")
	defer rabbitMqValidate.Destroy()

	StartHTTPServer()
}