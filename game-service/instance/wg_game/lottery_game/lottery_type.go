package lottery_game

type MemberGameAccount struct {
	Username string `db:"username"`
	Nickname string `db:"nickname"`
	Password string `db:"password"`
}

type (
	CreateUserApiReq struct {
		ApiId      string  `json:"api_id"`     // 代理商編號
		Prefix     *string `json:"prefix"`     // 代理商前綴
		Username   string  `json:"username"`   // 會員名稱
		Nickname   *string `json:"nickname"`   // 會員暱稱
		Handicaps  *string `json:"handicaps"`  // 盤口
		Proxyname  string  `json:"proxyname"`  // 代理帳號
		Experience string  `json:"experience"` // 體驗帳號
		Sign       string  `json:"sign"`       // 簽名
	}
	CreateUserApiRes struct {
		ErrorCode string `json:"error_code"` // 回傳代碼
		Username  string `json:"username"`   // 會員名稱
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
		Handicaps  *string `json:"handicaps"`  // 盤口 (已失效)
		VipLot     *int    `json:"vip_lot"`    // 會員vip等級
		DepositNum *int    `json:"depositnum"` // 最近1個月的存款次數
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
		BPoints    string  `json:"bpoints"`    // 轉前點數
		Points     string  `json:"points"`     // 轉入點數
		APoints    string  `json:"apoints"`    // 轉後點數
		TransferId *string `json:"transferid"` // 自定交易單號
		Sign       string  `json:"sign"`       // 簽名
	}
	TransferInRes struct {
		ErrorCode   string   `json:"error_code"` // 回傳代碼
		TransferId  *string  `json:"transferid"` // 交易序號
		PointsAfter *float64 `json:"points"`     // 交易後點數
	}
)

type (
	TransferOutReq struct {
		ApiId      string  `json:"api_id"`     // 代理商編號
		Username   string  `json:"username"`   // 會員名稱
		BPoints    string  `json:"bpoints"`    // 轉前點數
		Points     string  `json:"points"`     // 轉出點數
		APoints    string  `json:"apoints"`    // 轉後點數
		TransferId *string `json:"transferid"` // 自定交易單號
		Sign       string  `json:"sign"`       // 簽名
	}
	TransferOutRes struct {
		ErrorCode   string   `json:"error_code"` // 回傳代碼
		TransferId  *string  `json:"transferid"` // 交易單號
		PointsAfter *float64 `json:"points"`     // 交易後點數
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
		ProxyName string  `json:"proxyname"` // 代理帳號
		Type      *int    `json:"type"`      // 取得類型
		PriType   string  `json:"pri_type"`  // 注單狀態
		Page      *int    `json:"page"`      // 分頁
		Locale    *string `json:"locale"`    // 選擇語系
		PlayKey   *string `json:"playkey"`   // 彩種代碼
		Sign      string  `json:"sign"`      // 簽名
	}
	BuyListGetApiRes struct {
		ErrorCode string        `json:"error_code"` // 回傳代碼
		Pages     *int          `json:"pages"`      // 總頁數
		Data      []BuyListData `json:"data"`       // 注單資料
	}
	BuyListData struct {
		BuyId        string   `json:"buyid"`         // 注單編號
		Username     string   `json:"username"`      // 會員名稱
		Code         string   `json:"code"`          // 玩法類別
		PlayKey      string   `json:"playkey"`       // 彩票種類
		ListId       string   `json:"list_id"`       // 玩法代碼
		Period       string   `json:"period"`        // 期號
		Number       string   `json:"number"`        // 投注號碼
		Nums         string   `json:"nums"`          // 注數
		Money        string   `json:"money"`         // 下注金額
		MoneyOk      string   `json:"money_ok"`      // 有效量_投注
		PriMoney     *string  `json:"pri_money"`     // 中獎金額
		ZBuyRate     string   `json:"z_buy_rate"`    // 投注時賠率
		PriNumber    string   `json:"pri_number"`    // 中獎號碼
		Modes        string   `json:"modes"`         // 下注模式
		ZNumber      string   `json:"z_number"`      // 追單單碼
		Status       string   `json:"status"`        // 訂單狀態
		CreatedAt    string   `json:"created_at"`    // 下注時間
		PrizeTime    *string  `json:"prize_time"`    // 派彩時間
		PrizeDate    string   `json:"prize_date"`    // 報表時間
		Handicaps    string   `json:"handicaps"`     // 盤口值
		ErrorMsg     string   `json:"errormsg"`      // 錯誤訊息
		Currency     string   `json:"currency"`      // 幣值
		CurrencyDiff *float64 `json:"currency_diff"` // 匯差
		SuperMoney   string   `json:"super_money"`   // 超級彩金
		SuperStatus  string   `json:"super_status"`  // 超級彩金類別
		AfterPoints  string   `json:"afterpoints"`   // 下注後餘額
	}
)

