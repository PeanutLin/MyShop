package service

import (
	"context"
	"productshop/kitex_gen/shop/validate"
	"productshop/product_shop/common"
	"sync"
	"time"
)

type ValiedateService struct {
}

// 接入的用户的信息
type AccessControl struct {
	// 用户时间
	sourceArray map[int64]time.Time
	// 保证 sourceArray 的并发安全
	sync.RWMutex
}

var (
	// 用户接入控制句柄
	accessControl = &AccessControl{
		sourceArray: make(map[int64]time.Time),
	}
)

func NewValiedateService() *ValiedateService {
	return &ValiedateService{}
}

// 获取接入用户的信息
func (m *AccessControl) GetNewRecord(userID int64) time.Time {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	return m.sourceArray[userID]
}

// 设置接入用户的信息
func (m *AccessControl) SetNewRecord(userID int64) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourceArray[userID] = time.Now()
}

// 权限验证
func (m *AccessControl) CanBuyAgain(userID int64) bool {
	// 一个用户每隔 20s 才能抢购一次
	dataRecord := m.GetNewRecord(userID)
	if !dataRecord.IsZero() {
		if dataRecord.Add(time.Duration(common.Interval) * time.Second).After(time.Now()) {
			return false
		}
	}

	m.SetNewRecord(userID)
	return true
}

func (v *ValiedateService) Validate(ctx context.Context, req *validate.GetValidateReq) (resp *validate.GetValidateResp, err error) {
	// userCookie := req.GetUserCookie()
	// userID := req.GetUserID()

	// right := accessControl.CanBuyAgain(userID)
	// if !right {
	// 	return &validate.GetValidateResp{
	// 		IsSuccess: false,
	// 	}, nil
	// }

	return &validate.GetValidateResp{
		IsSuccess: true,
	}, nil
}
