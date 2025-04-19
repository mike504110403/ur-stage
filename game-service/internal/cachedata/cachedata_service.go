package cachedata

import (
	"game_service/internal/database"
	"strconv"
	"time"

	mlog "github.com/mike504110403/goutils/log"
)

// 商品資訊
func loadAgentData() {
	if db, err := database.GAME.DB(); err != nil {
		mlog.Error(err.Error())
	} else {
		queryStr := `
			SELECT id, agent_name, secrect_info
			FROM GameAgent
			WHERE is_enable = 1;
		`
		if rows, err := db.Query(queryStr); err != nil {
			mlog.Error(err.Error())
		} else {
			codeAgentIdMap.mu.Lock()
			codeAgentSecretMap.mu.Lock()
			codeAgentNameIdMap.mu.Lock()
			defer codeAgentIdMap.mu.Unlock()
			defer codeAgentSecretMap.mu.Unlock()
			defer codeAgentNameIdMap.mu.Unlock()
			codeAgentIdMap.Data = make(map[string]string)
			codeAgentSecretMap.Data = make(map[string]string)
			codeAgentNameIdMap.Data = make(map[string]int)
			for rows.Next() {
				var agentId int
				var agentName string
				var agentVal *string
				if err := rows.Scan(&agentId, &agentName, &agentVal); err != nil {
					mlog.Error(err.Error())
				} else {
					codeAgentIdMap.Data[agentName] = strconv.Itoa(agentId)
					codeAgentNameIdMap.Data[agentName] = agentId
					if agentVal == nil {
						continue
					}
					codeAgentSecretMap.Data[agentName] = *agentVal
				}
			}
			codeAgentIdMap.NextRefreshTime = time.Now().Add(cfg.RefreshDuration)
			codeAgentSecretMap.NextRefreshTime = time.Now().Add(cfg.RefreshDuration)
			codeAgentNameIdMap.NextRefreshTime = time.Now().Add(cfg.RefreshDuration)
		}
	}
}
