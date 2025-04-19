package payment

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
	htpay "wallet_service/instance/ht_pay"
	"wallet_service/internal/cachedata"
	"wallet_service/internal/database"
	"wallet_service/internal/prewithdraw"

	wallet "gitlab.com/gogogo2712128/common_moduals/dbModel/Wallet"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"

	"github.com/google/uuid"
	mlog "github.com/mike504110403/goutils/log"
)

// withdraw 提款請求
func withdraw(req withdrawReq) error {
	// 檢查可提款
	if !checkWithdraw(req.MemberId, req.Amount) {
		return errors.New("未滿足提款條件")
	}

	tx, err := database.WALLET.TX()
	if err != nil {
		return err
	}
	wallet, err := getWallet(tx, req.MemberId)
	if err != nil {
		tx.Rollback()
		return err
	}
	fee, err := prewithdraw.GetWithdrawFee(req.MemberId, req.Amount)
	if err != nil {
		tx.Rollback()
		return err
	} else {
		req.Fee = fee
	}
	wallet.LockAmount += req.Amount + fee
	wallet.Balance -= req.Amount + fee

	// 更新餘額
	if err = updateWallet(tx, wallet); err != nil {
		tx.Rollback()
		return err
	}
	// 提款設定
	if err = setWithdraw(tx, req); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// checkWithdraw 提款檢查
func checkWithdraw(mid int, amount float64) bool {
	return prewithdraw.IsWithdrawable(mid, amount)
}

// setWithdraw 提款設定
func setWithdraw(tx *sql.Tx, req withdrawReq) error {
	req.TradeNo = uuid.New().String()
	thirdPayRes, err := thirdPayWithdrawAPI(thirdWithDrawAPIReq{
		TradeNo:      req.TradeNo,
		Amount:       req.Amount,
		AccNo:        req.AccNo,
		WithdrawType: req.WithdrawType,
	})
	if err != nil {
		return err
	}

	// 三方 Id
	thirdPayIdMap := cachedata.ThirdPayIdMap()
	thirdPayIdStr, ok := thirdPayIdMap[req.WithdrawType]
	if !ok {
		return errors.New("無對應的支付方式")
	}
	thirdPayId, err := strconv.Atoi(thirdPayIdStr)
	if err != nil {
		return err
	}
	// 訂單狀態 - 起單
	initWithdraw := typeparam.TypeParam{
		MainType: "withdraw_status",
		SubType:  "init",
	}
	initState, err := initWithdraw.Get()
	if err != nil {
		return err
	}

	insertStr := `
		INSERT INTO WithdrawSet
		( 
			member_id,
			third_pay_id,
			trade_no,
			third_pay_trade_no,
			amount,
			fee,
			status
		)
		VALUES
		( ?, ?, ?, ?, ?, ?, ? )
	`
	_, err = tx.Exec(
		insertStr,
		req.MemberId,
		thirdPayId,
		req.TradeNo,
		thirdPayRes.ThirdPayTradeNo,
		req.Amount,
		req.Fee,
		initState,
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	return nil
}

// thirdPayWithdrawAPI 第三方提款APIs
func thirdPayWithdrawAPI(req thirdWithDrawAPIReq) (thirdWithDrawAPIRes, error) {
	resp := thirdWithDrawAPIRes{}
	resp.TradeNo = req.TradeNo

	switch req.WithdrawType {
	case "ht_trc20", "ht_erc20":
		var extra string
		if req.WithdrawType == "ht_trc20" {
			extra = `{"channel":"TRC20"}`
		} else {
			extra = `{"channel":"ERC20"}`
		}
		usdtReq := htpay.PlaceOrderOutReq{
			AccNo:      req.AccNo,
			Amount:     req.Amount,
			MhtOrderNo: req.TradeNo,
			Extra:      extra,
		}
		usdtRes, err := htpay.PayOutPlaceOrder(usdtReq)
		if err != nil {
			return resp, err
		}
		resp.ThirdPayTradeNo = usdtRes.Result.PfOrderNo
	default:
		return thirdWithDrawAPIRes{}, errors.New("無對應支付方式")
	}
	return resp, nil
}

// doneWithdraw 完成提款
func doneWithdraw(req PayoutResult) error {
	tx, err := database.WALLET.TX()
	if err != nil {
		return err
	}
	// 取得對應提款設定
	withdrawSet, err := getWithdrawSet(tx, req.MhtOrderNo)
	if err != nil {
		tx.Rollback()
		return err
	}
	if withdrawSet.Withdrawed {
		tx.Rollback()
		return errors.New("提款設定已完成")
	}
	// 更新提款設定
	if err = updateWithdrawSet(tx, withdrawSet, req); err != nil {
		tx.Rollback()
		return err
	}
	// 取得錢包
	wallet, err := getWallet(tx, withdrawSet.MemberId)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 更新錢包
	switch PayoutCallBackReqStatus(req.ResultCode) {
	// 0: 成功；1: 失敗
	case PayoutCallBackReqStatusFail:
		wallet.Balance += withdrawSet.Amount + withdrawSet.Fee
		wallet.LockAmount -= withdrawSet.Amount + withdrawSet.Fee
	case PayoutCallBackReqStatusSuccess:
		wallet.LockAmount -= withdrawSet.Amount + withdrawSet.Fee
	}
	// 更新錢包
	if err := updateWallet(tx, wallet); err != nil {
		tx.Rollback()
		return err
	}
	// 更新提款門檻
	if err := prewithdraw.UpdateWithdrawState(tx, withdrawSet.MemberId, withdrawSet.Amount); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// getWithdrawSet 取得提款設定
func getWithdrawSet(tx *sql.Tx, tradeNo string) (wallet.WithdrawSet, error) {
	queryStr := `
		SELECT * FROM WithdrawSet
		WHERE trade_no = ?
	`
	withdrawSet := wallet.WithdrawSet{}
	err := tx.QueryRow(queryStr, tradeNo).Scan(
		&withdrawSet.Id,
		&withdrawSet.MemberId,
		&withdrawSet.TradeNo,
		&withdrawSet.ThirdPayId,
		&withdrawSet.ThirdPayTradeNo,
		&withdrawSet.Amount,
		&withdrawSet.Fee,
		&withdrawSet.Withdrawed,
		&withdrawSet.Status,
		&withdrawSet.WithdrawDoneAt,
		&withdrawSet.CreateAt,
		&withdrawSet.UpdateAt,
	)
	if err != nil {
		tx.Rollback()
		mlog.Error(err.Error())
		return withdrawSet, err
	}
	return withdrawSet, nil
}

// 更新提款設定
func updateWithdrawSet(tx *sql.Tx, withdrawSet wallet.WithdrawSet, req PayoutResult) error {
	status := 0
	switch PayoutCallBackReqStatus(req.ResultCode) {
	// 0: 成功；1: 失敗
	case PayoutCallBackReqStatusFail:
		failType := typeparam.TypeParam{
			MainType: "withdraw_status",
			SubType:  "third_pay_fail",
		}
		failState, err := failType.Get()
		if err != nil {
			return err
		}
		status = failState
	case PayoutCallBackReqStatusSuccess:
		doneType := typeparam.TypeParam{
			MainType: "withdraw_status",
			SubType:  "done",
		}
		doneState, err := doneType.Get()
		if err != nil {
			return err
		}
		status = doneState
	}
	updateStr := `
		UPDATE WithdrawSet
		SET withdrawed = 1, withdraw_done_at = ?, status = ?
		WHERE id = ?
	`
	if _, err := tx.Exec(
		updateStr,
		time.Now(),
		status,
		withdrawSet.Id,
	); err != nil {
		return err
	}
	return nil
}

// 取得會員錢包資料
func getWallet(tx *sql.Tx, mid int) (wallet.Wallet, error) {
	queryStr := `
		SELECT * FROM Wallet
		WHERE member_id = ?
	`
	wallet := wallet.Wallet{}
	err := tx.QueryRow(queryStr, mid).Scan(
		&wallet.MemberId,
		&wallet.Balance,
		&wallet.LockAmount,
		&wallet.CreateAt,
		&wallet.UpdateAt,
	)
	if err != nil {
		return wallet, err
	}
	return wallet, nil
}

// 更新錢包資料
func updateWallet(tx *sql.Tx, wallet wallet.Wallet) error {
	updateStr := `
		UPDATE Wallet
		SET balance = ?,
			lock_amount = ?
		WHERE member_id = ?
	`
	_, err := tx.Exec(updateStr, wallet.Balance, wallet.LockAmount, wallet.MemberId)
	if err != nil {
		tx.Rollback()
		mlog.Error(err.Error())
		return err
	}
	return nil
}
