package mt_lottery

import (
	"game_service/pkg/encoder"
	"time"
)

type MtLotterySecrect_Info struct {
	Ischeck     bool
	RefreshTime time.Time
	HC          encoder.HashConfig
}

type MemberGameAccount struct {
	UserName     string `db:"username"`
	NickName     string `db:"nick_name"`
	GamePassword string `db:"game_password"`
}

type Response struct {
	Status bool   `json:"status"`
	Code   int    `json:"code"`
	ErrMsg string `json:"errMsg"`
	Data   string `json:"data"`
}

type (
	CreateUserReq struct {
		Account  string `json:"Account"`  // 會員帳號，6~20字串
		Password string `json:"Password"` // 會員密碼，6~20字串
	}
	CreateUserRes struct {
		Rows RowsDTO `json:"rows"`
	}
	RowsDTO struct {
		MemberAccount string `json:"MemberAccount"` // 會員帳號(長度限50)
		Enable        bool   `json:"Enable"`        // true:新增成功 false:代理狀態異常
	}
)

type (
	CheckPointReq struct {
		Account  string `json:"Account"`  // 會員帳號，6~20字串
		Password string `json:"Password"` // 會員密碼，6~20字串
		GameCode string `json:"GameCode"` // 遊戲代碼，請傳字串"lottery"
	}
	CheckPointRes struct {
		Rows CheckPointRowsDTO `json:"rows"`
	}
	CheckPointRowsDTO struct {
		MemberAccount string  `json:"MemberAccount"` // 會員帳號(長度限50)
		MemberPoint   float64 `json:"MemberPoint"`   // 會員可用點數(12 + 小數點4位)
	}
)

type (
	TransPointReq struct {
		Account    string  `json:"Account"`    // 會員帳號，6~20字串
		Password   string  `json:"Password"`   // 會員密碼，6~20字串
		Point      float64 `json:"Points"`     // 點數，正數:轉入 負數:轉出
		TradeOrder string  `json:"TradeOrder"` // 交易單號
	}
	TransPointRes struct {
		Rows TransPointRowsDTO `json:"rows"`
	}
	TransPointRowsDTO struct {
		MemberAccount     string  `json:"MemberAccount"`     // 會員帳號(長度限50)
		TranPoint         float64 `json:"TranPoint"`         // 轉入(出)點數(10 + 小數點4位)
		BeforeChangePoint float64 `json:"BeforeChangePoint"` // 轉入(出)後點數(10 + 小數點4位)
		AfterChangePoint  float64 `json:"AfterChangePoint"`  // 轉入(出)前點數(10 + 小數點4位)
		DateTran          string  `json:"DateTran"`          // 轉入(出)時間(長度限19)
		TranOrder         int     `json:"TranOrder"`         // Meta 所產生之交易單號(長度限19)
	}
)

type (
	TransactionLogReq struct {
		Account    string  `json:"Account"`    // 會員帳號，6~20字串
		Date       int     `json:"Date"`       // 查詢日期，時間戳 (Unix timestamp)
		Limit      *int    `json:"Limit"`      // 預設 1000
		TranOrder  *int    `json:"TranOrder"`  // 交易單號
		TradeOrder *string `json:"TradeOrder"` // 轉點時，填寫的 TradeOrder 單號
	}
	TransactionLogRes struct {
		TotalRows int                 `json:"totalRows"` // 總筆數
		OverRows  int                 `json:"overRows"`  // 剩餘筆數
		Limit     int                 `json:"limit"`     // 查詢筆數
		Rows      []TransactionLogDTO `json:"rows"`
	}
	TransactionLogDTO struct {
		Account    string `json:"Account"`    // 會員帳號(長度限50)
		TranPoint  string `json:"TranPoint"`  // 轉入(出)點數(10 + 小數點4位)，交易金額 正數:儲值 負數:提款
		DateTran   string `json:"DateTran"`   // 交易時間(長度限19)
		TranOrder  int64  `json:"TranOrder"`  // Meta 所產生之交易單號
		TradeOrder string `json:"TradeOrder"` // 介接商提供之交易單號
	}
)

