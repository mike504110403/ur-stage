package htpay

// 支付訂單 Pay
// /api/{version}/pay/placeOrder [POST]
const (
	PlaceOrderReqContentType = "application/x-www-form-urlencoded"
	PlaceOrderResContentType = "application/json;charset=utf-8"
)

type (
	PlaceOrderReq struct {
		Accno      string  `json:"accno"`                                                   // 付款人帳號
		Amount     float64 `json:"amount" validate:"required"`                              // 付款金額
		Attach     string  `json:"attach"`                                                  // 附加資訊 - 用於回傳時辨識
		ClientIp   string  `json:"clientip" validate:"required"`                            //TODO: 我方IP?/付款人IP?
		Currency   string  `json:"currency" validate:"required" default:"USDT"`             // 貨幣
		Extra      string  `json:"extra"`                                                   // 額外資訊
		MhtOrderNo string  `json:"mhtorderno" validate:"required"`                          // 我方訂單號
		MhtUserId  string  `json:"mhtuserid" validate:"required"`                           // 我方會員編號
		NotifyUrl  string  `json:"notifyurl"`                                               // 入款通知URL
		OpmhtId    string  `json:"opmhtid" validate:"required"`                             // 商戶ID
		PayerName  string  `json:"payername"`                                               // 付款人姓名
		PayerPhone string  `json:"payerphone"`                                              // 付款人電話
		PayType    string  `json:"paytype" validate:"required" default:"crypto_payasyougo"` // 付款方式
		Random     string  `json:"random" validate:"required"`                              // 隨機字串
		ReturnUrl  string  `json:"returnurl"`                                               // 付款完成後返回URL
		BankCode   string  `json:"bankcode"`                                                // 銀行代碼
		Sign       string  `json:"sign" validate:"required"`                                // 簽名
	}
	PlaceOrderRes struct {
		RtCode int    `json:"rtCode"` // 返回码，0：成功；其他：失败
		Msg    string `json:"msg"`    // 返回讯息
		Result *struct {
			PfOrderNo     string `json:"pforderno"` // 三方訂單號
			PayUrl        string `json:"payurl"`    // 支付URL
			PaymentDetail *struct {
				QrCodeRaw string `json:"qrcoderaw"` // 二维码原始内容
				Bank      string `json:"bank"`      // 支付银行
				Branch    string `json:"branch"`    // 支付分行
				Account   string `json:"account"`   // 支付帳號
				Holder    string `json:"holder"`    // 支付戶名
			} `json:"paymentdetail"` // 支付資訊
		} `json:"result"`
	}
)

// 支付訂單狀態 Pay Info
// /api/{version}/pay/getInfo [GET]
const PayInfoResContentType = "application/json;charset=utf-8"

type (
	PayInfoResRtCode       int
	PayInfoResStatus       int
	PayInfoResNotifyStatus int
)

const (
	PayInfoResRtCodeSuccess          PayInfoResRtCode       = 0
	PayInfoResRtCodeOrderNotExist    PayInfoResRtCode       = 202
	PayInfoResStatusUnpaid           PayInfoResStatus       = 0
	PayInfoResStatusPaid             PayInfoResStatus       = 1
	PayInfoResStatusTimeout          PayInfoResStatus       = 2
	PayInfoResStatusFailed           PayInfoResStatus       = 3
	PayInfoResNotifyStatusUnnotified PayInfoResNotifyStatus = 0
	PayInfoResNotifyStatusSuccess    PayInfoResNotifyStatus = 1
	PayInfoResNotifyStatusFailed     PayInfoResNotifyStatus = 2
)

type (
	PayInfoReq struct {
		MhtOrderNo string `json:"mhtorderno" validate:"required"` // 我方訂單號
		OpmhtId    string `json:"opmhtid" validate:"required"`    // 商戶ID
		Random     string `json:"random" validate:"required"`     // 隨機字串
		Sign       string `json:"sign" validate:"required"`       // 簽名
	}
	PayInfoRes struct {
		RtCode int    `json:"rtCode"` // 返回码，0：成功；202：订单不存在
		Msg    string `json:"msg"`    // 返回讯息
		Result *struct {
			PfOrderNo      string  `json:"pforderno"`      // 三方訂單號
			OrderAmount    float64 `json:"orderamount"`    // 訂單金額
			PaidAmount     float64 `json:"paidamount"`     // 支付金額
			Currency       string  `json:"currency"`       // 貨幣
			PayerName      string  `json:"payername"`      // 付款人姓名
			PayType        string  `json:"paytype"`        // 支付類型
			Accno          string  `json:"accno"`          // 支付人帳號
			Attach         string  `json:"attach"`         // 附加内容
			Note           string  `json:"note"`           // 備註
			OrderTime      string  `json:"ordertime"`      // 起單時間
			Status         int     `json:"status"`         // 訂單狀態，0：未付款，1：已付款，2：超时，3：付款失败
			SettleTime     string  `json:"settletime"`     // 訂單完成時間
			NotifyUrl      string  `json:"notifyurl"`      // 回調 URL
			NotifyStatus   int     `json:"notifystatus"`   // 回調狀態，0：未回調，1：回調成功，2：回調失败
			LastNotifyTime string  `json:"lastnotifytime"` // 最近一次回調时间
			Reference      string  `json:"reference"`      // 参考资讯
			FromIp         string  `json:"fromip"`         // 請求端 IP
		} `json:"result"`
	}
)

