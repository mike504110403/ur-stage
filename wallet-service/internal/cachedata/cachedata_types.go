package cachedata

import (
	"sync"
	"time"

	event "gitlab.com/gogogo2712128/common_moduals/dbModel/Point"
)

type Config struct {
	RefreshDuration time.Duration
	RetryDuration   time.Duration
}

// 三方支付資訊
type cacheThirdPaySecretMapData struct {
	mu              sync.RWMutex
	Data            map[string]string
	NextRefreshTime time.Time
}

type cacheThirdPayIdMapData struct {
	mu              sync.RWMutex
	Data            map[string]string
	NextRefreshTime time.Time
}

// 商品資訊
type cacheItemIdMapData struct {
	mu              sync.RWMutex
	Data            map[string]string
	NextRefreshTime time.Time
}

type cacheItemValMapData struct {
	mu              sync.RWMutex
	Data            map[string]string
	NextRefreshTime time.Time
}

// vip 資訊
type VipInfo struct {
	Threshold       int     // vip 流水門檻
	WithdrawTimes   int     // vip 免手續費提款次數
	WithdrawLimit   float64 // vip 每日提款上限
	WithdrawFeeRate float64 // vip 提款手續費比率
	RebateRate      string  // vip 返水比率
}

type cacheVipInfoMapData struct {
	mu              sync.RWMutex
	Data            map[int]VipInfo
	NextRefreshTime time.Time
}

// 活動資訊快取
type cacheEventMapData struct {
	mu              sync.RWMutex
	Data            map[int]event.PointEvent
	NextRefreshTime time.Time
}

type cacheEventNameMapData struct {
	mu              sync.RWMutex
	Data            map[string]event.PointEvent
	NextRefreshTime time.Time
}
