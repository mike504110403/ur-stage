package order

// 訂單起單請求
type (
	OrderInitReq struct {
		ItemType  string  `json:"item_type" validate:"required"`
		ItemCount int     `json:"item_count" validate:"required"`
		Amount    float64 `json:"amount" validate:"required"`
		PayType   string  `json:"pay_type" validate:"required"`
	}
	OrderInitRes struct {
		Url string `json:"url"`
	}
)

// 第三方訂單起單請求
type (
	thirdPayInitAPIReq struct {
		ThirdPayType string  `json:"third_pay_type"`
		Amount       float64 `json:"amount"`
		TradeNo      string  `json:"trade_no"`
	}

	thirdPayInitAPIRes struct {
		TradeNo         string `json:"tradeNo"`
		ThirdPayTradeNo string `json:"third_pay_tradeNo"`
		Url             string `json:"url"`
		Expiration      string `json:"expiration"`
	}
)

// 收款請求
const PayCallbackReqContentType = "application/x-www-form-urlencoded"

type PayCallbackReqStatus int

const (
	PayCallbackReqStatusSuccess PayCallbackReqStatus = 1
	PayCallbackReqStatusFail    PayCallbackReqStatus = 3
)

// 代收回調
type (
	PayCallbackReq struct {
		Accno      string  `json:"accno"`      // 付款人帳號
		Attach     string  `json:"attach"`     // 附加内容
		Currency   string  `json:"currency"`   // 貨幣名稱
		MhtOrderNo string  `json:"mhtorderno"` // 我方訂單號
		Note       string  `json:"note"`       // 備註
		PaidAmount float64 `json:"paidamount"` // 金额
		PayType    string  `json:"paytype"`    // 支付類型 這邊應該是 usdt
		PfOrderNo  string  `json:"pforderno"`  // 三方訂單號
		Random     string  `json:"random"`     // 隨機字串
		Status     int     `json:"status"`     // 订单状态，1: 成功；3: 失败
		Sign       string  `json:"sign"`       // 签名
	}
)

// 轉點設定
type (
	TransectionSetDto struct {
		MemberId           int     `json:"member_id"`
		Amount             float64 `json:"amount"`
		TransectionSrcType int     `json:"transection_src_type"`
		TransectionRelate  int     `json:"transection_relate"`
	}
)

type (
	AdminPaymentReq struct {
		UserName string  `json:"username" validate:"required"`
		Amount   float64 `json:"amount" validate:"required"`
	}
)
