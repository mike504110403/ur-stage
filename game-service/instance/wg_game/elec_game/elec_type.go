package elec_game

type MemberGameAccount struct {
	Username string `db:"username"`
	Nickname string `db:"nickname"`
	Password string `db:"password"`
}

type (
	CreateUserReq struct {
		ApiId      string  `json:"api_id"`     // 代理商編號
		Prefix     *string `json:"prefix"`     // 代理商前綴
		Username   string  `json:"username"`   // 會員名稱
		Nickname   *string `json:"nickname"`   // 會員暱稱
		Handicaps  *string `json:"handicaps"`  // 盤口
		Proxyname  string  `json:"proxyname"`  // 代理帳號
		Experience string  `json:"experience"` // 體驗帳號
		Sign       string  `json:"sign"`       // 簽名
	}
	CreateUserRes struct {
		ErrorCode string  `json:"error_code"` // 回傳代碼
		Username  *string `json:"username"`   // 會員名稱
	}
)

type (
	CheckUserReq struct {
		ApiId    string `json:"api_id"`   // 代理商編號
		Username string `json:"username"` // 會員名稱
		Sign     string `json:"sign"`     // 簽名
	}
	CheckUserRes struct {
		ErrorCode string `json:"error_code"` // 回傳代碼
	}
)

type (
	EditUserReq struct {
		ApiId       string  `json:"api_id"`       // 代理商編號
		Username    string  `json:"username"`     // 會員名稱
		Nickname    *string `json:"nickname"`     // 暱稱 (已失效)
		Enabled     *int    `json:"enabled"`      // 彩票登入狀態
		EnabledLive *int    `json:"enabled_live"` // 真人登入狀態
		EnabledElec *int    `json:"enabled_elec"` // 電子登入狀態
		Sign        string  `json:"sign"`         // 簽名
	}
	EditUserRes struct {
		ErrorCode string `json:"error_code"` // 回傳代碼
	}
)

type (
	TransferCheckReq struct {
		ApiId      string `json:"api_id"`     // 代理商編號
		Username   string `json:"username"`   // 會員名稱
		TransferId string `json:"transferid"` // 交易單號
		Sign       string `json:"sign"`       // 簽名
	}
	TransferCheckRes struct {
		ErrorCode string `json:"error_code"` // 回傳代碼
	}
)

type (
	ForwardGameReq struct {
		ApiId      string  `json:"api_id"`     // 代理商編號
		Username   string  `json:"username"`   // 會員名稱
		Nickname   *string `json:"nickname"`   // 會員暱稱 (已失效)
		Experience *string `json:"experience"` // 體驗帳號 (已失效)
		Mobile     *string `json:"mobile"`     // 手機裝置
		Locale     *string `json:"locale"`     // 語系
		Sign       string  `json:"sign"`       // 簽名
	}
	ForwardGameRes struct {
		ErrorCode string `json:"error_code"` // 回傳代碼
		URL       string `json:"url"`        // 進入遊戲連結
	}
)

type (
	TransferInReq struct {
		ApiId      string  `json:"api_id"`     // 代理商編號
		Username   string  `json:"username"`   // 會員名稱
		Bpoints    string  `json:"bpoints"`    // 轉前點數
		Points     string  `json:"points"`     // 轉入點數
		Apoints    string  `json:"apoints"`    // 轉後點數
		TransferId *string `json:"transferid"` // 自定交易單號 (非必填)
		Sign       string  `json:"sign"`       // 簽名
	}
	TransferInRes struct {
		ErrorCode  string   `json:"error_code"` // 回傳代碼
		TransferId *string  `json:"transferid"` // 交易序號
		Points     *float64 `json:"points"`     // 交易後點數
	}
)

type (
	TransferOutReq struct {
		ApiId      string  `json:"api_id"`     // 代理商編號
		Username   string  `json:"username"`   // 會員名稱
		Bpoints    string  `json:"bpoints"`    // 轉前點數
		Points     string  `json:"points"`     // 轉出點數
		Apoints    string  `json:"apoints"`    // 轉後點數
		TransferId *string `json:"transferid"` // 自定交易單號 (非必填)
		Sign       string  `json:"sign"`       // 簽名
	}
	TransferOutRes struct {
		ErrorCode  string   `json:"error_code"` // 回傳代碼
		TransferId *string  `json:"transferid"` // 交易單號
		Points     *float64 `json:"points"`     // 交易後點數
	}
)

