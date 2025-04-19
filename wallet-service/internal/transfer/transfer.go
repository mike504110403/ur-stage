package transfer

import (
	"encoding/json"
	"fmt"

	Redis "wallet_service/internal/redis"

	"github.com/gomodule/redigo/redis"
	mlog "github.com/mike504110403/goutils/log"
)

// ProcessTransferRequests : 處理轉帳請求
func ProcessTransferRequests(transferChannel string) {
	rdb := Redis.GetRedisConn()
	defer rdb.Close()

	pubsub := redis.PubSubConn{Conn: rdb}
	if err := pubsub.Subscribe(transferChannel); err != nil {
		mlog.Fatal(fmt.Sprintf("Failed to subscribe to transfer channel: %v", err))
	}
	defer pubsub.Close()
	// 處理訂閱消息
	for {
		switch v := pubsub.Receive().(type) {
		case redis.Message:
			var req TransferRequest
			if err := json.Unmarshal(v.Data, &req); err != nil {
				mlog.Error(fmt.Sprintf("Error unmarshalling transfer request: %v", err))
				continue
			}
			// 轉點事務分流
			switch req.Type {
			case MainToSub:
				if err := transferMainToSub(rdb, req.MemberId, req.GameId, req.Amount); err != nil {
					mlog.Error(fmt.Sprintf("Error executing main-to-sub transfer: %v", err))
				}
			case SubToMain:
				if err := transferSubToMain(rdb, req.MemberId, req.GameId, req.Amount); err != nil {
					mlog.Error(fmt.Sprintf("Error executing sub-to-main transfer: %v", err))
				}
			case GameToSub:
				if err := processGameTransaction(rdb, req.MemberId, req.GameId, req.Amount); err != nil {
					mlog.Error(fmt.Sprintf("Error executing game-to-sub transaction: %v", err))
				}
			case BringOut:
				if err := transferBringOut(rdb, req.MemberId); err != nil {
					mlog.Error(fmt.Sprintf("Error executing bring-out transfer: %v", err))
				}
			case BringIn:
				if err := transferBringIn(rdb, req.MemberId, req.GameId); err != nil {
					mlog.Error(fmt.Sprintf("Error executing bring-in transfer: %v", err))
				}
			case WithDraw:
				amount := req.Amount
				if err := transferTransection(rdb, req.MemberId, amount, string(WithDraw)); err != nil {
					mlog.Error(fmt.Sprintf("Error executing withdraw transfer: %v", err))
				}
			case Transection:
				amount := -req.Amount
				if err := transferTransection(rdb, req.MemberId, amount, string(Transection)); err != nil {
					mlog.Error(fmt.Sprintf("Error executing transection transfer: %v", err))
				}
			default:
				mlog.Error(fmt.Sprintf("Unknown transfer type: %s", req.Type))
			}
		case redis.Subscription:
			mlog.Info(fmt.Sprintf("Subscribed to channel: %s", v.Channel))
		case error:
			mlog.Error(fmt.Sprintf("Error with pubsub connection: %v", v))
			return
		}
	}
}
