package apollo

import (
	"time"
)

// TODO: 所以這邊應該是針對最實用的call api/db 的操作所用的結構
type MemberGameAccount struct {
	MemberId int    `db:"member_id"`
	Username string `db:"username"`
	NickName string `db:"nickname"`
	Password string `db:"password"`
}

type Apollo struct {
	Member *ApolloMember
}

type Register struct {
	Register *ApolloRegister
}

type ApolloMember struct {
	MemberId    int    `db:"member_id"`
	Username    string `db:"username"`
	MemberState int    `db:"member_state"`
}

type ApolloRegister struct {
	NickName string
	Username string
	Password string
}

type Cash struct {
	Amt      string `json:"amt"`
	DoMan    string `json:"doMan"`
	BetId    string `json:"betId"`
	Remark   string `json:"remark"`
	UpTime   string `json:"upTime"`
	GameType string `json:"gameType"`
}

type CashList struct {
	CashList []Cash `json:"cashlist"`
}

// func AccountExist
type (
	ExistReq struct {
		Agid     int    `json:"agid"`     // 代理ID
		Username string `json:"username"` // 帳號
	}
	ExistRes struct {
		Result string `json:"result"`
	}
)

// func AccountRegister
type (
	RegisterInput struct {
		Username string
		Password string
		Nickname string
	}
	RegisterReq struct {
		Agid     int    `json:"agid"`
		Username string `json:"username"` // 帳號 4~30 碼, 僅可輸入英文字母以及數字的組合,開頭必須為英文字
		Password string `json:"password"` // 密碼(必填) 英數混合， 6~20 碼，無特殊字符
		Nickname string `json:"nickname"` // 玩家昵稱，最大 30 字元英數字夾雜,特殊字元僅可包含 _ . ! # $ *
	}
	RegisterRes struct {
		Status int `json:"status"`
		Error  any `json:"error"`
	}
)

// func Deposit
type (
	DepositInput struct {
		Username string
		Gametype string
		Amt      float64
		Refno    string
		Type     string
	}
	DepositReq struct {
		Agid     int     `json:"agid"`
		Username string  `json:"username"`
		Gametype string  `json:"gametype"` // FH(捕魚額度)、APL_E(電子遊戲)、LO(彩票含真人额度)，預設為彩票額度
		Amt      float64 `json:"amt"`      // 金額
		Refno    string  `json:"refno"`    // 交易編號
		Type     string  `json:"type"`     // 類別 :存入 ,提出
	}
	DepositRes struct {
		Status int           `json:"status"`
		Error  any           `json:"error"`
		Result DepositResult `json:"result"`
	}
	DepositResult struct {
		Refno     string  `json:"refno"`     // 交易編號
		PaymentId string  `json:"paymentId"` // 編碼後交易編號
		Balance   float64 `json:"balance"`   // 餘額
		Amt       float64 `json:"amt"`       // 金額
	}
)

// func GetQuota
type (
	GetquotaReq struct {
		Agid     int    `json:"agid"`
		Username string `json:"username"`
		Gametype string `json:"gametype"` // (非必填) 遊戲代碼 EX:FH(捕魚額度) APL_E(電子遊戲額度) 預設為彩票遊戲
	}
	GetquotaRes struct {
		Status int            `json:"status"`
		Error  any            `json:"error"`
		Result GetquotaResult `json:"result"`
	}
	GetquotaResult struct {
		Balance float64 `json:"balance"` // 額度
	}
)

