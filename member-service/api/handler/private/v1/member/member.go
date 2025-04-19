package member

import (
	"database/sql"
	"member_service/internal/database"
	"sort"
	"time"

	mlog "github.com/mike504110403/goutils/log"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

// GetMemberStatus: 取得會員基本狀態
func GetMemberStatus(mid int) (StatusRes, error) {
	status := StatusRes{}
	memberState := 0

	memberDb, err := database.MEMBER.DB()
	if err != nil {
		return status, err
	}
	queryStr := `
		SELECT m.id, l.vip_level, l.member_level, m.member_state
		FROM Member as m
		JOIN MemberLevel as l ON m.id = l.member_id
		WHERE m.id = ?
	`
	err = memberDb.QueryRow(queryStr, mid).Scan(&status.MemberId, &status.VipLevel, &status.MemberLevel, &memberState)
	if err != nil {
		return status, err
	}

	walletDb, err := database.WALLET.DB()
	if err != nil {
		return status, err
	}
	queryStr = `
		SELECT balance, lock_amount
		FROM Wallet
		WHERE member_id = ?
	`
	err = walletDb.QueryRow(queryStr, mid).Scan(&status.Balance, &status.LockBalance)
	if err != nil {
		return status, err
	}
	state := string(typeparam.Find(memberState).SubType)
	switch state {
	case "enable":
		status.Status = string(Enable)
		status.Role = string(Normal)
	case "disable":
		status.Status = string(Disable)
		status.Role = string(Normal)
	case "admin":
		status.Status = string(Enable)
		status.Role = string(Admin)
	case "tester":
		status.Status = string(Enable)
		status.Role = string(Tester)
	default:
		status.Status = string(Disable)
		status.Role = string(Normal)
	}

	return status, nil
}

// GetMemberInfo: 取得會員基本資料
func GetMemberInfo(mid int) (InfoRes, error) {
	info := InfoRes{}
	info.MemberId = mid
	db, err := database.MEMBER.DB()
	if err != nil {
		mlog.Error("資料庫取得錯誤")
		return info, err
	}
	queryStr := `
		SELECT
			m.username,
			mi.nick_name, 
			mi.name,
			mi.gender,
			mi.birthday, 
			m.create_at, 
			ms.phone, 
			ms.email
		FROM 
			Member as m
		JOIN 
			MemberInfo as mi ON m.id = mi.member_id
		JOIN 
			Security as ms ON m.id = ms.member_id
		WHERE 
			m.id = ?
	`
	err = db.QueryRow(queryStr, mid).Scan(
		&info.Basic.Username,
		&info.Basic.NickName,
		&info.Basic.Name,
		&info.Basic.Gender,
		&info.Basic.Birthday,
		&info.Basic.RegisterDate,
		&info.Security.Phone,
		&info.Security.Email,
	)
	if err != nil {
		mlog.Error("資料庫查詢錯誤")
		return info, err
	}
	return info, nil
}

// EditMemberInfo: 更新會員基本資料
func EditMemberInfo(mid int, info BasicInfo) error {
	db, err := database.MEMBER.DB()
	if err != nil {
		mlog.Error("資料庫取得錯誤")
		return err
	}
	queryStr := `
		UPDATE MemberInfo
		SET name = ?, nick_name = ?, gender = ?, birthday = ?
		WHERE member_id = ?
	`
	_, err = db.Exec(queryStr, info.Name, info.NickName, info.Gender, info.Birthday, mid)
	if err != nil {
		mlog.Error("資料庫更新錯誤")
		return err
	}
	return nil
}

// GetTransRecords: 取得會員交易紀錄
func GetTransRecords(mid int, req TransRecordReq) ([]TransRecordRes, error) {
	result := []TransRecordRes{}
	if req.Type == nil {
		orderRecords, err := getOrderRecords(mid, req)
		if err != nil {
			return result, err
		}
		withdrawRecords, err := getWithdrawRecords(mid, req)
		if err != nil {
			return result, err
		}
		result = append(result, orderRecords...)
		result = append(result, withdrawRecords...)
	} else if *req.Type == 1 {
		withdrawRecords, err := getWithdrawRecords(mid, req)
		if err != nil && err != sql.ErrNoRows {
			mlog.Error("資料庫取得錯誤")
			return result, err
		}
		result = append(result, withdrawRecords...)
	} else {
		orderRecords, err := getOrderRecords(mid, req)
		if err != nil && err != sql.ErrNoRows {
			mlog.Error("資料庫取得錯誤")
			return result, err
		}
		result = append(result, orderRecords...)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TransDate.After(result[j].TransDate)
	})
	return result, nil
}