type (
	GiftListGetApiv2Req struct {
		ApiId     string  `json:"api_id"`    // 代理商編號
		Date      *string `json:"date"`      // 時間 (YYYY-MM-DD)
		ProxyName string  `json:"proxyname"` // 代理帳號
		Locale    *string `json:"locale"`    // 選擇語系
		Sign      string  `json:"sign"`      // 簽名
	}
	GiftListGetApiv2Res struct {
		ErrorCode string         `json:"error_code"` // 回傳代碼
		Data      []GiftListData `json:"data"`       // 打賞資料
	}
	GiftListData struct {
		DeptId      string `json:"dept_id"`     // 部門代碼
		BuyId       string `json:"buyid"`       // 注單編號
		Time        string `json:"time"`        // 打賞時間
		Username    string `json:"username"`    // 會員名稱
		PlayKey     string `json:"playkey"`     // 彩票種類
		FullName    string `json:"fullname"`    // 彩票名稱
		ListId      string `json:"list_id"`     // 項目代碼
		ListName    string `json:"listname"`    // 項目名稱
		CamgirlId   string `json:"camgirl_id"`  // 主播代碼
		GirlName    string `json:"girlname"`    // 主播名稱
		Money       string `json:"money"`       // 打賞金額
		AfterPoints string `json:"afterpoints"` // 打賞後餘額
	}
)

type (
	ProxyWinloseGetReq struct {
		ApiId     string  `json:"api_id"`    // 代理商編號
		Date      string  `json:"date"`      // 時間 (YYYY-MM-DD)
		PriType   string  `json:"pri_type"`  // 注單狀態
		ProxyName *string `json:"proxyname"` // 代理帳號
		Sign      string  `json:"sign"`      // 簽名
	}
	ProxyWinloseGetRes struct {
		ErrorCode string             `json:"error_code"` // 回傳代碼
		Data      []ProxyWinloseData `json:"data"`       // 回傳資料
	}
	ProxyWinloseData struct {
		ProxyName  string `json:"proxyname"`  // 代理帳號
		Currency   string `json:"currency"`   // 幣種
		TotalCount string `json:"totalcount"` // 投注單數
		Money      string `json:"money"`      // 投注金額
		Wait       string `json:"wait"`       // 未派彩金額
		WinLose    string `json:"winlose"`    // 輸贏金額
		MoneyOk    string `json:"money_ok"`   // 有效量_投注
	}
)

type (
	UserWinloseGetReq struct {
		ApiId     string  `json:"api_id"`    // 代理商編號
		Date      string  `json:"date"`      // 時間 (YYYY-MM-DD)
		PriType   string  `json:"pri_type"`  // 注單狀態
		ProxyName *string `json:"proxyname"` // 代理帳號
		Sign      string  `json:"sign"`      // 簽名
	}
	UserWinloseGetRes struct {
		ErrorCode string            `json:"error_code"` // 回傳代碼
		Data      []UserWinloseData `json:"data"`       // 回傳資料
	}
	UserWinloseData struct {
		PlayKey    string `json:"playkey"`    // 彩種代碼
		Username   string `json:"username"`   // 會員帳號
		Currency   string `json:"currency"`   // 幣種
		TotalCount string `json:"totalcount"` // 投注單數
		Money      string `json:"money"`      // 投注金額
		Wait       string `json:"wait"`       // 未派彩金額
		WinLose    string `json:"winlose"`    // 輸贏金額
		MoneyOk    string `json:"money_ok"`   // 有效量
	}
)

type (
	BuySingleGetApiReq struct {
		ApiId     string `json:"api_id"`    // 代理商編號
		ProxyName string `json:"proxyname"` // 代理帳號
		BuyId     string `json:"buyid"`     // 單號
		Sign      string `json:"sign"`      // 簽名
	}
	BuySingleGetApiRes struct {
		ErrorCode string `json:"error_code"` // 回傳代碼
		Data      string `json:"data"`       // 回傳資料 (直接回應開獎號 html 結構 內含css link file)
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
