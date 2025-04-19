package bonus

import (
	"errors"
	"wallet_service/internal/cachedata"
	"wallet_service/internal/database"

	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

const FIRST_DEPOSIT_EVENT_NAME = "first-deposit"

// 首存優惠
func firstDepositapply(mid int) error {
	// 取得優惠類型
	bonusSrcType := typeparam.TypeParam{
		MainType: typeparam.MainType("tx_src_type"),
		SubType:  typeparam.SubType("bonus"),
	}
	bonusTypeInt, err := bonusSrcType.Get()
	if err != nil {
		return err
	}
	// 取得優惠活動
	eventMap := cachedata.EventNameMap()
	event, ok := eventMap[FIRST_DEPOSIT_EVENT_NAME]
	if !ok {
		return errors.New("優惠活動不存在")
	}
	// 取得首存金額
	amount, err := getFirstDepositAmount(mid)
	if err != nil {
		return err
	}
	// 執行優惠
	pointDb, err := database.POINT.DB()
	if err != nil {
		return err
	}

	if _, err := pointDb.Exec(
		"CALL First_Deposit_BSP(?, ?, ?, ?)",
		mid,
		amount,
		bonusTypeInt,
		event.Id,
	); err != nil {
		return err
	} else {
		var rowsAffected int
		if err = pointDb.QueryRow("SELECT @p_rows_affected").Scan(&rowsAffected); err != nil {
			return err
		}
		if rowsAffected == 0 {
			return errors.New("首存優惠資格不符")
		}
	}

	return nil
}

// 取得首存金額
func getFirstDepositAmount(mid int) (float64, error) {
	// 取得轉點類型
	transTypeMap, err := typeparam.MainType("tx_src_type").Map()
	if err != nil {
		return 0, err
	}
	orderSrcType := int(transTypeMap[typeparam.SubType("order")])
	// adminSrcType := int(transTypeMap[typeparam.SubType("admin")])

	db, err := database.WALLET.DB()
	if err != nil {
		return 0, err
	}
	queryStr := `
		SELECT amount
		FROM TransectionSet
		WHERE 
			member_id = ?
				AND transection_src_type = ?
				AND transected = 1
		ORDER BY trans_done_at ASC
		LIMIT 1
	`
	var amount float64
	if err := db.QueryRow(
		queryStr,
		mid,
		orderSrcType,
		// adminSrcType, // 移除 adminSrcType
	).Scan(&amount); err != nil {
		return 0, err
	}
	return amount, nil
}
