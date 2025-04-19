package wallet

import (
	"database/sql"
	"errors"
	"wallet_service/internal/cachedata"
	"wallet_service/internal/prewithdraw"

	mlog "github.com/mike504110403/goutils/log"
)

// 提款限制
func getWithdrawGate(mid int) (WithdrawGateRes, error) {
	res := WithdrawGateRes{}
	// 取得提款門檻相關資訊
	vipInfo, err := prewithdraw.GetVipInfoMap(mid)
	if err != nil {
		return res, err
	}
	gateInfo, err := prewithdraw.GetEventGate(mid)
	if err != nil {
		return res, err
	}

	res.WithdrawableTimes = func() int {
		if vipInfo.WithdrawTimes == -1 {
			return -1
		}
		var thisTimes int
		if gateInfo.TihsdayWithdrawTimes != nil {
			thisTimes = *gateInfo.TihsdayWithdrawTimes
		} else {
			thisTimes = 0
		}
		return vipInfo.WithdrawTimes - thisTimes
	}()

	res.WithdrawableAmount = func() float64 {
		if vipInfo.WithdrawLimit == -1 {
			return -1
		}
		var thisAmount float64
		if gateInfo.ThisdayWithdrawAmount != nil {
			thisAmount = *gateInfo.ThisdayWithdrawAmount
		} else {
			thisAmount = 0
		}
		return vipInfo.WithdrawLimit - thisAmount
	}()
	return res, nil
}

func getWithdrawGateFlow(mid int) ([]WithdrawEvent, error) {
	withdrawEvents := []WithdrawEvent{}
	eventMap := cachedata.EventMap()
	events, err := prewithdraw.GetMemberEvents(mid)
	if err != nil {
		return withdrawEvents, err
	}

	for _, event := range events {
		withdrawEvent := WithdrawEvent{
			Name:     eventMap[event.EventId].Name,
			Amount:   event.Amount,
			Gate:     event.Gate,
			GateFlow: event.GateFlow,
			CreateAt: event.CreateAt,
		}
		withdrawEvents = append(withdrawEvents, withdrawEvent)
	}
	return withdrawEvents, nil

}

// 新增轉點紀錄
func insertTransferLogs(tx *sql.Tx, req TransferReq) error {
	// 轉帳紀錄
	transferLogStr := `
		INSERT INTO TransferLogs (member_id, agent_id, amount)
		VALUES (?, ?, ?)
	`
	_, err := tx.Exec(transferLogStr, req.MemberId, req.AgentId, req.Amount)
	if err != nil {
		return err
	}

	return nil
}

// sql 大轉小轉帳紀錄
func mainToSubUpDateSQL(tx *sql.Tx, req TransferReq) error {
	_, err := tx.Exec("UPDATE Wallet SET balance = balance - ? WHERE member_id = ?", req.Amount, req.MemberId)
	if err != nil {
		return err
	}
	gupsertStr := `
		INSERT INTO GameWallet (member_id, agent_id, balance) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE balance = balance + VALUES(balance)
	`
	_, err = tx.Exec(gupsertStr, req.MemberId, req.AgentId, req.Amount)
	if err != nil {
		return err
	}

	if err := insertTransferLogs(tx, TransferReq{
		MemberId: req.MemberId,
		AgentId:  req.AgentId,
		Amount:   -req.Amount,
	}); err != nil {
		return err
	}

	return nil
}

// sql 小轉大轉帳紀錄
func subToMainUpDateSQL(tx *sql.Tx, req TransferReq) error {
	var err error
	//檢查小表是否存在，查詢帳目，比對帳目是否符合轉帳點數，檢查大表示是否存在
	var balance float64
	err = tx.QueryRow("SELECT balance FROM GameWallet WHERE member_id = ? AND agent_id = ?", req.MemberId, req.AgentId).Scan(&balance)
	if err != nil {
		return err
	} else {
		if balance < req.Amount {
			return errors.New("餘額不足")
		} else {
			_, err := tx.Exec("UPDATE GameWallet SET balance = balance - ? WHERE member_id = ?", req.Amount, req.MemberId)
			if err != nil {
				return err
			}
		}
	}

	gupsertStr := `
		UPDATE Wallet SET balance = balance + ? WHERE member_id = ?
	`
	_, err = tx.Exec(gupsertStr, req.Amount, req.MemberId)
	if err != nil {
		return err
	}

	if err := insertTransferLogs(tx, TransferReq{
		MemberId: req.MemberId,
		AgentId:  req.AgentId,
		Amount:   req.Amount,
	}); err != nil {
		return err
	}

	return nil
}

// sql 遊戲錢包轉帳紀錄
func gameToSubUpdateSQL(tx *sql.Tx, req TransferReq) error {
	updateStr := `
		INSERT INTO GameWallet (member_id, agent_id, balance)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
		balance = VALUES(balance);
	`
	_, err := tx.Exec(updateStr, req.MemberId, req.AgentId, req.Amount)
	if err != nil {
		return err
	}

	return nil
}

// sql 大錢包全部轉入遊戲錢包
func bringInUpdateSQL(tx *sql.Tx, req TransferReq) error {
	var balance float64
	queryStr := `
		SELECT balance
		FROM Wallet
		WHERE member_id = ?
		FOR UPDATE
	`
	err := tx.QueryRow(queryStr, req.MemberId).Scan(&balance)
	if err != nil {
		return err
	}
	updateStr := `
		UPDATE Wallet SET balance = 0 WHERE member_id = ?
	`
	_, err = tx.Exec(updateStr, req.MemberId)
	if err != nil {
		return err
	}

	gupsertStr := `
		INSERT INTO GameWallet (member_id, agent_id, balance)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		balance = balance + VALUES(balance)
	`
	_, err = tx.Exec(gupsertStr, req.MemberId, req.AgentId, balance)
	if err != nil {
		return err
	}

	if err := insertTransferLogs(tx, TransferReq{
		MemberId: req.MemberId,
		AgentId:  req.AgentId,
		Amount:   -balance,
	}); err != nil {
		return err
	}

	return nil
}

// sql 遊戲錢包全部轉出
func bringOutUpdateSQL(tx *sql.Tx, req TransferReq) error {
	var totalBalance float64

	// 計算總額
	for aid, a := range req.GameWalletMap {
		totalBalance += a.TransBalance
		// 轉點紀錄
		if err := insertTransferLogs(tx, TransferReq{
			MemberId: req.MemberId,
			AgentId:  aid,
			Amount:   a.TransBalance,
		}); err != nil {
			mlog.Error(err.Error())
			continue
		}
	}
	// 更新主錢包餘額
	updateMainWalletStr := `
		UPDATE Wallet SET balance = balance + ? WHERE member_id = ?
	`
	_, err := tx.Exec(updateMainWalletStr, totalBalance, req.MemberId)
	if err != nil {
		return err
	}

	return nil
}
