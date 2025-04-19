package wallet

import (
	"wallet_service/internal/database"
	"wallet_service/internal/transfer"

	mlog "github.com/mike504110403/goutils/log"
)

// 轉帳並記錄
func transferAndRecord(req TransferReq) error {
	// 開啟 SQL 交易
	tx, err := database.WALLET.TX()
	if err != nil {
		mlog.Error("資料庫取得錯誤")
		return err
	}

	switch transfer.TransferType(req.Type) {
	// 遊戲點數更新
	case transfer.GameToSub:
		if err := gameToSubUpdateSQL(tx, req); err != nil {
			tx.Rollback()
			return err
		}
	// 大轉小
	case transfer.MainToSub:
		if err := mainToSubUpDateSQL(tx, req); err != nil {
			tx.Rollback()
			return err
		}
	// 小轉大
	case transfer.SubToMain:
		if err := subToMainUpDateSQL(tx, req); err != nil {
			tx.Rollback()
			return err
		}
	// 全吸入遊戲
	case transfer.BringIn:
		if err := bringInUpdateSQL(tx, req); err != nil {
			tx.Rollback()
			return err
		}
	// 全吸出遊戲
	case transfer.BringOut:
		if err := bringOutUpdateSQL(tx, req); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交 SQL 交易
	if err := tx.Commit(); err != nil {
		mlog.Error("交易提交錯誤")
		return err
	}

	return nil
}
