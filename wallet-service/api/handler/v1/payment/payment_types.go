package payment

type (
	withdrawReq struct {
		MemberId     int     `json:"member_id"`
		TradeNo      string  `json:"trade_no"`
		AccNo        string  `json:"acc_no" validate:"required"`
		WithdrawType string  `json:"withdraw_type" validate:"required"`
		Amount       float64 `json:"amount" validate:"required"`
		Fee          float64 `json:"fee"`
	}
)

// thirdWithDrawAPIReq 第三方提款API請求
type (
	thirdWithDrawAPIReq struct {
		AccNo        string  `json:"acc_no"`
		TradeNo      string  `json:"trade_no"`
		WithdrawType string  `json:"withdraw_type"`
		Amount       float64 `json:"amount"`
	}
	thirdWithDrawAPIRes struct {
		TradeNo         string `json:"trade_no"`
		ThirdPayTradeNo string `json:"third_pay_trade_no"`
		Result          struct {
			Success bool `json:"success"`
		} `json:"result"`
	}
)

// 代付回調 - 給第三方支付系統用
const PayOutCallBackReqContentType = "application/x-www-form-urlencoded"

type PayoutCallBackReqStatus int

const (
	PayoutCallBackReqStatusSuccess PayoutCallBackReqStatus = 0
	PayoutCallBackReqStatusFail    PayoutCallBackReqStatus = 1
)

type (
	PayoutResult struct {
		AfterBalance float64 `json:"afterbalance"` // 剩餘餘額
		Amount       string  `json:"amount"`       // 代付金額
		Currency     string  `json:"currency"`     // 貨幣
		MhtOrderNo   string  `json:"mhtorderno"`   // 我方訂單號
		Note         string  `json:"note"`         // 备注
		PayoutTime   string  `json:"payouttime"`   // 代付完成時間
		PfOrderNo    string  `json:"pforderno"`    // 三方訂單號
		Random       string  `json:"random"`       // 隨機碼
		ResultCode   int     `json:"resultcode"`   // 訂單狀態 0: 成功；1: 失敗
		Sign         string  `json:"sign"`         // 簽名
	}
)