// getOrderRecords: 取得會員交易紀錄
func getOrderRecords(mid int, req TransRecordReq) ([]TransRecordRes, error) {
	result, trans := []TransRecordRes{}, TransRecordRes{}
	db, err := database.ORDER.DB()
	if err != nil {
		mlog.Error("資料庫取得錯誤")
		return result, err
	}
	// TODO: hard code 到時候要改
	queryStr := `
		SELECT 
    	CASE
        	WHEN o.id IS NULL THEN t.create_at
        	ELSE o.created_at
    	END AS created_at,
    	CASE
        	WHEN o.id IS NULL THEN ''
        	ELSE  o.trade_no
    	END AS trade_no,
    	t.trans_done_at,
    	t.amount,
    	CASE
        	WHEN o.id IS NULL THEN 8
        	ELSE o.order_status
    	END AS order_status
		FROM
    		Wallet.TransectionSet AS t
        LEFT JOIN
    		OrderPurchase AS o ON t.transection_relate = o.id
        AND transection_src_type = 9
		WHERE 
			t.member_id = ? 
			AND (
				(o.id IS NULL AND t.create_at >= ? AND t.create_at <= ?)
				OR (o.id IS NOT NULL AND o.created_at >= ? AND o.created_at <= ?)
			)
	`
	rows, err := db.Query(queryStr, mid, req.StartDate, req.EndDate, req.StartDate, req.EndDate)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var trade_no string
		var create_at time.Time
		var order_done_at *time.Time
		var amount float64
		var order_status int
		err = rows.Scan(
			&create_at,
			&trade_no,
			&order_done_at,
			&amount,
			&order_status,
		)
		trans.TransDate = create_at
		trans.OrderSeq = trade_no
		trans.TransDoneDate = order_done_at
		trans.Amount = amount
		switch string(typeparam.Find(order_status).SubType) {
		case "init":
			trans.Status = "待处理"
		case "done":
			trans.Status = "已完成"
		}
		trans.TransType = "存点"
		if err != nil {
			mlog.Error("資料庫取得錯誤")
			return result, err
		}
		result = append(result, trans)
	}
	return result, nil
}

// getWithdrawRecords: 取得會員提款紀錄
func getWithdrawRecords(mid int, req TransRecordReq) ([]TransRecordRes, error) {
	result, trans := []TransRecordRes{}, TransRecordRes{}

	db, err := database.WALLET.DB()
	if err != nil {
		mlog.Error("資料庫取得錯誤")
		return []TransRecordRes{}, err
	}
	queryStr := `
		SELECT create_at, trade_no, withdraw_done_at, amount, fee, status
		FROM WithdrawSet
		WHERE member_id = ? 
			AND create_at >= ? 
			AND create_at <= ?
	`
	rows, err := db.Query(queryStr, mid, req.StartDate, req.EndDate)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var create_at time.Time
		var trade_no string
		var withdraw_done_at *time.Time
		var amount, fee float64
		var status int
		err = rows.Scan(
			&create_at,
			&trade_no,
			&withdraw_done_at,
			&amount,
			&fee,
			&status,
		)
		trans.TransDate = create_at
		trans.OrderSeq = trade_no
		trans.TransDoneDate = withdraw_done_at
		trans.Amount = amount
		trans.Fee = fee
		switch string(typeparam.Find(status).SubType) {
		case "init":
			trans.Status = "待处理"
		case "done":
			trans.Status = "已完成"
		}

		trans.TransType = "托售"
		if err != nil {
			mlog.Error("資料庫取得錯誤")
			return result, err
		}
		result = append(result, trans)
	}
	return result, nil
}

// GetBetRecords: 取得會員投注紀錄
func GetBetRecords(mid int, req BetRecordReq) ([]BetRecordRes, error) {
	result, bet := []BetRecordRes{}, BetRecordRes{}
	gameType := typeparam.TypeParam{
		MainType: "game_type",
		SubType:  "live_donate",
	}
	donateType, err := gameType.Get()
	if err != nil {
		mlog.Error("資料庫取得錯誤")
		return result, err
	}
	db, err := database.GAME.DB()
	if err != nil {
		mlog.Error("資料庫取得錯誤")
		return result, err
	}
	queryStr := `
		SELECT game_type_id, bet, win_lose, bet_at
		FROM BetRecord
		WHERE member_id = ? AND game_type_id != ? AND bet_at >= ? AND bet_at <= ?
	`
	var rows *sql.Rows
	if req.GameType != nil {
		queryStr += " AND game_type_id = ?"
		rows, err = db.Query(queryStr, mid, donateType, req.StartDate, req.EndDate, *req.GameType)
		if err != nil {
			mlog.Error("資料庫查詢錯誤")
			return result, err
		}
	} else {
		rows, err = db.Query(queryStr, mid, donateType, req.StartDate, req.EndDate)
		if err != nil {
			mlog.Error("資料庫查詢錯誤")
			return result, err
		}
	}

	for rows.Next() {
		err = rows.Scan(
			&bet.GameTypeId,
			&bet.EffectBet,
			&bet.WinLose,
			&bet.BetAt,
		)
		if err != nil {
			mlog.Error("資料庫取得錯誤")
			return result, err
		}
		switch string(typeparam.Find(bet.GameTypeId).SubType) {
		case "elect":
			bet.GameType = "电子"
		case "lottery":
			bet.GameType = "彩票"
		case "sport":
			bet.GameType = "体育"
		case "live":
			bet.GameType = "真人"
		}
		result = append(result, bet)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].BetAt.After(result[j].BetAt)
	})
	return result, nil
}
