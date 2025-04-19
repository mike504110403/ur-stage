package cachedata

import (
	"time"
)

var cfg = Config{
	RefreshDuration: time.Minute * 30,
	RetryDuration:   time.Second * 3,
}

var codeAgentIdMap = cacheAgentIdMapData{
	Data:            map[string]string{},
	NextRefreshTime: time.Now(),
}

var codeAgentSecretMap = cacheAgentSecretMapData{
	Data:            map[string]string{},
	NextRefreshTime: time.Now(),
}

var codeAgentNameIdMap = cacheAgentNameIdMapData{
	Data:            map[string]int{},
	NextRefreshTime: time.Now(),
}

func Init(initCfg Config) {
	cfg = initCfg
	loadAgentData()
	go refreshDataPeriodically()
}

// 定時更新快取資料
func refreshDataPeriodically() {
	for {
		time.Sleep(cfg.RefreshDuration)
		loadAgentData()
	}
}

// 取得遊戲商ID
func AgentIdMap() map[string]string {
	codeAgentIdMap.mu.RLock()
	defer codeAgentIdMap.mu.RUnlock()
	return codeAgentIdMap.Data
}

// 取得遊戲商名稱ID
func AgentNameIdMap() map[string]int {
	codeAgentNameIdMap.mu.RLock()
	defer codeAgentNameIdMap.mu.RUnlock()
	return codeAgentNameIdMap.Data
}

// 取得遊戲商密鑰
func AgentSecretMap() map[string]string {
	codeAgentSecretMap.mu.RLock()
	defer codeAgentSecretMap.mu.RUnlock()
	return codeAgentSecretMap.Data
}
