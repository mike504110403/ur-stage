package prewithdraw

import (
	"database/sql"
	"wallet_service/internal/cachedata"
	"wallet_service/internal/database"

	event "gitlab.com/gogogo2712128/common_moduals/dbModel/Point"
	wallet "gitlab.com/gogogo2712128/common_moduals/dbModel/Wallet"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

// 取得提款門檻相關資訊
func GetEventGate(mid int) (EventGateDto, error) {
	gateInfo := EventGateDto{}
	walletDb, err := database.WALLET.DB()
	if err != nil {
		return EventGateDto{}, err
	}
	queryStr := `
		SELECT
			w.member_id,
			w.balance,
			wg.thisday_withdraw_times,
			wg.thisday_withdraw_amount
		FROM Wallet AS w
		LEFT JOIN WithdrawGate wg ON w.member_id = wg.member_id
		WHERE w.member_id = ?
	`
	if err := walletDb.QueryRow(queryStr, mid).Scan(
		&gateInfo.MemberId,
		&gateInfo.Amount,
		&gateInfo.TihsdayWithdrawTimes,
		&gateInfo.ThisdayWithdrawAmount,
	); err != nil {
		return EventGateDto{}, err
	}
	if events, err := GetMemberEvents(mid); err != nil {
		return EventGateDto{}, err
	} else {
		gateInfo.Events = events
	}

	return gateInfo, nil
}

// 取得會員點數限制事件
func GetMemberEvents(mid int) ([]event.MemberEvents, error) {
	events := []event.MemberEvents{}
	eventDb, err := database.POINT.DB()
	if err != nil {
		return events, err
	}
	queryStr := `
		SELECT
			id,
			member_id,
			event_id,
			amount,
			gate,
			gate_flow,
			is_release,
			create_at,
			update_at
		FROM MemberEvents
		WHERE
			member_id = ? AND is_release = 0 AND gate_flow < 0
		ORDER BY create_at ASC
	`
	rows, err := eventDb.Query(queryStr, mid)
	if err != nil {
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		event := event.MemberEvents{}
		if err := rows.Scan(
			&event.Id,
			&event.MemberId,
			&event.EventId,
			&event.Amount,
			&event.Gate,
			&event.GateFlow,
			&event.IsRelease,
			&event.CreateAt,
			&event.UpdateAt,
		); err != nil {
			return events, err
		}
		events = append(events, event)
	}
	return events, nil
}

// 新增提款事件
func InsertMemberEvents(tx *sql.Tx, event event.MemberEvents) error {
	insertStr := `
		INSERT INTO MemberEvents (member_id, event_id, amount, gate, gate_flow)
		VALUES (?, ?, ?, ?, ?)
	`
	if _, err := tx.Exec(
		insertStr,
		event.MemberId,
		event.EventId,
		event.Amount,
		event.Gate,
		event.GateFlow,
	); err != nil {
		return err
	} else {
		return nil
	}
}

// 計算提款手續費
func GetWithdrawFee(mid int, amount float64) (float64, error) {
	// 取得vip等級資訊
	cacheVipMap, err := GetVipInfoMap(mid)
	if err != nil {
		return 0, err
	}
	// 計算提款手續費
	return ceilFloat(amount * cacheVipMap.WithdrawFeeRate / 100), nil
}

// 取得vip等級資訊
func GetVipInfoMap(mid int) (cachedata.VipInfo, error) {
	vipInfo := cachedata.VipInfo{}
	// 會員等級
	memberDb, err := database.MEMBER.DB()
	if err != nil {
		return vipInfo, err
	}
	queryStr := `
			SELECT vip_level
			FROM MemberLevel
			WHERE member_id = ?
		`
	var vipLevel int
	if err := memberDb.QueryRow(queryStr, mid).Scan(&vipLevel); err != nil {
		return vipInfo, err
	}
	// 等級提款資訊
	cacheVipMap, ok := cachedata.VipInfoMap()[vipLevel]
	if !ok {
		return vipInfo, nil
	} else {
		vipInfo = cacheVipMap
	}
	return vipInfo, nil
}

// 對應 event_id
func relateEventId(tansSet wallet.TransectionSet) (int, error) {
	transMap, err := typeparam.MainType("tx_src_type").Map()
	if err != nil {
		return 0, err
	}
	r_transMap := typeparam.GetReversedMap(transMap)
	eventNameMap := cachedata.EventNameMap()

	var event_id int
	switch r_transMap[tansSet.TransectionSrcType] {
	case "order":
		event_id = eventNameMap[NORMAL_POINT_EVENT_NAME].Id
	case "admin":
		event_id = eventNameMap[ADMIN_POINT_EVENT_NAME].Id
	case "bonus":
		event_id = *tansSet.TransectionRelate
	}
	return event_id, nil
}
