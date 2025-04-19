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

// 活動資訊快取
type cacheEventMapData struct {
	mu              sync.RWMutex
	Data            map[int]event.PointEvent
	NextRefreshTime time.Time
}
