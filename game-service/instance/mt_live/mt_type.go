package mt_live

import (
	"game_service/pkg/encoder"
	"time"
)

type MtLiveSecrect_Info struct {
	Ischeck     bool
	RefreshTime time.Time
	CE          encoder.CashEncryption
}

// TODO: 所以這邊應該是針對最實用的call api/db 的操作所用的結構
type MemberGameAccount struct {
	UserName     string `db:"username"`
	NickName     string `db:"nick_name"`
	GamePassword string `db:"game_password"`
}

type (
	CreateUserReq struct {
		SystemCode  string `json:"system_code"`   // 系統代碼
		WebId       string `json:"web_id"`        // 站台代碼，即代理唯一識別碼ID
		UserId      string `json:"user_id"`       // 玩家的唯一識別碼
		UserName    string `json:"user_name"`     // 玩家名稱，4~16字元，可使用英文字母、數字等字元
		Currency    string `json:"currency"`      // 幣別
		LimitGroup  *int   `json:"limit_group"`   // (非必填)限注範本(無為預設)
		LimitDayWin *int   `json:"limit_day_win"` // (非必填)每日贏額上限(萬) (無為預設)
	}
	CreateUserRes struct {
		Code      string  `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string  `json:"message"`   // 錯誤訊息
		Timestamp int64   `json:"timestamp"` // 時間戳記
		Data      DataDTO `json:"data"`      // 回傳資料
	}
	DataDTO struct {
		Result int `json:"result"` // 0:失敗, 1:成功
	}
)

type (
	EditUserReq struct {
		SystemCode string `json:"system_code"` // 系統代碼
		WebId      string `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		UserId     string `json:"user_id"`     // 玩家的唯一識別碼
		UserName   string `json:"user_name"`   // 玩家名稱，4~16字元，可使用英文字母、數字等字元
		LimitGroup *int   `json:"limit_group"` // (非必填)限注範本(無為預設)
	}
	EditUserRes struct {
		Code      string  `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string  `json:"message"`   // 錯誤訊息
		Timestamp int64   `json:"timestamp"` // 時間戳記
		Data      DataDTO `json:"data"`      // 回傳資料
	}
)

type (
	CheckUserReq struct {
		SystemCode string `json:"system_code"` // 系統代碼
		WebId      string `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		UserId     string `json:"user_id"`     // 玩家的唯一識別碼
	}
	CheckUserRes struct {
		Code      string           `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string           `json:"message"`   // 錯誤訊息
		Timestamp int64            `json:"timestamp"` // 時間戳記
		Data      CheckUserDataDTO `json:"data"`      // 回傳資料
	}
	CheckUserDataDTO struct {
		Result int `json:"result"` // 0:失敗, 1:成功
	}
)

type (
	GetURLTokenReq struct {
		SystemCode string  `json:"system_code"` // 系統代碼
		WebId      string  `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		UserId     string  `json:"user_id"`     // 玩家的唯一識別碼
		Language   string  `json:"language"`    // 語系
		ExitAction *string `json:"exit_action"` // 離開遊戲時導向特定網址
	}
	GetURLTokenRes struct {
		Code      string             `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string             `json:"message"`   // 錯誤訊息
		Timestamp int64              `json:"timestamp"` // 時間戳記
		Data      GetURLTokenDataDTO `json:"data"`      // 回傳資料
	}
	GetURLTokenDataDTO struct {
		Token string `json:"token"` // token key
		Url   string `json:"url"`   // 遊戲平台網址
	}
)

type (
	DepositReq struct {
		SystemCode string  `json:"system_code"` // 系統代碼
		WebId      string  `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		UserId     string  `json:"user_id"`     // 玩家的唯一識別碼
		Balance    float64 `json:"balance"`     // 存款額度
		TransferId *string `json:"transfer_id"` // 交易編號可空值，空值時系統將自動生成
	}
	DepositRes struct {
		Code      string         `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string         `json:"message"`   // 錯誤訊息
		Timestamp int64          `json:"timestamp"` // 時間戳記
		Data      DepositDataDTO `json:"data"`      // 回傳資料
	}
	DepositDataDTO struct {
		UserId     string `json:"user_id"`     // 玩家的唯一識別碼
		Balance    string `json:"balance"`     // 點數
		TransferId string `json:"transfer_id"` // 交易編號
	}
)

type (
	WithdrawReq struct {
		SystemCode string  `json:"system_code"` // 系統代碼
		WebId      string  `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		UserId     string  `json:"user_id"`     // 玩家的唯一識別碼
		Balance    float64 `json:"balance"`     // 提款額度
		TransferId *string `json:"transfer_id"` // 交易編號可空值，空值時系統將自動生成
	}
	WithdrawRes struct {
		Code      string          `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string          `json:"message"`   // 錯誤訊息
		Timestamp int64           `json:"timestamp"` // 時間戳記
		Data      WithdrawDataDTO `json:"data"`      // 回傳資料
	}
	WithdrawDataDTO struct {
		UserId     string `json:"user_id"`     // 玩家的唯一識別碼
		Balance    string `json:"balance"`     // 點數
		TransferId string `json:"transfer_id"` // 交易編號
	}
)

type (
	GetBalanceReq struct {
		SystemCode string `json:"system_code"` // 系統代碼
		WebId      string `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		UserId     string `json:"user_id"`     // 玩家的唯一識別碼
	}
	GetBalanceRes struct {
		Code      string            `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string            `json:"message"`   // 錯誤訊息
		Timestamp int64             `json:"timestamp"` // 時間戳記
		Data      GetBalanceDataDTO `json:"data"`      // 回傳資料
	}
	GetBalanceDataDTO struct {
		UserId  string `json:"user_id"` // 玩家的唯一識別碼
		Balance string `json:"balance"` // 點數
	}
)

type (
	KickoutReq struct {
		SystemCode string `json:"system_code"` // 系統代碼
		WebId      string `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		UserId     string `json:"user_id"`     // 玩家的唯一識別碼
	}
	KickoutRes struct {
		Code      string         `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string         `json:"message"`   // 錯誤訊息
		Timestamp int64          `json:"timestamp"` // 時間戳記
		Data      KickoutDataDTO `json:"data"`      // 回傳資料
	}
	KickoutDataDTO struct {
		UserId string `json:"user_id"` // 玩家的唯一識別碼
		Result int    `json:"result"`  // 是否踢出結果(0:失敗:1成功)
	}
)

type (
	PlayerOnlineListReq struct {
		SystemCode string `json:"system_code"` // 系統代碼
		WebId      string `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
	}
	PlayerOnlineListRes struct {
		Code      string                  `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string                  `json:"message"`   // 錯誤訊息
		Timestamp int64                   `json:"timestamp"` // 時間戳記
		Data      PlayerOnlineListDataDTO `json:"data"`      // 回傳資料
	}
	PlayerOnlineListDataDTO struct {
		Webid     string `json:"web_id"`     // 站台代碼
		Userid    string `json:"user_id"`    // 玩家ID
		LoginTime string `json:"login_time"` // 登入時間
	}
)

type (
	GetBetRecordReq struct {
		SystemCode string `json:"system_code"` // 系統代碼
		WebId      string `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		UserId     string `json:"user_id"`     // 玩家的唯一識別碼
		StartTime  string `json:"start_time"`  // 開始時間,格式:”2018-01-01 02:00:00”
		EndTime    string `json:"end_time"`    // 結束時間,格式:”2018-01-01 02:00:00”，目前結束時間與起始時間不能相差超過 24 小時
		Page       int    `json:"page"`        // 起始頁，正整數，從 1 開始計數
		PageSize   *int   `json:"page_size"`   // 每頁筆數，正整數，最大值100
	}
	GetBetRecordRes struct {
		Code      string              `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string              `json:"message"`   // 錯誤訊息
		Timestamp int64               `json:"timestamp"` // 時間戳記
		Data      GetBetRecordDataDTO `json:"data"`      // 回傳資料
	}
	GetBetRecordDataDTO struct {
		SystemCode  string    `json:"system_code"`  // 系統代碼
		WebId       string    `json:"web_id"`       // 站台代碼
		CurrentPage int       `json:"current_page"` // 當前頁碼
		TotalPage   int       `json:"total_page"`   // 總頁數
		TotalCount  int       `json:"total_count"`  // 總筆數
		List        []ListDTO `json:"list"`         // 資料列表
	}
	ListDTO struct {
		UserId     string `json:"user_id"`     // 玩家的唯一識別碼
		Sn         string `json:"sn"`          // 注單編號
		GameCode   string `json:"game_code"`   // 遊戲代碼
		TableCode  string `json:"table_code"`  // 桌號
		GameName   string `json:"game_name"`   // 遊戲局號
		PlayCode   string `json:"play_code"`   // 玩法代碼
		PlayName   string `json:"play_name"`   // 玩法名稱
		Odds       string `json:"odds"`        // 賠率
		OrderMoney string `json:"order_money"` // 投注金額
		ValidMoney string `json:"valid_money"` // 有效投注金額
		WinMoney   string `json:"win_money"`   // 中獎金額 (不含本金)
		Profit     string `json:"profit"`      // 輸贏金額
		OrderTime  string `json:"order_time"`  // 下注時間
		SettleTime string `json:"settle_time"` // 結算時間
		SettleDate string `json:"settle_date"` // 歸帳日期
		Status     int    `json:"status"`      // 注單狀態 (0:待結算 ,1:取消 ,2:未中獎,3: 中獎 ,4:和局)
		Currency   string `json:"currency"`    // 幣別
		IP         string `json:"ip"`          // IP
	}
)

type (
	GetTransationRecordReq struct {
		SystemCode string  `json:"system_code"` // 系統代碼
		WebId      string  `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		UserId     *string `json:"user_id"`     // 玩家的唯一識別碼
		StartTime  string  `json:"start_time"`  // 開始時間,格式:”2018-01-01 02:00:00”
		EndTime    string  `json:"end_time"`    // 結束時間,格式:”2018-01-01 02:00:00”，目前結束時間與起始時間不能相差超過 24 小時
		Page       string  `json:"page"`        // 起始頁，正整數，從 1 開始計數
		PageSize   *int    `json:"page_size"`   // 每頁筆數，正整數，最大值100
	}
	GetTransationRecordRes struct {
		Code      string                     `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string                     `json:"message"`   // 錯誤訊息
		Timestamp int64                      `json:"timestamp"` // 時間戳記
		Data      GetTransationRecordDataDTO `json:"data"`      // 回傳資料
	}
	GetTransationRecordDataDTO struct {
		SystemCode  string                           `json:"system_code"`  // 系統代碼
		WebId       string                           `json:"web_id"`       // 站台代碼
		CurrentPage string                           `json:"current_page"` // 當前頁碼
		TotalPage   int                              `json:"total_page"`   // 總頁數
		TotalCount  int                              `json:"total_count"`  // 總筆數
		List        []GetTransationRecordDataListDTO `json:"list"`         // 資料列表
	}
	GetTransationRecordDataListDTO struct {
		Type   string `json:"type"`    // 交易類型(1: 轉入,2: 轉出)
		Sn     string `json:"sn"`      // 交易編號
		Amount string `json:"amount"`  // 交易金額
		Time   string `json:"time"`    // 交易時間
		UserId string `json:"user_id"` // 玩家的唯一識別碼
	}
)

type (
	FindTransationRecordReq struct {
		SystemCode string `json:"system_code"` // 系統代碼
		WebId      string `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		TransferSN string `json:"transfer_sn"` // 交易編號
	}
	FindTransationRecordRes struct {
		Code      string                      `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string                      `json:"message"`   // 錯誤訊息
		Timestamp int64                       `json:"timestamp"` // 時間戳記
		Data      FindTransationRecordDataDTO `json:"data"`      // 回傳資料
	}
	FindTransationRecordDataDTO struct {
		SystemCode string                        `json:"system_code"` // 系統代碼
		WebId      string                        `json:"web_id"`      // 站台代碼
		List       []FindTransationRecordListDTO `json:"list"`        // 資料列表
	}
	FindTransationRecordListDTO struct {
		Type   string `json:"type"`
		SN     string `json:"sn"`
		Amount string `json:"amount"`
		Time   string `json:"time"`
		UserID string `json:"user_id"`
	}
)

// 無法驗證注單異動狀態
type (
	GetUpdateBetRecordReq struct {
		SystemCode string  `json:"system_code"` // 系統代碼
		WebId      string  `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		UserId     *string `json:"user_id"`     // 玩家的唯一識別碼
		StartTime  string  `json:"start_time"`  // 開始時間,格式:”2018-01-01 02:00:00”
		EndTime    string  `json:"end_time"`    // 結束時間,格式:”2018-01-01 02:00:00”，目前結束時間與起始時間不能相差超過 24 小時
		Page       int     `json:"page"`        // 起始頁，正整數，從 1 開始計數
		PageSize   int     `json:"page_size"`   // 每頁筆數，正整數，最大值100
	}
	GetUpdateBetRecordRes struct {
		Code      string                    `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string                    `json:"message"`   // 錯誤訊息
		Timestamp int64                     `json:"timestamp"` // 時間戳記
		Data      GetUpdateBetRecordDataDTO `json:"data"`      // 回傳資料
	}
	GetUpdateBetRecordDataDTO struct {
		SystemCode  string                      `json:"system_code"`  // 系統代碼
		WebId       string                      `json:"web_id"`       // 站台代碼
		CurrentPage int                         `json:"current_page"` // 目前頁碼
		TotalPage   int                         `json:"total_page"`   // 總頁數
		TotalCount  int                         `json:"total_count"`  // 總筆數
		List        []GetUpdateBetRecordListDTO `json:"list"`         // 資料列表
	}
	GetUpdateBetRecordListDTO struct {
		UserId     string  `json:"user_id"`     // 玩家的唯一識別碼
		Sn         string  `json:"sn"`          // 注單編號
		GameCode   string  `json:"game_code"`   // 遊戲代碼
		TableCode  string  `json:"table_code"`  // 桌號
		GameName   string  `json:"game_name"`   // 遊戲局號
		PlayCode   string  `json:"play_code"`   // 玩法代碼
		PlayName   string  `json:"play_name"`   // 玩法名稱
		Odds       string  `json:"odds"`        // 賠率
		OrderMoney float64 `json:"order_money"` // 投注金額
		ValidMoney float64 `json:"valid_money"` // 有效投注金額
		WinMoney   float64 `json:"win_money"`   // 中獎金額 (不含本金)
		Profit     float64 `json:"profit"`      // 輸贏金額
		OrderTime  string  `json:"order_time"`  // 下注時間
		SettleTime string  `json:"settle_time"` // 結算時間
		SettleDate string  `json:"settle_date"` // 歸帳日期
		Status     string  `json:"status"`      // 注單狀態 (0:待結算 ,1:取消 ,2:未中獎,3: 中獎 ,4:和局)
		Currency   string  `json:"currency"`    // 幣別
		IP         string  `json:"ip"`          // IP
	}
)

type (
	DonateRecordReq struct {
		SystemCode string `json:"system_code" binding:"required"`
		WebID      string `json:"web_id" binding:"required"`
		UserID     string `json:"user_id"`
		StartTime  string `json:"start_time" binding:"required"`
		EndTime    string `json:"end_time" binding:"required"`
		Page       int    `json:"page" binding:"required"`
		PageSize   *int   `json:"page_size"`
	}
	DonateRecordRes struct {
		Code      string           `json:"code"`
		Message   string           `json:"message"`
		Timestamp int              `json:"timestamp"`
		Data      DonateRecordData `json:"data"`
	}
	DonateRecordData struct {
		SystemCode  string         `json:"system_code"`
		WebID       string         `json:"web_id"`
		CurrentPage int            `json:"current_page"`
		TotalPage   int            `json:"total_page"`
		TotalCount  int            `json:"total_count"`
		List        []DonateRecord `json:"list"`
	}
	DonateRecord struct {
		UserID     string `json:"user_id"`
		SN         string `json:"sn"`
		GameCode   string `json:"game_code"`
		GameName   string `json:"game_name"`
		TableCode  string `json:"table_code"`
		Money      int    `json:"money"`
		Time       string `json:"time"`
		DealerID   string `json:"dealer_id"`
		DealerName string `json:"delaer_name"`
		GiftID     string `json:"gift_id"`
		GiftName   string `json:"gift_name"`
	}
)

type (
	BetListResult struct {
	}
)
