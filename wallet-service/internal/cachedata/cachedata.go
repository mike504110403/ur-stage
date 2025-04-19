package cachedata

import (
	"time"

	event "gitlab.com/gogogo2712128/common_moduals/dbModel/Point"
)

var cfg = Config{
	RefreshDuration: time.Minute * 3,
	RetryDuration:   time.Second * 3,
}

var codeThirdSecretMap = cacheThirdPaySecretMapData{
	Data:            map[string]string{},
	NextRefreshTime: time.Now(),
}

var codeThirdIdMap = cacheThirdPayIdMapData{
	Data:            map[string]string{},
	NextRefreshTime: time.Now(),
}

var codeItemIdMap = cacheItemIdMapData{
	Data:            map[string]string{},
	NextRefreshTime: time.Now(),
}

var codeItemValMap = cacheItemValMapData{
	Data:            map[string]string{},
	NextRefreshTime: time.Now(),
}

var codeVipInfoMap = cacheVipInfoMapData{
	Data:            map[int]VipInfo{},
	NextRefreshTime: time.Now(),
}

var eventMap = cacheEventMapData{
	Data:            map[int]event.PointEvent{},
	NextRefreshTime: time.Now(),
}

var eventNameMap = cacheEventNameMapData{
	Data:            map[string]event.PointEvent{},
	NextRefreshTime: time.Now(),
}

func Init(initCfg Config) {
	cfg = initCfg
	loadThirdPayData()
	loadItemData()
	loadVipData()
	loadEventData()
	go refreshDataPeriodically()
}

// 定時更新快取資料
func refreshDataPeriodically() {
	for {
		time.Sleep(cfg.RefreshDuration)
		loadThirdPayData()
		loadItemData()
		loadVipData()
		loadEventData()
	}
}

// ThirdPayMap 返回第三方支付信息
func ThirdPaySecretMap() map[string]string {
	codeThirdSecretMap.mu.RLock()
	defer codeThirdSecretMap.mu.RUnlock()
	return codeThirdSecretMap.Data
}

// ThirdPayIdMap 返回第三方支付信息
func ThirdPayIdMap() map[string]string {
	codeThirdIdMap.mu.RLock()
	defer codeThirdIdMap.mu.RUnlock()
	return codeThirdIdMap.Data
}

// ItemIdMap 返回商品信息
func ItemIdMap() map[string]string {
	codeItemIdMap.mu.RLock()
	defer codeItemIdMap.mu.RUnlock()
	return codeItemIdMap.Data
}

// ItemValMap 返回商品信息
func ItemValMap() map[string]string {
	codeItemValMap.mu.RLock()
	defer codeItemValMap.mu.RUnlock()
	return codeItemValMap.Data
}

// VipInfoMap 返回 vip 信息
func VipInfoMap() map[int]VipInfo {
	codeVipInfoMap.mu.RLock()
	defer codeVipInfoMap.mu.RUnlock()
	return codeVipInfoMap.Data
}

// EventMap : 活動資訊快取
func EventMap() map[int]event.PointEvent {
	eventMap.mu.RLock()
	defer eventMap.mu.RUnlock()
	return eventMap.Data
}

// EventNameMap : 活動資訊快取
func EventNameMap() map[string]event.PointEvent {
	eventNameMap.mu.RLock()
	defer eventNameMap.mu.RUnlock()
	return eventNameMap.Data
}
