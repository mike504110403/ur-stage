package sa_live

import "time"

type SaLiveSecrect_Info struct {
	Ischeck     bool
	RefreshTime time.Time
	SE          Secret
}
type Secret struct {
	SecretKey    string `json:"secretKey"`
	Md5Key       string `json:"md5Key"`
	EncryptKey   string `json:"encryptKey"`
	CurrencyType string `json:"currencyType"`
	LobbyCode    string `json:"lobbyCode"`
	Lang         string `json:"lang"`
	LoginURL     string `json:"loginURL"`
	HttpDomain   string `json:"httpDomain"`
}

type (
	RegUserInfoReq struct {
		Method       string `json:"method"`
		Key          string `json:"Key"`
		Time         string `json:"Time"`
		Username     string `json:"Username"`
		CurrencyType string `json:"CurrencyType"`
	}
	RegUserInfoRes struct {
		ErrorMsgId int    `xml:"ErrorMsgId"`
		ErrorMsg   string `xml:"ErrorMsg"`
	}
)

type (
	LoginUrl struct {
		Username string `json:"username"`
		Token    string `json:"token"`
		Lobby    string `json:"lobby"`
		Lang     string `json:"lang"`
	}
	LoginReq struct {
		Method       string `json:"method"`
		Key          string `json:"Key"`
		Time         string `json:"Time"`
		Username     string `json:"Username"`
		CurrencyType string `json:"CurrencyType"`
	}
	LoginRes struct {
		Token       string `xml:"Token"`
		DisplayName string `xml:"DisplayName"`
		ErrorMsgId  int    `xml:"ErrorMsgId"`
		ErrorMsg    string `xml:"ErrorMsg"`
	}
)

type (
	// 請求結構體
	VerifyUsernameReq struct {
		Method   string `json:"method"`   // 必須是 "VerifyUsername"
		Key      string `json:"Key"`      // 密鑰
		Time     string `json:"Time"`     // 當前時間格式 "yyyyMMddHHmmss"
		Username string `json:"Username"` // 用戶名
	}

	// 響應結構體
	VerifyUsernameRes struct {
		Username   string `xml:"Username"`   // 用戶名
		IsExist    bool   `xml:"IsExist"`    // 用戶是否存在
		ErrorMsgId int    `xml:"ErrorMsgId"` // 錯誤信息 ID
		ErrorMsg   string `xml:"ErrorMsg"`   // 錯誤信息詳細
	}
)

type (
	// 請求結構體
	KickUserReq struct {
		Method   string `json:"method"`   // 必須是 "KickUser"
		Key      string `json:"Key"`      // 密鑰
		Time     string `json:"Time"`     // 當前時間格式 "yyyyMMddHHmmss"
		Username string `json:"Username"` // 用戶名
	}

	// 響應結構體
	KickUserRes struct {
		ErrorMsgId int    `xml:"ErrorMsgId"` // 錯誤信息 ID
		ErrorMsg   string `xml:"ErrorMsg"`   // 錯誤信息詳細
	}
)

type (
	// 請求結構體
	CreditBalanceDVReq struct {
		Method       string  `json:"method"`       // 必須是 "CreditBalanceDV"
		Key          string  `json:"Key"`          // 密鑰
		Time         string  `json:"Time"`         // 當前時間格式 "yyyyMMddHHmmss"
		Username     string  `json:"Username"`     // 用戶名
		OrderId      string  `json:"OrderId"`      // 訂單 ID
		CreditAmount float64 `json:"CreditAmount"` // 信用額，精確到分
		CurrencyType string  `json:"CurrencyType"` // 幣種
	}

	// 響應結構體
	CreditBalanceDVRes struct {
		Username     string  `xml:"Username"`     // 用戶名
		Balance      float64 `xml:"Balance"`      // 現存結餘
		CreditAmount float64 `xml:"CreditAmount"` // 信用額，精確到分
		OrderId      string  `xml:"OrderId"`      // 訂單 ID
		ErrorMsgId   int     `xml:"ErrorMsgId"`   // 錯誤信息 ID
		ErrorMsg     string  `xml:"ErrorMsg"`     // 錯誤信息詳細
	}
)