type (
	PointUserReq struct {
		ApiId    string `json:"api_id"`   // 代理商編號
		Username string `json:"username"` // 會員名稱
		Sign     string `json:"sign"`     // 簽名
	}
	PointUserRes struct {
		ErrorCode string   `json:"error_code"` // 回傳代碼
		Points    *float64 `json:"points"`     // 現有點數
	}
)

type (
	BuyListGetApiReq struct {
		ApiId     string  `json:"api_id"`    // 代理商編號
		Date      *string `json:"date"`      // 時間 (YYYY-MM-DD)
		Username  *string `json:"username"`  // 會員名稱
		Proxyname string  `json:"proxyname"` // 代理帳號
		Type      *int    `json:"type"`      // 取得類型
		PriType   string  `json:"pri_type"`  // 注單狀態
		Locale    *string `json:"locale"`    // 選擇語系
		Sign      string  `json:"sign"`      // 簽名
	}
	BuyListGetApiRes struct {
		ErrorCode string        `json:"error_code"` // 回傳代碼
		Pages     *int          `json:"pages"`      // 總頁數
		Data      []BettingData `json:"data"`       // 注單資料
	}
	BettingData struct {
		BuyId     string `json:"buyid"`      // 注單編號
		Username  string `json:"username"`   // 會員名稱
		PlayKey   string `json:"playkey"`    // 遊戲代碼
		Period    string `json:"period"`     // 遊戲局號
		Number    string `json:"number"`     // 投注內容
		Money     string `json:"money"`      // 下注金額
		SelfPoint string `json:"selfpoint"`  // 自身返點
		PriMoney  string `json:"pri_money"`  // 中獎金額
		MoneyOk   string `json:"money_ok"`   // 有效投注
		WinLose   string `json:"winlose"`    // 輸贏
		Status    string `json:"status"`     // 注單狀態
		PrizeDate string `json:"prize_date"` // 歸帳日
		CreatedAt string `json:"created_at"` // 下注時間
	}
)

type (
	ProxyWinloseGetReq struct {
		ApiId     string  `json:"api_id"`    // 代理商編號
		Date      string  `json:"date"`      // 時間 (YYYY-MM-DD)
		PriType   string  `json:"pri_type"`  // 注單狀態
		Proxyname *string `json:"proxyname"` // 代理帳號
		Sign      string  `json:"sign"`      // 簽名
	}
	ProxyWinloseGetRes struct {
		ErrorCode string             `json:"error_code"` // 回傳代碼
		Data      []ProxyWinloseData `json:"data"`       // 回傳資料
	}
	ProxyWinloseData struct {
		Proxyname  string  `json:"proxyname"`  // 代理帳號
		Currency   string  `json:"currency"`   // 幣種
		TotalCount int     `json:"totalcount"` // 投注單數
		Money      int     `json:"money"`      // 投注金額
		Wait       int     `json:"wait"`       // 未派彩金額
		WinLose    float64 `json:"winlose"`    // 輸贏金額
		MoneyOk    float64 `json:"money_ok"`   // 有效投注
	}
)

type (
	BuySingleGetApiReq struct {
		ApiId     string  `json:"api_id"`    // 代理商編號
		Proxyname *string `json:"proxyname"` // 代理帳號
		BuyId     string  `json:"buyid"`     // 單號
		Sign      string  `json:"sign"`      // 簽名
	}
	BuySingleGetApiRes struct {
		ErrorCode string `json:"error_code"` // 回傳代碼
		Data      string `json:"data"`       // 回傳資料 (HTML 結構)
	}
)

type (
	KickUserReq struct {
		ApiId    string `json:"api_id"`   // 代理商編號
		Username string `json:"username"` // 會員帳號
		Sign     string `json:"sign"`     // 簽名
	}
	KickUserRes struct {
		ErrorCode string `json:"error_code"` // 回傳代碼
	}
)
