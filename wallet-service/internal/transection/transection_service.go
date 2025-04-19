package transection

import (
	"database/sql"
	"fmt"
	"time"
	"wallet_service/internal/database"
	"wallet_service/internal/prewithdraw"

	wallet "gitlab.com/gogogo2712128/common_moduals/dbModel/Wallet"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"

	mlog "github.com/mike504110403/goutils/log"
)

// DealUnTransectionSet: 處理未交易的交易設定
func DealUnTransectionSet() error {
	db, err := database.WALLET.DB()
	if err != nil {
		mlog.Error(fmt.Sprintf("取得資料庫錯誤: %v", err))
		return err
	}
	if data, err := GetUnTransectionSet(db); err != nil {
		if err != sql.ErrNoRows {
			mlog.Error(fmt.Sprintf("取得未轉點交易失敗: %v", err))
			return err
		}
	} else {
		if err := SetPointTransection(db, data); err != nil {
			mlog.Error(fmt.Sprintf("轉點失敗: %v", err))
			return err
		}

	}
	return nil
}

// GetUnTransectionSet: 取得尚未交易的交易設定
func GetUnTransectionSet(db *sql.DB) (wallet.TransectionSet, error) {
	untrans := wallet.TransectionSet{}
	queryStr := `
		SELECT id, member_id, amount, gate_amount, transection_src_type, transection_relate
		FROM TransectionSet
		WHERE transected = 0
		LIMIT 1	
		FOR UPDATE
		SKIP LOCKED
	`
	if err := db.QueryRow(queryStr).Scan(
		&untrans.Id,
		&untrans.MemberId,
		&untrans.Amount,
		&untrans.GateAmount,
		&untrans.TransectionSrcType,
		&untrans.TransectionRelate,
	); err != nil {
		return untrans, err
	}

	return untrans, nil
}

// PointTransection: 轉點交易操作
func SetPointTransection(db *sql.DB, transSet wallet.TransectionSet) error {
	walletTx, err := db.Begin()
	if err != nil {
		mlog.Error(fmt.Sprintf("transectionCron begin tx error: %v", err))
		return err
	}

	updateSetStr := `
		UPDATE TransectionSet
		SET transected = true, trans_done_at = ?
		WHERE id = ? AND transected = false
	`
	_, err = walletTx.Exec(
		updateSetStr,
		time.Now(),
		transSet.Id,
	)
	if err != nil {
		walletTx.Rollback()
		mlog.Error(fmt.Sprintf("排成入點失敗: %v", err))
		return err
	}

	updateWalletStr := `
		UPDATE Wallet
		SET balance = balance + ?
		WHERE member_id = ?
	`

	_, err = walletTx.Exec(
		updateWalletStr,
		transSet.Amount,
		transSet.MemberId,
	)
	if err != nil {
		walletTx.Rollback()
		mlog.Error(fmt.Sprintf("錢包入點失敗: %v", err))
		return err
	}
	// 訂單完成
	transTypeOrder := typeparam.TypeParam{
		MainType: "tx_src_type",
		SubType:  "order",
	}
	orderType, err := transTypeOrder.Get()
	if err != nil {
		mlog.Error(fmt.Sprintf("取得訂單轉點來源類型失敗: %v", err))
		return err
	}
	// 訂單事務
	orderTx, err := database.ORDER.TX()
	if err != nil {
		mlog.Error(fmt.Sprintf("取得訂單資料庫失敗: %v", err))
		walletTx.Rollback()
		return err
	}
	if orderType == transSet.TransectionSrcType {
		err = doneOrder(orderTx, transSet)
		if err != nil {
			walletTx.Rollback()
			mlog.Error(fmt.Sprintf("訂單完成失敗: %v", err))
			return err
		}
	}
	// 更新提款門檻
	pointTx, err := database.POINT.TX()
	if err != nil {
		mlog.Error(fmt.Sprintf("取得提款資料庫失敗: %v", err))
		walletTx.Rollback()
		orderTx.Rollback()
		return err
	}
	if err = prewithdraw.InsertPointMemberEvent(pointTx, transSet); err != nil {
		walletTx.Rollback()
		return err
	}

	if err := pointTx.Commit(); err != nil {
		mlog.Error(fmt.Sprintf("提款門檻更新失敗: %v", err))
		walletTx.Rollback()
		orderTx.Rollback()
		return err
	}
	if err := walletTx.Commit(); err != nil {
		mlog.Error(fmt.Sprintf("錢包交易完成失敗: %v", err))
		orderTx.Rollback()
		return err
	}
	return orderTx.Commit()
}

// doneOrder: 訂單完成
func doneOrder(tx *sql.Tx, tranSet wallet.TransectionSet) error {
	orderDone := typeparam.TypeParam{
		MainType: "order_status",
		SubType:  "done",
	}
	doneState, err := orderDone.Get()
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	updateStr := `
		UPDATE OrderPurchase SET order_status = ?, order_done_at = ? WHERE id  = ?
	`
	_, err = tx.Exec(
		updateStr,
		doneState,
		tranSet.TransDoneAt,
		tranSet.TransectionRelate,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
