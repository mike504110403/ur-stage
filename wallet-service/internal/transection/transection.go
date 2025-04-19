package transection

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"wallet_service/internal/redis"
	"wallet_service/internal/transfer"

	mlog "github.com/mike504110403/goutils/log"
)

type Config struct {
	Fequency time.Duration
}

var (
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
)

// Init 啟動交易排程
func Init(initCfg Config) {
	ctx, cancel := context.WithCancel(context.Background())
	cancelFunc = cancel
	go func() {
		frequency := initCfg.Fequency
		for {
			select {
			case <-ctx.Done():
				mlog.Info("交易排程已停止")
				return
			default:
				wg.Add(1)
				if err := DealUnTransectionSet(); err != nil {
					mlog.Error(fmt.Sprintf("交易排成錯誤: %v", err))
				}
				wg.Done()
				time.Sleep(frequency)
			}
		}
	}()
}

// SetPointTranection 設定轉點交易
func SetPointTranection(db *sql.DB, transSet TransectionSetReq) error {
	insertStr := `
		INSERT INTO TransectionSet
		( member_id, amount, gate_amount, transection_src_type, transection_relate )
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := db.Exec(
		insertStr,
		transSet.MemberId,
		transSet.Amount,
		transSet.Amount,
		transSet.TransectionSrcType,
		transSet.TransectionRelate,
	)
	if err != nil {
		return err
	}
	return nil
}

// 發佈轉帳請求
func PublishTransection(transReq transfer.TransferRequest) error {
	conn := redis.RedisPool.Get()
	defer conn.Close()

	reqJSON, err := json.Marshal(transReq)
	if err != nil {
		return err
	}
	if err := conn.Send("PUBLISH", transReq.Type, reqJSON); err != nil {
		return err
	}

	return nil
}
