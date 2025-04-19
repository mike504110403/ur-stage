package wallet

import (
	"database/sql"
	"member_service/internal/database"

	mlog "github.com/mike504110403/goutils/log"
)

func getWalletDetail(mid int) ([]WalletDetail, error) {
	// TODO: 遊戲列表
	result, gameMap, detail := []WalletDetail{}, map[string]float64{}, WalletDetail{}
	db, err := database.WALLET.DB()
	if err != nil {
		mlog.Error("資料庫取得錯誤")
		return result, err
	}
	queryStr := `
		SELECT game_id, balance
		FROM GameWallet
		WHERE member_id = ?
	`
	rows, err := db.Query(queryStr, mid)
	if err != nil {
		// TODO: 這段要拿掉
		if err == sql.ErrNoRows {
			return result, nil
		}
		mlog.Error("資料庫查詢錯誤")
		return result, err
	}
	for rows.Next() {
		err = rows.Scan(&detail.Game, &detail.Balance)
		if err != nil {
			mlog.Error("資料庫查詢錯誤")
			return result, err
		}
		gameMap[detail.Game] = detail.Balance
	}

	for game, balance := range gameMap {
		// TODO: 遊戲id轉名稱
		result = append(result, WalletDetail{
			Game:    game,
			Balance: balance,
		})
	}

	return result, nil
}

// getCenterWallet 取得中心錢包
func getCenterWallet(mid int) (CenterWalletRes, error) {
	centerWallet := CenterWalletRes{}
	db, err := database.WALLET.DB()
	if err != nil {
		mlog.Error("資料庫取得錯誤")
		return centerWallet, err
	}
	queryStr := `
		SELECT balance, lock_amount
		FROM Wallet
		WHERE member_id = ?
	`
	err = db.QueryRow(queryStr, mid).Scan(&centerWallet.Balance, &centerWallet.LockBalance)
	if err != nil {
		mlog.Error("資料庫查詢錯誤")
		return centerWallet, err
	}

	return centerWallet, nil
}
