package common

import (
	"io/ioutil"
	"net/http"
)

// 模拟请求
func GetCurl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	// 获取 Uid
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	// 获取 sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}

	// 模拟接口访问，
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}

	// 手动指定，排查多余cookies
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	// 添加 cookie 到模拟的请求中
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	// 获取返回结果
	response, err = client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	return
}