// func Login
type (
	Login struct {
		Username string
		Password string
		Gametype string
	}
	LoginInput struct {
		Username  string
		Password  string
		Gametype  string
		Language  string
		Ismobile  string
		ReturnUrl string
	}
	LoginReq struct {
		Agid      int    `json:"agid"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		Gametype  string `json:"gametype"`
		Language  string `json:"language"`  // (非必填) 遊戲語系( 預設中文語系 cn:簡中、big5:繁中、en:英文 )
		Ismobile  string `json:"ismobile"`  // (非必填) 是否為手機板（ 若遊戲類型填入LO或MOBILE，則不需要填入。預設為N ）Y / N
		ReturnUrl string `json:"returnUrl"` // 從MOBILE版的遊戲中退出時返回的地址，需使用UrlEncode進行編碼
	}
	LoginRes struct {
		Status int         `json:"status"`
		Error  any         `json:"error"`
		Result LoginResult `json:"result"`
	}
	LoginResult struct {
		Site  string `json:"site"`  // 站台
		Url   string `json:"url"`   // 登入成功後跳轉的頁面
		Query string `json:"query"` // 登入成功請求參數
	}
)

// func BetReport
type (
	// BetReportInput struct {
	// 	Username    string
	// 	Startdate   time.Time
	// 	Enddate     time.Time
	// 	Gametype    string
	// 	Inquirymode any
	// }
	BetReportReq struct {
		Agid        int       `json:"agid"`
		Username    string    `json:"username"`
		Startdate   time.Time `json:"startdate"`   // 開始時間 ex:2014-01-01 01:01:01
		Enddate     time.Time `json:"enddate"`     // 結束時間 ex:2014-01-02 02:02:02
		Gametype    string    `json:"gametype"`    // (非必填) 遊戲代碼 EX:LO(彩票遊戲) FH(捕魚游戲) APL_E(電子遊戲)預設為彩票遊戲
		Inquirymode any       `json:"inquirymode"` // (非必填) 使用報表模式 EX:gr3(注單下注時間抓取注單並統一betlist為陣列)gr2(注單更新時間抓取注單)gr(注單下注時間抓取注單)
	}
	BetReportRes struct {
		Status int             `json:"status"`
		Error  any             `json:"error"`
		Result BetReportResult `json:"result"`
	}
	BetReportResult struct {
		Lottory LotttoryResult `json:"lottory"`
	}
	LotttoryResult struct {
		BetList []BetListResult `json:"betiist"` // TODO:要確認是否為列表
	}
	BetListResult struct {
		CurId      string    `json:"curId"`      // 幣別
		GameId     string    `json:"gameId"`     // 遊戲帳號
		GameType   string    `json:"gameType"`   // 遊戲類別
		GameName   string    `json:"gameName"`   // 遊戲名稱
		Username   string    `json:"username"`   // 會員帳號
		BetId      string    `json:"betId"`      // 注單編號
		EventId    string    `json:"eventId"`    // 期數編號
		Status     int       `json:"status"`     // 注單狀態: 0:未開獎, 1:已開獎, 3:取消單
		BetTime    time.Time `json:"betTime"`    // 下注時間
		Result     string    `json:"result"`     // 開獎結果
		ResultCode string    `json:"resultCode"` // 開獎號碼 PS:代碼參考 遊戲代碼表
		BillTime   time.Time `json:"billTime"`   // 結算時間
		Amt        float64   `json:"amt"`        // 下注金額
		Payout     float64   `json:"payout"`     // 輸贏金額
		PlayId     string    `json:"playId"`     // 玩法代碼 PS:代碼參考 遊戲代碼表
		Number     string    `json:"number"`     // 下注號碼
		Selection  string    `json:"selection"`  // 下注內容
		Odds       float64   `json:"odds"`       // 賠率
		Rollback   float64   `json:"rollback"`   // 退水比
		Ip         string    `json:"ip"`         // 會員下注IP
		Device     string    `json:"device"`     // 會員下注裝置
		Line       string    `json:"line"`       // 下注盤口
	}
)

// func ChangePassword
type (
	AccountChangePasswordReq struct {
		Agid     int    `json:"agid"`
		Username string `json:"username"`
		Password string `json:"password"` // 新密碼 (須為英數混合6~20碼且無特殊字符)
	}
	AccountChangePasswordRes struct {
		Status int `json:"status"`
		Error  any `json:"error"`
	}
)

// func AccountBook
type (
	GetAccountBookReq struct {
		Agid     int    `json:"agid"`
		Username string `json:"username"`
		Gametype string `json:"gametype"` // (非必填) 遊戲代碼 EX:FH(捕魚賬本)LO(彩票含真人賬本)預設為彩票賬本
	}
	GetAccountBookRes struct {
		Status int            `json:"status"`
		Error  any            `json:"error"`
		Result CashListResult `json:"result"`
	}
	CashListResult struct {
		Cashlist CashListDTO `json:"cashlist"`
	}
	CashListDTO struct {
		Amt      float64   `json:"amt"`      // 金額
		DoMan    string    `json:"doMan"`    // 操作者名稱
		BetId    string    `json:"betId"`    // 下注單號
		Remark   string    `json:"remark"`   // 內容
		UpTime   time.Time `json:"upTime"`   // 時間
		GameType string    `json:"gameType"` // 遊戲類別
	}
)

// func GetPeriodResult
type (
	PeriodResultInput struct {
		Agid     int // TODO: 新增agid欄位排除與其他結構衝突
		Gametype string
		Gamedate time.Time
		Gamenum  string
	}
	PeriodResultReq struct {
		Gametype string    `json:"gametype"`
		Gamedate time.Time `json:"gamedate"`
		Gamenum  string    `json:"gamenum"` // 遊戲期數
	}
	PeriodResultRes struct {
		Status int          `json:"status"`
		Error  any          `json:"error"`
		Result ResultResult `json:"result"`
	}
	ResultResult struct {
		Data ResultDTO `json:"data"`
	}
	ResultDTO struct {
		GameNum    string `json:"gameNum"`    // 遊戲期數
		GameResult string `json:"gameResult"` // 遊戲賽果
	}
)

// func GetRealReport
type (
	RealReportInput struct {
		Username  string
		Startdate time.Time
		Enddate   time.Time
		Gametype  string
	}
	RealReportReq struct {
		Agid      int       `json:"agid"`
		Username  string    `json:"username"`
		Startdate time.Time `json:"startdate"` // 開始時間 ex:2014-01-01 01:01:01
		Enddate   time.Time `json:"enddate"`   // 結束時間 ex:2014-01-02 02:02:02
		Gametype  string    `json:"gametype"`  // (非必填) 遊戲代碼 EX:LO(彩票遊戲)FH(捕魚游戲)預設為彩票遊戲
	}
	RealReportRes struct {
		Status string           `json:"status"`
		Error  any              `json:"error"`
		Result RealReportResult `json:"result"`
	}
	RealReportResult struct {
		Lottery RealReportDTO `json:"lottery"`
	}
	RealReportDTO struct {
		Username  string              `json:"username"`
		Ordergold RealReportOrderGold `json:"ordergold"`
	}
	RealReportOrderGold struct {
		GameName string  `json:"gameName"` // 遊戲名稱
		Realgold float64 `json:"realgold"` // 有效投注額
		Win      float64 `json:"win"`      // 輸贏金額
	}
)

// func CheckRefno
type (
	RefnoReq struct {
		Agid     int    `json:"agid"`
		Username string `json:"username"`
		Refno    string `json:"refno"` // 交易編號 自定義ex: 201507010001 ；限30字元以內
	}
	RefnoRes struct {
		Status string      `json:"status"`
		Error  any         `json:"error"`
		Result RefnoResult `json:"result"`
	}
	RefnoResult struct {
		Type string  `json:"type"` // 類別 :存入 ,提出
		Amt  float64 `json:"amt"`  // 金額
	}
)

// func GetTotalRealGold
type (
	TotalRealGoldInput struct {
		Gametype  string
		Startdate time.Time
		Enddate   time.Time
	}
	TotalRealGoldReq struct {
		Agid      int       `json:"agid"`
		Gametype  string    `json:"gametype"`  // (非必填) 遊戲代碼 EX:LO(彩票遊戲)FH(捕魚游戲)預設為彩票遊戲
		Startdate time.Time `json:"startdate"` // 開始時間 ex:2014-01-01 01:01:01
		Enddate   time.Time `json:"enddate"`   // 結束時間 ex:2014-01-02 02:02:02
	}
	TotalRealGoldRes struct {
		Status string              `json:"status"`
		Error  any                 `json:"error"`
		Result TotalRealGoldResult `json:"result"`
	}
	TotalRealGoldResult struct {
		Gametype      string  `json:"gametype"`      // 遊戲代碼(非必填) PS: 代碼參考 遊戲代碼表
		Num           string  `json:"num"`           // 下注筆數
		Totalgold     float64 `json:"totalgold"`     // 總下注額
		Totalrealgold float64 `json:"totalrealgold"` // 總有效投注額
		Totalwingold  float64 `json:"totalwingold"`  // 總派彩
	}
)

// TODO: fhdetaiils struct
type (
	FhDetailsReq struct {
		Agid  int    `json:"agid"`
		Betid string `json:"betid"` // 注單編號
	}
	FhDetailsRes struct {
		Status string          `json:"status"`
		Error  any             `json:"error"`
		Result FhDetailsResult `json:"result"`
	}
	FhDetailsResult struct {
		Id            string    `json:"id"`            // 細單編號
		TableId       string    `json:"tableId"`       // 桌號
		CreateTime    time.Time `json:"createTime"`    // 建立時間
		BeforeBalance float64   `json:"beforeBalance"` // 遊戲前餘額
		AfterBalance  float64   `json:"afterBalance"`  // 遊戲後餘額
		Bet           float64   `json:"bet"`           // 押注
		BetWin        float64   `json:"betWin"`        // 贏分
		WinLoss       float64   `json:"winLoss"`       // 輸贏
		ProcessStatus string    `json:"processStatus"` // 遊戲狀態
		FishSpecies   string    `json:"fishSpecies"`   // 擊中魚種
	}
)

// func GetRecReport
type (
	RecReportInput struct {
		Username  string
		Startdate time.Time
		Enddate   time.Time
	}
	RecReportReq struct {
		Agid      int       `json:"agid"`
		Username  string    `json:"username"`
		Startdate time.Time `json:"startdate"` // 開始時間 ex:2014-01-01 01:01:01
		Enddate   time.Time `json:"enddate"`   // 結束時間 ex:2014-01-02 02:02:02
	}
	RecReportRes struct {
		Status string          `json:"status"`
		Error  any             `json:"error"`
		Result RecReportResult `json:"result"`
	}
	RecReportResult struct {
		TotalNum       string  `json:"totalnum"`       // 總單數
		TotalBetAmount float64 `json:"totalbetamount"` // 總下注金額
		TotalWinGold   float64 `json:"totalwingold"`   // 總派彩
		TotalRealGold  float64 `json:"totalrealgold"`  // 總有效金額
	}
)

// func Logout
type (
	AccountLogoutReq struct {
		Agid     int    `json:"agid"`
		Username string `json:"username"`
	}
	AccountLogoutRes struct {
		Status string `json:"status"`
		Error  any    `json:"error"`
	}
)

// func OnlineStatus
type (
	OnlineStatusReq struct {
		Agid     int    `json:"agid"`
		Username string `json:"username"`
		Gametype string `json:"gametype"` // 遊戲類型。該接口目前只支援捕魚 FH
	}
	OnlineStatusRes struct {
		Status string             `json:"status"`
		Error  any                `json:"error"`
		Result OnlineStatusResult `json:"result"`
	}
	OnlineStatusResult struct {
		OnlineStatus string `json:"onlineStatus"` // 線上狀態 0為離線，1為在線
	}
)

// func GetLine
type (
	GetLineReq struct {
		Agid     int    `json:"agid"`
		Username string `json:"username"`
	}
	GetLineRes struct {
		Status string        `json:"status"`
		Error  any           `json:"error"`
		Result GetLineResult `json:"result"`
	}
	GetLineResult struct {
		Account string `json:"account"` // 帳號
		ALine   string `json:"Aline"`   // A盤口、Y開啟、N關閉
		BLine   string `json:"Bline"`
		CLine   string `json:"Cline"`
		DLine   string `json:"Dline"`
	}
)

// func PreDeposit
type (
	PreDepositInput struct {
		Username string
		Gametype string
		Amt      float64
		Refno    string
		Type     string
	}
	PreDepositReq struct {
		Agid     int     `json:"agid"`
		Username string  `json:"username"`
		Gametype string  `json:"gametype"` // FH(捕魚額度)、APL_E(電子遊戲)、LO(彩券含真人額度)，預設為彩票額度
		Amt      float64 `json:"amt"`      // 額度，建議使用整數int為主。若小於等於0即報錯。
		Refno    string  `json:"refno"`    // 交易編號，自訂ex: aaabbb123456 ；限30字以內
		Type     string  `json:"type"`     // 類別，IN：存入、OUT：提出
	}
	PreDepositRes struct {
		Status string           `json:"status"`
		Error  any              `json:"error"`
		Result PreDepositResult `json:"result"`
	}
	PreDepositResult struct {
		Preno string `json:"preno"` // 此值為ck_deposit 確認預存提介面 所需字符串；存活時間為60秒
	}
)

// func CkDeposit
type (
	CheckDepositReq struct {
		Agid  int    `json:"agid"`
		Preno string `json:"preno"` // 該值需向pre_deposit 預存提接口 取得該字串
	}
	CheckDepositRes struct {
		Status string             `json:"status"`
		Error  any                `json:"error"`
		Result CheckDepositResult `json:"result"`
	}
	CheckDepositResult struct {
		Refno     string  `json:"Refno"`     // 交易編號
		PayMentId string  `json:"paymentId"` // 編碼後交易編號
		Balance   float64 `json:"balance"`   // 餘額
		Amt       float64 `json:"amt"`       // 金額
	}
)
