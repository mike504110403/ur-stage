package htpay

import (
	"encoding/json"
	"fmt"
	"wallet_service/internal/apicaller"
	"wallet_service/internal/cachedata"
	"wallet_service/internal/encoder"

	"github.com/google/uuid"
	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

type Config struct {
	HT_CALL_BACK_DOMAIN string
	Ht_Secret           Ht_Secret
}

type Ht_Secret struct {
	OpmhtId string `json:"opmht_id"`
	Domain  string `json:"domain"`
	Key     string `json:"key"`
}

var cfg Config

func Init(initCfg Config) error {
	cfg = initCfg
	cachedataMay := cachedata.ThirdPaySecretMap()
	secretJson, ok := cachedataMay["ht_trc20"]
	if !ok {
		return fmt.Errorf("ht_pay secret not found")
	}
	var ht_sercet Ht_Secret
	err := json.Unmarshal([]byte(secretJson), &ht_sercet)
	if err != nil {
		return err
	}
	cfg.Ht_Secret = ht_sercet
	return nil
}

// PlaceOrderReq : 訂單起單請求
func PlaceOrder(req PlaceOrderReq) (PlaceOrderRes, error) {
	res := PlaceOrderRes{}
	url := cfg.Ht_Secret.Domain + "/api/v2/pay/placeorder"
	req.Currency = "USDT"
	req.OpmhtId = cfg.Ht_Secret.OpmhtId
	req.PayType = "crypto_payasyougo"
	req.NotifyUrl = cfg.HT_CALL_BACK_DOMAIN + "/api/private/v1/wallet/pay/htCallback"
	req.Random = uuid.New().String()

	// 組請求字串
	queryString, err := ToQueryString(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求參數轉換錯誤: %s", err))
		return res, err
	}
	reqSign, err := encoder.Hmac384_Encoder(queryString, []byte(cfg.Ht_Secret.Key))
	if err != nil {
		mlog.Error(fmt.Sprintf("簽名加簽錯誤: %s", err))
		return res, err
	} else {
		req.Sign = reqSign
	}

	reqString, err := ToQueryString(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求參數轉換錯誤: %s", err))
		return res, err
	}
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

	if res.RtCode != 0 {
		return res, fmt.Errorf("響應錯誤: %s", res.Msg)
	} else {
		return res, nil
	}
}

// 取得訂單資訊
func GetInfo(req PayInfoReq) (PayInfoRes, error) {
	res := PayInfoRes{}

	url := cfg.Ht_Secret.Domain + "/api/v2/pay/getInfo"
	req.OpmhtId = cfg.Ht_Secret.OpmhtId
	req.Random = uuid.New().String()

	queryString, err := ToQueryString(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求參數轉換錯誤: %s", err))
		return res, err
	}
	reqSign, err := encoder.Hmac384_Encoder(queryString, []byte(cfg.Ht_Secret.Key))
	if err != nil {
		mlog.Error(fmt.Sprintf("簽名加簽錯誤: %s", err))
		return res, err
	} else {
		req.Sign = reqSign
	}
	req.Sign = reqSign
	reqStr, err := ToQueryString(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求參數轉換錯誤: %s", err))
		return res, err
	}
	url = url + "?" + reqStr
	if err := apicaller.SendGetRequest(url, func(resp *fasthttp.Response) error {
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
	return res, nil
}

// 提款訂單起單
func PayOutPlaceOrder(req PlaceOrderOutReq) (PlaceOrderOutRes, error) {
	res := PlaceOrderOutRes{}
	url := cfg.Ht_Secret.Domain + "/api/v2/payout/placeorder"
	req.AccType = "crypto"
	req.Currency = "USDT"
	req.NotifyUrl = cfg.HT_CALL_BACK_DOMAIN + "/api/private/v1/wallet/payout/htCallback"
	req.OpmhtId = cfg.Ht_Secret.OpmhtId
	req.Random = uuid.New().String()

	// 組請求字串
	queryString, err := ToQueryString(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求參數轉換錯誤: %s", err))
		return res, err
	}
	reqSign, err := encoder.Hmac384_Encoder(queryString, []byte(cfg.Ht_Secret.Key))
	if err != nil {
		mlog.Error(fmt.Sprintf("簽名加簽錯誤: %s", err))
		return res, err
	} else {
		req.Sign = reqSign
	}

	reqString, err := ToQueryString(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求參數轉換錯誤: %s", err))
		return res, err
	}
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

	if res.RtCode != 0 {
		return res, fmt.Errorf("響應錯誤: %s", res.Msg)
	} else {
		return res, nil
	}
}