type (
	LoginReq struct {
		Account     string `json:"Account"`     // 會員帳號，6~20字串
		Password    string `json:"Password"`    // 會員密碼，6~20字串
		GameCode    string `json:"GameCode"`    // 遊戲代碼
		RedirectUrl string `json:"RedirectUrl"` // 客戶欲返回自家平台網址
		Lang        string `json:"Lang"`        // 語系
	}
	LoginRes struct {
		TotalRows int        `json:"totalRows"` // 總筆數
		Rows      []LoginDTO `json:"rows"`
	}
	LoginDTO struct {
		Token string `json:"Token"` // 進入遊戲的 token
		Url   string `json:"Url"`   // 進入遊戲的連結
	}
)

type (
	LogoutReq struct {
		Account  string `json:"Account"`  // 會員帳號
		Password string `json:"Password"` // 會員密碼
		GameCode string `json:"GameCode"` // 遊戲代碼
	}
	LogoutRes struct {
		Logout bool `json:"Logout"` // true:登出成功 false:代理狀態異常
	}
)

type (
	HandicapCheckReq struct {
		Account string `json:"Account"` // 會員帳號
	}
	AgentHandicap struct {
		Id     string `json:"Id"`     // 限紅設定使用編號
		Name   string `json:"Name"`   // 名稱
		Type   string `json:"Type"`   // 1: 一般 2: VIP
		BetMin string `json:"BetMin"` // 最小下注
		BetMax string `json:"BetMax"` // 最大下注
	}
	MemberHandicap struct {
		Type    string `json:"Type"`    // 1: 一般 2: VIP
		Content string `json:"Content"` // 盤口 ID
	}
	HandicapCheckRes struct {
		AgentHandicap  []AgentHandicap  `json:"AgentHandicap"`  // 代理盤口列表
		MemberHandicap []MemberHandicap `json:"MemberHandicap"` // 會員盤口列表
	}
)

type (
	MemberHandicapSettingReq struct {
		Account string `json:"Account"` // 會員帳號
		Normal  string `json:"normal"`  // 一般限紅
		VIP     string `json:"vip"`     // VIP 限紅
	}
	MemberHandicapSettingRes struct {
		Rows struct {
			Success bool `json:"Success"` // true:設定成功 false:設定失敗
		} `json:"rows"`
	}
)

type (
	BetOrderReq struct {
		Date       int64   `json:"Date"`       // 查詢日期 (Unix timestamp)
		Account    *string `json:"Account"`    // 會員帳號
		Limit      *int    `json:"Limit"`      // 筆數
		LastSerial *string `json:"LastSerial"` // 最後一筆流水序號
		GameTypeId int     `json:"GameTypeId"` // 遊戲類別 ID
		GameId     int     `json:"GameId"`     // 遊戲 ID
		Collect    *int    `json:"Collect"`    // 是否已對獎
	}

	BetOrderRes struct {
		TotalRows int           `json:"totalRows"` // 總筆數
		OverRows  int           `json:"overRows"`  // 剩餘筆數
		Limit     int           `json:"limit"`     // 查詢筆數
		BetTotal  int           `json:"betTotal"`  // 投注總金額
		WinTotal  int           `json:"winTotal"`  // 玩家輸贏總金額
		BetCount  int           `json:"betCount"`  // 注單總數量
		Rows      []BetOrderDTO `json:"rows"`      // 資料列
	}

	BetOrderDTO struct {
		Serial           string  `json:"Serial"`           // 流水編號
		No               string  `json:"No"`               // 注單號
		GameTypeId       string  `json:"GameTypeId"`       // 遊戲類別 ID
		GameId           string  `json:"GameId"`           // 遊戲 ID
		LotteryType      string  `json:"LotteryType"`      // 彩種代碼
		LotteryPlayGroup string  `json:"LotteryPlayGroup"` // 玩法群組編碼
		LotteryPlay      string  `json:"LotteryPlay"`      // 玩法代碼
		Issue            string  `json:"Issue"`            // 期數
		Account          string  `json:"Account"`          // 會員帳號
		Rate             string  `json:"Rate"`             // 匯率
		BetTotal         int     `json:"BetTotal"`         // 下注金額
		BetCount         string  `json:"BetCount"`         // 下注數
		Odds             string  `json:"Odds"`             // 下注時賠率
		Odds2            string  `json:"Odds2"`            // 第二組賠率
		Odds3            string  `json:"Odds3"`            // 第三組賠率
		Winnings         int     `json:"Winnings"`         // 實際中獎金額
		AfterPoints      float64 `json:"AfterPoints"`      // 投注後金額
		OpenResult       *string `json:"OpenResult"`       // 開獎號碼
		OpenResult2      *string `json:"OpenResult2"`      // 第二組開獎號碼
		Type             string  `json:"Type"`             // 下注類別 (0: 正常 / 1:追單)
		Collect          string  `json:"Collect"`          // 是否已對獎 (0:未對獎 1:已對獎)
		BetValid         int     `json:"BetValid"`         // 有效下注
		LotteryContent   string  `json:"LotteryContent"`   // 下注項目
		Status           string  `json:"Status"`           // 狀態 (1:取消 2:未開獎 3:未中獎 4:中獎 5:和局)
		Chase            *string `json:"Chase"`            // 追單 UUID
		DateCurrent      string  `json:"DateCurrent"`      // 本期所屬日期
		DateClosing      string  `json:"DateClosing"`      // 封盤時間
		DateDraw         string  `json:"DateDraw"`         // 開獎時間
		DateCreate       string  `json:"DateCreate"`       // 下注時間
		DateUpdate       string  `json:"DateUpdate"`       // 修改時間
		EffectiveBet     int     `json:"EffectiveBet"`     // 有效下注
	}
)