// 代付訂單 Payout
// /api/{version}/payout/placeOrder [POST]
const PlaceOrderOutReqContentType = "application/x-www-form-urlencoded"
const PlaceOrderOutResContentType = "application/json;charset=utf-8"

type (
	PlaceOrderOutReq struct {
		AccNo      string  `json:"accno"`                    // 收款人帳號
		AccType    string  `json:"acctype" default:"crypto"` // 帳戶類型
		Amount     float64 `json:"amount"`                   // 金额
		Currency   string  `json:"currency" default:"USDT"`  // 幣別
		MhtOrderNo string  `json:"mhtorderno"`               // 我方訂單號
		NotifyUrl  string  `json:"notifyurl"`                // 回調 URL
		OpmhtId    string  `json:"opmhtid"`                  //
		Random     string  `json:"random"`                   // 隨機字串
		Remark     string  `json:"remark"`                   // 備註
		Extra      string  `json:"extra"`                    // 擴充欄位
		Sign       string  `json:"sign"`                     // 簽名
		// 以下非 USDT 代付必填
		// AccCityName string `json:"acccityname"`             // 代付城市名称
		// AccName     string `json:"accname"`                 // 收款人帳戶名
		// AccProvince string `json:"accprovince"`             // 开户行所在省份名称
		// BankBranch  string `json:"bankbranch"`              // 分行名称或代码，日本网银转账必填且提供分行代码
		// BankCode    string `json:"bankcode"`                // 银行代码，请参阅 Bank List
		// PayerPhone  string `json:"payerphone"`              // 收款方手机号
	}
	PlaceOrderOutRes struct {
		RtCode int    `json:"rtCode"` // 返回码，0：成功；其他：失败
		Msg    string `json:"msg"`    // 返回讯息
		Result *struct {
			PfOrderNo    string  `json:"pforderno"`    // 三方訂單號
			Afterbalance float64 `json:"afterbalance"` // 支付後剩餘可代付金額
		} `json:"result"`
	}
)

// 代付訂單狀態 Payout Info
// /api/{version}/payout/getInfo [GET]
const PayoutInfoResContentType = "application/json;charset=utf-8"

type (
	PayoutInfoResRtCode       int
	PayoutInfoResStatus       int
	PayoutInfoResNotifyStatus int
)

const (
	PayoutInfoResRtCodeSuccess          PayoutInfoResRtCode       = 0
	PayoutInfoResRtCodeOrderNotExist    PayoutInfoResRtCode       = 202
	PayoutInfoResStatusProcessing       PayoutInfoResStatus       = 0
	PayoutInfoResStatusSuccess          PayoutInfoResStatus       = 1
	PayoutInfoResStatusFailed           PayoutInfoResStatus       = 2
	PayoutInfoResNotifyStatusUnnotified PayoutInfoResNotifyStatus = 0
	PayoutInfoResNotifyStatusSuccess    PayoutInfoResNotifyStatus = 1
	PayoutInfoResNotifyStatusFailed     PayoutInfoResNotifyStatus = 2
)

type (
	PayoutInfoReq struct {
		MhtOrderNo string `json:"mhtorderno" validate:"required"` // 我方訂單號
		OpmhtId    string `json:"opmhtid" validate:"required"`    // 商戶ID
		Random     string `json:"random" validate:"required"`     // 隨機字串
		Sign       string `json:"sign" validate:"required"`       // 簽名
	}
	PayoutInfoRes struct {
		RtCode int    `json:"rtCode"` // 返回码：0 表示成功；其他值表示失败
		Msg    string `json:"msg"`    // 返回讯息
		Result *struct {
			PforOrderNo    string  `json:"pforderno"`      // 三方訂單號
			OrderAmount    float64 `json:"orderamount"`    // 訂單金額
			PaidAmount     float64 `json:"paidamount"`     // 支付金額
			Currency       string  `json:"currency"`       // 貨幣
			AccType        string  `json:"acctype"`        // 帳戶類型
			Remark         string  `json:"remark"`         // 備註
			OrderTime      string  `json:"ordertime"`      // 訂單創建時間
			Status         int     `json:"status"`         // 訂單狀態：0 處理中，1 成功，2 失敗
			SettleTime     string  `json:"settletime"`     // 訂單完成時間
			NotifyUrl      string  `json:"notifyurl"`      // 回調 url
			NotifyStatus   int     `json:"notifystatus"`   // 回調狀態：0 未回調，1 回調成功，2 回調失敗
			LastNotifyTime string  `json:"lastnotifytime"` // 最近一次回調时间
			BeforeBalance  float64 `json:"beforebalance"`  // 支付前剩余可代付金额
			AfterBalance   float64 `json:"afterbalance"`   // 支付後剩余可代付金额
			FromIp         string  `json:"fromip"`         // 客户端 IP
			// 以下非 USDT 代付必填
			// BankCode       string  `json:"bankcode"`       // 銀行代碼
			// AccProvince    string  `json:"accprovince"`    // 开户行所在省份名称
			// AccCityName    string  `json:"acccityname"`    // 开户行所在城市名称
			// AccName        string  `json:"accname"`        // 收款人帐户名
			// AccNo          string  `json:"accno"`          // 收款人帐号/地址
		} `json:"result"`
	}
)