type (
	GetUserStatusDVReq struct {
		Method   string `json:"method"`   // 必須是 "GetUserStatusDV"
		Key      string `json:"Key"`      // 密鑰
		Time     string `json:"Time"`     // 當前時間格式 "yyyyMMddHHmmss"
		Username string `json:"Username"` // 用戶名
	}

	GetUserStatusDVRes struct {
		IsSuccess    bool    `xml:"IsSuccess"`    // 成功與否
		Username     string  `xml:"Username"`     // 用戶名
		Balance      float64 `xml:"Balance"`      // 現存結餘
		Online       bool    `xml:"Online"`       // 是否在線
		Betted       bool    `xml:"Betted"`       // 是否下注
		BettedAmount float64 `xml:"BettedAmount"` // 總下注額
		MaxBalance   float64 `xml:"MaxBalance"`   // 最大帳戶餘額
		MaxWinning   float64 `xml:"MaxWinning"`   // 最大帳戶贏額
		ErrorMsgId   int     `xml:"ErrorMsgId"`   // 錯誤信息 ID
		ErrorMsg     string  `xml:"ErrorMsg"`     // 錯誤信息詳細
	}
)

type (
	// 請求結構體
	DebitBalanceDVReq struct {
		Method      string  `json:"method"`      // 必須是 "DebitBalanceDV"
		Key         string  `json:"Key"`         // 密鑰
		Time        string  `json:"Time"`        // 當前時間格式 "yyyyMMddHHmmss"
		Username    string  `json:"Username"`    // 用戶名
		OrderId     string  `json:"OrderId"`     // 訂單 ID
		DebitAmount float64 `json:"DebitAmount"` // 借記賬戶額，精確到分
	}

	// 響應結構體
	DebitBalanceDVRes struct {
		Username    string  `xml:"Username"`    // 用戶名
		Balance     float64 `xml:"Balance"`     // 現存結餘
		DebitAmount float64 `xml:"DebitAmount"` // 借記賬戶額，精確到分
		OrderId     string  `xml:"OrderId"`     // 訂單 ID
		ErrorMsgId  int     `xml:"ErrorMsgId"`  // 錯誤信息 ID
		ErrorMsg    string  `xml:"ErrorMsg"`    // 錯誤信息詳細
	}
)

type (
	// 請求結構體
	GetAllBetDetailsForTimeIntervalDVReq struct {
		Method string `json:"method"` // 必須是 "GetAllBetDetailsForTimeIntervalDV"
		Key    string `json:"Key"`    // 密鑰
		Time   string `json:"Time"`   // 當前時間格式 "yyyyMMddHHmmss"
		//Username *string `json:"Username"` // 用戶名 (非必填)
		FromTime string `json:"FromTime"` // 日期詳細信息 "yyyy-MM-dd HH:mm:ss"
		ToTime   string `json:"ToTime"`   // 日期詳細信息 "yyyy-MM-dd HH:mm:ss"
	}

	// 投注信息詳細結構體
	BetDetail struct {
		BetTime       string  `xml:"BetTime"`       // 投注時間
		PayoutTime    string  `xml:"PayoutTime"`    // 結算時間
		Username      string  `xml:"Username"`      // 用戶名
		HostID        int     `xml:"HostID"`        // 桌台ID
		Detail        string  `xml:"Detail"`        // 保留
		GameID        string  `xml:"GameID"`        // 遊戲編號
		Round         int     `xml:"Round"`         // 局
		Set           int     `xml:"Set"`           // 靴
		BetID         int64   `xml:"BetID"`         // 投注編號
		Currency      string  `xml:"Currency"`      // 幣種
		BetAmount     float64 `xml:"BetAmount"`     // 投注金額
		Rolling       float64 `xml:"Rolling"`       // 有效投注額/洗碼量
		ResultAmount  float64 `xml:"ResultAmount"`  // 輸贏金額
		Balance       float64 `xml:"Balance"`       // 投注後的餘額
		GameType      string  `xml:"GameType"`      // 遊戲類型
		BetType       int     `xml:"BetType"`       // 真人遊戲: 不同的投注類型
		BetSource     int     `xml:"BetSource"`     // 投注來源
		TransactionID int64   `xml:"TransactionID"` // 單一錢包下注交易編號
		//GameResult    string  `xml:"GameResult"`    // 遊戲結果
	}

	// 響應結構體
	GetAllBetDetailsForTimeIntervalDVRes struct {
		NumOfRecord   int         `xml:"NumOfRecord"`             // 返回的記錄數
		BetDetailList []BetDetail `xml:"BetDetailList>BetDetail"` // 投注信息詳細
		ErrorMsgId    int         `xml:"ErrorMsgId"`              // 錯誤信息 ID
		ErrorMsg      string      `xml:"ErrorMsg"`                // 錯誤信息詳細
	}
)
