package order

import (
	"encoding/json"
	"fmt"
	"testing"

	htpay "wallet_service/instance/ht_pay"
	"wallet_service/internal/apicaller"
	"wallet_service/internal/encoder"

	"github.com/google/uuid"
	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

const Key = "rU7cJlu29UCTDuokVMcPusRbfz5aA943zhtOJ8MJslOqkFvthLpRSKtlqfsctfHO"

func TestHtCallbackHandler(t *testing.T) {
	req := PayCallbackReq{
		Accno:      "123456",
		Attach:     "test",
		Currency:   "USDT",
		MhtOrderNo: "2a373638-ac6b-40aa-9aa4-3a06e4845992",
		Note:       "test",
		PaidAmount: 1,
		PayType:    "test",
		PfOrderNo:  "RVTJ240905165232",
		Random:     uuid.New().String(),
		Status:     1,
	}
	enStr, err := htpay.ToQueryString(req)
	if err != nil {
		t.Error(err)
	}

	ened, err := encoder.Hmac384_Encoder(enStr, []byte(Key))
	if err != nil {
		t.Error(err)
	}

	req.Sign = ened
	url := "http://localhost:3001/v1/order/htCallback"

	reqString, err := htpay.ToQueryString(req)
	if err != nil {
		t.Error(err)
	}
	res := ""

	// 發送請求
	if err = apicaller.SendPostRequest(url, reqString, func(resp *fasthttp.Response) error {
		if resp.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("響應狀態碼錯誤: %d", resp.StatusCode())
		}
		err := json.Unmarshal(resp.Body(), &res)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		mlog.Error(fmt.Sprintf("發送請求錯誤: %s", err))
	}

}
