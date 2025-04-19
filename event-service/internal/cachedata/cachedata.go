package cachedata

import (
	"time"

	event "gitlab.com/gogogo2712128/common_moduals/dbModel/Point"
)

var cfg = Config{
	RefreshDuration: time.Minute * 3,
	RetryDuration:   time.Second * 3,
}

var eventMap = cacheEventMapData{
	Data:            map[int]event.PointEvent{},
	NextRefreshTime: time.Now(),
}

func Init(initCfg Config) {
	cfg = initCfg
	loadEventData()
	go refreshDataPeriodically()
}

// 定時更新快取資料
func refreshDataPeriodically() {
	for {
		time.Sleep(cfg.RefreshDuration)
		loadEventData()
	}
}

// EventMap : 活動資訊快取
func EventMap() map[int]event.PointEvent {
	eventMap.mu.RLock()
	defer eventMap.mu.RUnlock()
	return eventMap.Data
}
