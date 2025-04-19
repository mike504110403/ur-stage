package dg_live

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

const Agent_account = "DGTE010809"
const Key = "56f3506d04d249c19ad8ff1ac101b74c"
const AgentFix = "3213"
const AfterFix = "PKZ"
const Domain = "http://apidoc.dg99web.com"
const Url = "/v2/api/signup"

func TestConect(t *testing.T) {
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(Domain + Url)
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sign := GenerateMD5(Agent_account + Key + timestamp)
	req.Header.Set("agent", Agent_account)
	req.Header.Set("sign", sign)
	req.Header.Set("time", timestamp)

	type SingUpReq struct {
		UserName     string  `json:"username"`
		Password     string  `json:"password"`
		CurrencyName string  `json:"currencyName"`
		WinLimit     float64 `json:"winLimit"`
	}

	type SingUpRes struct {
		CodeId int    `json:"codeId"`
		Msg    string `json:"msg"`
	}

	bodyReq := SingUpReq{
		UserName:     "test",
		Password:     GenerateMD5("123456"),
		CurrencyName: "USDT",
		WinLimit:     0,
	}

	body, err := json.Marshal(bodyReq)
	if err != nil {
		t.Error(err.Error())
	}

	// 設置請求正文
	req.Header.SetContentType("application/json")
	req.SetBody(body)

	apiclient := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		ReadTimeout:              30 * time.Second,
		WriteTimeout:             30 * time.Second,
	}
	// 創建一個新的響應對象
	resp := fasthttp.AcquireResponse()

	res := SingUpRes{}
	// 發送請求
	if err := apiclient.Do(req, resp); err != nil {
		t.Error(err.Error())
	} else {
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			t.Error(err.Error())
		} else {
			t.Log()
		}
	}
}

func TestLogIn(t *testing.T) {
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(Domain + "/v2/api/login")
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sign := GenerateMD5(Agent_account + Key + timestamp)
	req.Header.Set("agent", Agent_account)
	req.Header.Set("sign", sign)
	req.Header.Set("time", timestamp)

	type LoginReq struct {
		UserName     string  `json:"username"`
		Password     string  `json:"password"`
		CurrencyName string  `json:"currencyName"`
		WinLimit     float64 `json:"winLimit"`
		Language     string  `json:"language"`
	}

	type LoginRes struct {
		CodeId  int    `json:"codeId"`
		Msg     string `json:"msg"`
		Token   string `json:"token"`
		Domains string `json:"domains"`
		List    []any  `json:"list"`
	}

	bodyReq := LoginReq{
		UserName:     "test",
		Password:     GenerateMD5("123456"),
		CurrencyName: "USDT",
		WinLimit:     0,
		Language:     "cn",
	}

	body, err := json.Marshal(bodyReq)
	if err != nil {
		t.Error(err.Error())
	}

	// 設置請求正文
	req.Header.SetContentType("application/json")
	req.SetBody(body)

	apiclient := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		ReadTimeout:              30 * time.Second,
		WriteTimeout:             30 * time.Second,
	}
	// 創建一個新的響應對象
	resp := fasthttp.AcquireResponse()

	res := LoginRes{}
	// 發送請求
	if err := apiclient.Do(req, resp); err != nil {
		t.Error(err.Error())
	} else {
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			t.Error(err.Error())
		} else {
			fmt.Println(res.List[0])
			t.Log()
		}
	}
}

func TestFree(t *testing.T) {
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(Domain + "/v2/api/free")
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sign := GenerateMD5(Agent_account + Key + timestamp)
	req.Header.Set("agent", Agent_account)
	req.Header.Set("sign", sign)
	req.Header.Set("time", timestamp)

	type FreeReq struct {
		Language string `json:"language"`
	}

	type FreeRes struct {
		Url string `json:"url"`
	}

	bodyReq := FreeReq{
		Language: "cn",
	}

	body, err := json.Marshal(bodyReq)
	if err != nil {
		t.Error(err.Error())
	}

	// 設置請求正文
	req.Header.SetContentType("application/json")
	req.SetBody(body)

	apiclient := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		ReadTimeout:              30 * time.Second,
		WriteTimeout:             30 * time.Second,
	}
	// 創建一個新的響應對象
	resp := fasthttp.AcquireResponse()

	res := FreeRes{}
	// 發送請求
	if err := apiclient.Do(req, resp); err != nil {
		t.Error(err.Error())
	} else {
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			t.Error(err.Error())
		} else {
			t.Log()
		}
	}
}

func TestOnline(t *testing.T) {
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(Domain + "/v2/api/free")
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sign := GenerateMD5(Agent_account + Key + timestamp)
	req.Header.Set("agent", Agent_account)
	req.Header.Set("sign", sign)
	req.Header.Set("time", timestamp)

	type OnlineRes struct {
		UserName string  `json:"username"`
		NickName string  `json:"nickname"`
		Currency string  `json:"currencyName"`
		Ip       string  `json:"ip"`
		Device   string  `json:"device"`
		Login    string  `json:"login"`
		MemberId int     `json:"memberId"`
		Balance  float64 `json:"balance"`
	}

	// 設置請求正文
	req.Header.SetContentType("application/json")

	apiclient := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		ReadTimeout:              30 * time.Second,
		WriteTimeout:             30 * time.Second,
	}
	// 創建一個新的響應對象
	resp := fasthttp.AcquireResponse()

	res := OnlineRes{}
	// 發送請求
	if err := apiclient.Do(req, resp); err != nil {
		t.Error(err.Error())
	} else {
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			t.Error(err.Error())
		} else {
			t.Log()
		}
	}
}

func TestKickOut(t *testing.T) {
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(Domain + "/v2/api/free")
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	sign := GenerateMD5(Agent_account + Key + timestamp)
	req.Header.Set("agent", Agent_account)
	req.Header.Set("sign", sign)
	req.Header.Set("time", timestamp)

	type FreeReq struct {
		Language string `json:"language"`
	}

	type FreeRes struct {
		Url string `json:"url"`
	}

	bodyReq := FreeReq{
		Language: "cn",
	}

	body, err := json.Marshal(bodyReq)
	if err != nil {
		t.Error(err.Error())
	}

	// 設置請求正文
	req.Header.SetContentType("application/json")
	req.SetBody(body)

	apiclient := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		ReadTimeout:              30 * time.Second,
		WriteTimeout:             30 * time.Second,
	}
	// 創建一個新的響應對象
	resp := fasthttp.AcquireResponse()

	res := FreeRes{}
	// 發送請求
	if err := apiclient.Do(req, resp); err != nil {
		t.Error(err.Error())
	} else {
		if err := json.Unmarshal(resp.Body(), &res); err != nil {
			t.Error(err.Error())
		} else {
			t.Log()
		}
	}
}
