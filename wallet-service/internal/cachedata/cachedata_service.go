package cachedata

import (
	"time"
	"wallet_service/internal/database"

	event "gitlab.com/gogogo2712128/common_moduals/dbModel/Point"

	mlog "github.com/mike504110403/goutils/log"
)

// 三方支付資訊
func loadThirdPayData() {
	if db, err := database.ORDER.DB(); err != nil {
		mlog.Error(err.Error())
	} else {
		queryStr := `
			SELECT id, client_name, secret_info 
			FROM ThirdPay
			WHERE is_enable = 1
		`
		if rows, err := db.Query(queryStr); err != nil {
			mlog.Error(err.Error())
		} else {
			codeThirdSecretMap.mu.Lock()
			codeThirdIdMap.mu.Lock()
			defer codeThirdIdMap.mu.Unlock()
			defer codeThirdSecretMap.mu.Unlock()
			codeThirdSecretMap.Data = make(map[string]string)
			codeThirdIdMap.Data = make(map[string]string)
			for rows.Next() {
				var agentId string
				var agentName string
				var secretInfo *string
				if err := rows.Scan(&agentId, &agentName, &secretInfo); err != nil {
					mlog.Error(err.Error())
				} else {
					codeThirdIdMap.Data[agentName] = agentId
					if secretInfo == nil {
						continue
					}
					codeThirdSecretMap.Data[agentName] = *secretInfo
				}
			}
			codeThirdSecretMap.NextRefreshTime = time.Now().Add(cfg.RefreshDuration)
		}
	}
}

// 商品資訊
func loadItemData() {
	if db, err := database.ORDER.DB(); err != nil {
		mlog.Error(err.Error())
	} else {
		queryStr := `
			SELECT id, item_name, item_value
			FROM Item
			WHERE status = 1
		`
		if rows, err := db.Query(queryStr); err != nil {
			mlog.Error(err.Error())
		} else {
			codeItemIdMap.mu.Lock()
			codeItemValMap.mu.Lock()
			defer codeItemIdMap.mu.Unlock()
			defer codeItemValMap.mu.Unlock()
			codeItemIdMap.Data = make(map[string]string)
			codeItemValMap.Data = make(map[string]string)
			for rows.Next() {
				var itemId string
				var itemName string
				var itemVal string
				if err := rows.Scan(&itemId, &itemName, &itemVal); err != nil {
					mlog.Error(err.Error())
				} else {
					codeItemIdMap.Data[itemName] = itemId
					codeItemValMap.Data[itemName] = itemVal
				}
			}
			codeItemIdMap.NextRefreshTime = time.Now().Add(cfg.RefreshDuration)
			codeItemValMap.NextRefreshTime = time.Now().Add(cfg.RefreshDuration)
		}
	}
}

// vip 資訊
func loadVipData() {
	if db, err := database.SETTING.DB(); err != nil {
		mlog.Error(err.Error())
	} else {
		queryStr := `
			SELECT 
				vip_level,
				vip_threshold,
				vip_withdraw_times,
				vip_withdraw_limit,
				vip_withdraw_fee_rate,
				vip_rebate_rate
			FROM VIP_Details
		`
		if rows, err := db.Query(queryStr); err != nil {
			mlog.Error(err.Error())
		} else {
			codeVipInfoMap.mu.Lock()
			defer codeVipInfoMap.mu.Unlock()
			codeVipInfoMap.Data = make(map[int]VipInfo)
			for rows.Next() {
				var vipLevel int
				var threshold int
				var withdrawTimes int
				var withdrawLimit float64
				var withdrawFeeRate float64
				var rebateRate string
				if err := rows.Scan(
					&vipLevel,
					&threshold,
					&withdrawTimes,
					&withdrawLimit,
					&withdrawFeeRate,
					&rebateRate,
				); err != nil {
					mlog.Error(err.Error())
				} else {
					codeVipInfoMap.Data[vipLevel] = VipInfo{
						Threshold:       threshold,
						WithdrawTimes:   withdrawTimes,
						WithdrawLimit:   withdrawLimit,
						WithdrawFeeRate: withdrawFeeRate,
						RebateRate:      rebateRate,
					}
				}
			}
			codeVipInfoMap.NextRefreshTime = time.Now().Add(cfg.RefreshDuration)
		}
	}
}

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
			eventNameMap.mu.Lock()
			defer eventNameMap.mu.Unlock()
			eventMap.Data = make(map[int]event.PointEvent)
			eventNameMap.Data = make(map[string]event.PointEvent)
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
					eventNameMap.Data[event.Name] = event
				}
			}
			eventMap.NextRefreshTime = time.Now().Add(cfg.RefreshDuration)
			eventNameMap.NextRefreshTime = time.Now().Add(cfg.RefreshDuration)
		}
	}
}