type (
	BetOrderV2Req struct {
		StartTime  int64   `json:"StartTime"`  // 開始時間 (Unix timestamp)
		EndTime    int64   `json:"EndTime"`    // 結束時間 (Unix timestamp)
		Account    string  `json:"Account"`    // 會員帳號
		Limit      *int    `json:"Limit"`      // 筆數
		LastSerial *string `json:"LastSerial"` // 最後一筆流水序號
		GameTypeId int     `json:"GameTypeId"` // 遊戲類別 ID
		GameId     int     `json:"GameId"`     // 遊戲 ID
		Collect    *int    `json:"Collect"`    // 是否已對獎
		DateType   int     `json:"DateType"`   // 查詢時間模式 (1:下注時間, 2:修改時間)
	}

	BetOrderV2Res struct {
		TotalRows int             `json:"totalRows"` // 總筆數
		OverRows  int             `json:"overRows"`  // 剩餘筆數
		Limit     int             `json:"limit"`     // 查詢筆數
		BetTotal  float64         `json:"betTotal"`  // 投注總金額
		WinTotal  float64         `json:"winTotal"`  // 玩家輸贏總金額
		BetCount  int             `json:"betCount"`  // 注單總數量
		Rows      []BetOrderV2DTO `json:"rows"`      // 資料列
	}

	BetOrderV2DTO struct {
		Serial           string  `json:"Serial"`           // 流水編號
		No               string  `json:"No"`               // 注單號
		GameTypeId       string  `json:"GameTypeId"`       // 遊戲類別 ID
		GameId           string  `json:"GameId"`           // 遊戲 ID
		LotteryType      string  `json:"LotteryType"`      // 彩種代碼
		LotteryPlayGroup string  `json:"LotteryPlayGroup"` // 玩法群組編碼
		LotteryPlay      string  `json:"LotteryPlay"`      // 玩法代碼
		Issue            string  `json:"Issue"`            // 期數
		Account          string  `json:"Account"`          // 會員帳號
		Rate             string  `json:"Rate"`             // 匯率
		BetTotal         float64 `json:"BetTotal"`         // 下注金額
		BetCount         string  `json:"BetCount"`         // 下注數
		Odds             string  `json:"Odds"`             // 下注時賠率
		Odds2            string  `json:"Odds2"`            // 第二組賠率
		Odds3            string  `json:"Odds3"`            // 第三組賠率
		Winnings         float64 `json:"Winnings"`         // 實際中獎金額
		AfterPoints      float64 `json:"AfterPoints"`      // 投注後金額
		OpenResult       *string `json:"OpenResult"`       // 開獎號碼
		OpenResult2      *string `json:"OpenResult2"`      // 第二組開獎號碼
		Type             string  `json:"Type"`             // 下注類別 (0: 正常 / 1:追單)
		Collect          string  `json:"Collect"`          // 是否已對獎 (0:未對獎 1:已對獎)
		BetValid         float64 `json:"BetValid"`         // 有效下注
		LotteryContent   string  `json:"LotteryContent"`   // 下注項目
		Status           string  `json:"Status"`           // 狀態 (1:取消 2:未開獎 3:未中獎 4:中獎 5:和局)
		Chase            *string `json:"Chase"`            // 追單 UUID
		DateCurrent      string  `json:"DateCurrent"`      // 本期所屬日期
		DateClosing      string  `json:"DateClosing"`      // 封盤時間
		DateDraw         string  `json:"DateDraw"`         // 開獎時間
		DateCreate       string  `json:"DateCreate"`       // 下注時間
		DateUpdate       string  `json:"DateUpdate"`       // 修改時間
		EffectiveBet     float64 `json:"EffectiveBet"`     // 有效下注
	}
)

