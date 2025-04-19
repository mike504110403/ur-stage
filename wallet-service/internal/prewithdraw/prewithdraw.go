package prewithdraw

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"wallet_service/internal/cachedata"
	"wallet_service/internal/database"

	point "gitlab.com/gogogo2712128/common_moduals/dbModel/Point"
	wallet "gitlab.com/gogogo2712128/common_moduals/dbModel/Wallet"
)

const NORMAL_POINT_EVENT_NAME = "normal-point"
const ADMIN_POINT_EVENT_NAME = "admin-point"

// InsertPointMemberEvent 新增提款門檻
func InsertPointMemberEvent(tx *sql.Tx, transSet wallet.TransectionSet) error {
	event_id, err := relateEventId(transSet)
	if err != nil {
		return err
	}
	eventMap := cachedata.EventMap()
	var gate float64
	event, ok := eventMap[event_id]
	if !ok {
		return errors.New("提款事件不存在")
	} else {
		if event.WithdrawRuleSp != nil {
			if g, err := getRuleGateBySp(*event.WithdrawRuleSp, transSet.GateAmount); err != nil {
				return err
			} else {
				gate = g
			}
		} else {
			gate = transSet.GateAmount
		}
	}
	if err := InsertMemberEvents(tx, point.MemberEvents{
		MemberId: transSet.MemberId,
		EventId:  event_id,
		Amount:   transSet.GateAmount,
		Gate:     gate,
		GateFlow: -gate,
	}); err != nil {
		return err
	}
	return nil
}

// 透過 SP 取得提款門檻
func getRuleGateBySp(sp string, amount float64) (float64, error) {
	var gate float64
	db, err := database.POINT.DB()
	if err != nil {
		return gate, err
	}
	queryStr := fmt.Sprintf("CALL %s(?, @p_gate)", sp)
	if _, err = db.Exec(queryStr, amount); err != nil {
		return gate, err
	}
	if err = db.QueryRow("SELECT @p_gate").Scan(&gate); err != nil {
		return gate, err
	}
	return gate, nil
}

// UpdateWithdrawState 更新提款狀態
func UpdateWithdrawState(tx *sql.Tx, mid int, amount float64) error {
	upsertStr := `
		INSERT INTO WithdrawGate (member_id, thisday_withdraw_times, thisday_withdraw_amount)
		VALUES (?, 1, ?)
		ON DUPLICATE KEY UPDATE 
			thisday_withdraw_times = thisday_withdraw_times + 1,
			thisday_withdraw_amount = thisday_withdraw_amount + ?
	`
	if _, err := tx.Exec(upsertStr, mid, amount, amount); err != nil {
		return err
	} else {
		return nil
	}
}

// WithdrawableAmount 可提款金額
func WithdrawableAmount(mid int) (float64, error) {
	// 取得vip等級資訊
	cacheVipMap, err := GetVipInfoMap(mid)
	if err != nil {
		return 0, err
	}
	// 取得提款門檻相關資訊
	gateInfo, err := GetEventGate(mid)
	if err != nil {
		return 0, err
	} else {
		var gateAmount float64
		for _, event := range gateInfo.Events {
			gateAmount += event.Amount
		}
		// 提款流水門檻
		balance := gateInfo.Amount - gateAmount
		thisdayWithdrawtimes := func() int {
			if gateInfo.TihsdayWithdrawTimes == nil {
				return 0
			} else {
				return *gateInfo.TihsdayWithdrawTimes
			}
		}()

		thisdaywithdrawamount := func() float64 {
			if gateInfo.ThisdayWithdrawAmount == nil {
				return 0
			} else {
				return *gateInfo.ThisdayWithdrawAmount
			}
		}()
		if balance < 0 || (cacheVipMap.WithdrawTimes-thisdayWithdrawtimes) < 1 {
			return 0, nil
		} else {
			// 實際可提款金額 - 以現有餘額計算若要手需費 實際上可提多少
			realWithdrawAmount := calRealWithdrawAmount(cacheVipMap, balance)
			// 提款額度限額 - 已提款金額 - 實際提款金額
			if cacheVipMap.WithdrawLimit == -1 {
				return realWithdrawAmount, nil
			} else if cacheVipMap.WithdrawLimit-thisdaywithdrawamount <= 0 {
				return 0, errors.New("超出每日提款限額")
			} else if balance-(cacheVipMap.WithdrawLimit-thisdaywithdrawamount) > 0 {
				return cacheVipMap.WithdrawLimit - thisdaywithdrawamount, nil
			} else {
				return realWithdrawAmount, nil
			}
		}
	}
}

// IsWithdrawable 是否可提款
func IsWithdrawable(mid int, amount float64) bool {
	withdrawable, err := WithdrawableAmount(mid)
	if err != nil {
		return false
	} else {
		return withdrawable > amount
	}
}

// 計算實際提款金額
func calRealWithdrawAmount(cacheVipMap cachedata.VipInfo, amount float64) float64 {
	return roundFloat(amount / (1 + cacheVipMap.WithdrawFeeRate/100))
}

// 無條件捨去到小數第二位
func roundFloat(value float64) float64 {
	return math.Floor(value*100) / 100
}

// 無條件進位到小數點後第二位
func ceilFloat(value float64) float64 {
	return math.Ceil(value*100) / 100
}
