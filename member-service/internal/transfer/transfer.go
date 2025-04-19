package transfer

import (
	"database/sql"
	"encoding/json"
	"member_service/internal/database"
	"member_service/internal/redis"

	mlog "github.com/mike504110403/goutils/log"
)

// 發佈轉帳請求
func PublishGameToSub(transferChannel string, mid int, amount float64, gameId int) error {
	conn := redis.RedisPool.Get()
	defer conn.Close()

	tx, err := database.WALLET.TX()
	if err != nil {
		mlog.Error("資料庫取得錯誤")
		return err
	}

	req := TransferRequest{
		MemberId: mid,
		Amount:   amount,
		Type:     GameToSub,
		GameId:   gameId,
	}
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return err
	}
	if err := conn.Send("PUBLISH", transferChannel, reqJSON); err != nil {
		return err
	}
	go updateGameWallet(tx, mid, amount, gameId)

	return nil
}

func updateGameWallet(tx *sql.Tx, mid int, amount float64, gameId int) error {
	_, err := tx.Exec("UPDATE GameWallet SET balance = balance + ? WHERE member_id = ? AND game_id = ?", amount, mid, gameId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
