package cachedata

import (
	"sync"
	"time"
)

type Config struct {
	RefreshDuration time.Duration
	RetryDuration   time.Duration
}

// 遊戲商資訊
type cacheAgentIdMapData struct {
	mu              sync.RWMutex
	Data            map[string]string
	NextRefreshTime time.Time
}

type cacheAgentSecretMapData struct {
	mu              sync.RWMutex
	Data            map[string]string
	NextRefreshTime time.Time
}

type cacheAgentNameIdMapData struct {
	mu              sync.RWMutex
	Data            map[string]int
	NextRefreshTime time.Time
}
