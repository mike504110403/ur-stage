package cachedata

import (
	"event_service/internal/database"
	"time"

	event "gitlab.com/gogogo2712128/common_moduals/dbModel/Point"

	mlog "github.com/mike504110403/goutils/log"
)

// 活動資訊
func loadEventData() {
	if db, err := database.POINT.DB(); err != nil {
		mlog.Error(err.Error())
	} else {
		queryStr := `
			SELECT 
				id,
				name,
				calc_bonus_rule,
				calc_withdraw_rule,
				start_time,
				end_time,
				is_enable,
				create_at,
				update_at
			FROM PointEvent
		`
		if rows, err := db.Query(queryStr); err != nil {
			mlog.Error(err.Error())
		} else {
			eventMap.mu.Lock()
			defer eventMap.mu.Unlock()
			eventMap.Data = make(map[int]event.PointEvent)
			for rows.Next() {
				var event event.PointEvent
				if err := rows.Scan(
					&event.Id,
					&event.Name,
					&event.BonusRuleSp,
					&event.WithdrawRuleSp,
					&event.StartTime,
					&event.EndTime,
					&event.IsEnable,
					&event.CreateAt,
					&event.UpdateAt,
				); err != nil {
					mlog.Error(err.Error())
				} else {
					eventMap.Data[event.Id] = event
				}
			}
			eventMap.NextRefreshTime = time.Now().Add(cfg.RefreshDuration)
		}
	}
}
