package apollo

import "time"

type BetRecord struct {
	Id            int       `db:"id"`            // 注單紀錄ID
	CurId         string    `db:"curId"`         // 幣別
	GameId        string    `db:"gameId"`        // 遊戲帳號
	GameType      string    `db:"gameType"`      // 遊戲類別
	GameName      string    `db:"gameName"`      // 遊戲名稱
	Username      string    `db:"username"`      // 會員帳號
	BetId         string    `db:"betId"`         // 注單編號 - 會員下注唯一編號 // 設索引
	EventId       string    `db:"eventId"`       // 期數編號
	Status        int       `db:"status"`        // 注單狀態: 0:未開獎, 1:已開獎, 3:取消單
	BetTime       time.Time `db:"betTime"`       // 下注時間
	Result        string    `db:"result"`        // 開獎結果
	ResultCode    string    `db:"resultCode"`    // 開獎號碼 PS:代碼參考 遊戲代碼表
	BillTime      time.Time `db:"billTime"`      // 結算時間
	Amt           float64   `db:"amt"`           // 下注金額
	Payout        float64   `db:"payout"`        // 輸贏金額
	PlayId        string    `db:"playId"`        // 玩法代碼 PS:代碼參考 遊戲代碼表
	Number        string    `db:"number"`        // 下注號碼
	Selection     string    `db:"selection"`     // 下注內容
	Odds          float64   `db:"odds"`          // 賠率
	Rollback      float64   `db:"rollback"`      // 退水比
	Ip            string    `db:"ip"`            // 會員下注IP
	Device        string    `db:"device"`        // 會員下注裝置
	Line          string    `db:"line"`          // 下注盤口
	RecordConfirm bool      `db:"record_status"` // 是否已確認
}

type LocalBetRecord struct {
	Id          int `db:"id"`            // 注單紀錄ID
	BetRecordId int `db:"bet_record_id"` // 遊戲商注單ID
	BetId       int `db:"bet_id"`        // 本地注單ID
}