type (
	CheckDateModifyReq struct {
		Date       int64 `json:"Date"`       // 查詢日期 (Unix timestamp)
		GameTypeId int   `json:"GameTypeId"` // 遊戲類別 ID
		GameId     int   `json:"GameId"`     // 遊戲 ID
	}

	CheckDateModifyRes struct {
		TotalRows int                  `json:"totalRows"` // 總筆數
		OverRows  int                  `json:"overRows"`  // 剩餘筆數
		Limit     int                  `json:"limit"`     // 查詢筆數
		BetTotal  int                  `json:"betTotal"`  // 投注總金額
		WinTotal  int                  `json:"winTotal"`  // 玩家輸贏總金額
		BetCount  int                  `json:"betCount"`  // 注單總數量
		Rows      []CheckDateModifyDTO `json:"rows"`      // 資料列
	}

	CheckDateModifyDTO struct {
		Serial           string  `json:"Serial"`           // 流水編號
		No               string  `json:"No"`               // 注單號
		GameTypeId       string  `json:"GameTypeId"`       // 遊戲類別 ID
		GameId           string  `json:"GameId"`           // 遊戲 ID
		LotteryType      string  `json:"LotteryType"`      // 彩種代碼
		LotteryPlayGroup string  `json:"LotteryPlayGroup"` // 玩法群組編碼
		LotteryPlay      string  `json:"LotteryPlay"`      // 玩法代碼
		Issue            string  `json:"Issue"`            // 期數
		Account          string  `json:"Account"`          // 會員帳號
		Rate             string  `json:"Rate"`             // 匯率
		BetTotal         int     `json:"BetTotal"`         // 下注金額
		BetCount         string  `json:"BetCount"`         // 下注數
		Odds             string  `json:"Odds"`             // 下注時賠率
		Odds2            string  `json:"Odds2"`            // 第二組賠率
		Odds3            string  `json:"Odds3"`            // 第三組賠率
		Winnings         int     `json:"Winnings"`         // 實際中獎金額
		AfterPoints      int     `json:"AfterPoints"`      // 投注後金額
		OpenResult       *string `json:"OpenResult"`       // 開獎號碼
		OpenResult2      *string `json:"OpenResult2"`      // 第二組開獎號碼
		Type             string  `json:"Type"`             // 下注類別 (0: 正常 / 1:追單)
		Collect          string  `json:"Collect"`          // 是否已對獎 (0:未對獎 1:已對獎)
		BetValid         int     `json:"BetValid"`         // 有效下注
		LotteryContent   string  `json:"LotteryContent"`   // 下注項目
		Status           string  `json:"Status"`           // 狀態 (1:取消 2:未開獎 3:未中獎 4:中獎 5:和局)
		Chase            *string `json:"Chase"`            // 追單 UUID
		DateCurrent      string  `json:"DateCurrent"`      // 本期所屬日期
		DateClosing      string  `json:"DateClosing"`      // 封盤時間
		DateDraw         string  `json:"DateDraw"`         // 開獎時間
		DateCreate       string  `json:"DateCreate"`       // 下注時間
		DateUpdate       string  `json:"DateUpdate"`       // 修改時間
		EffectiveBet     int     `json:"EffectiveBet"`     // 有效下注
	}
)

type (
	ModifyNickNameReq struct {
		Account  string `json:"Account"`  // 會員帳號
		Password string `json:"Password"` // 會員密碼
		NickName string `json:"NickName"` // 玩家新暱稱
	}

	ModifyNickNameRes struct {
		Rows struct {
			Success bool `json:"Success"` // true:修改成功 false:狀態異常
		} `json:"rows"`
	}
)

type (
	ModifyPasswordReq struct {
		Account     string `json:"Account"`     // 會員帳號
		Password    string `json:"Password"`    // 會員密碼
		NewPassword string `json:"NewPassword"` // 玩家新密碼
	}

	ModifyPasswordRes struct {
		Rows struct {
			Success bool `json:"Success"` // true:修改成功 false:狀態異常
		} `json:"rows"`
	}
)
