package apicaller

import (
	"github.com/valyala/fasthttp"
)

type ResponseHandler func(*fasthttp.Response) error

// SendPostRequest 發送POST請求並返回響應
func SendPostRequest(url string, reqBody string, handler ResponseHandler) error {
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	req.Header.SetContentType("application/x-www-form-urlencoded")

	// 設置請求正文
	req.SetBodyString(reqBody)

	// 創建一個新的響應對象
	resp := fasthttp.AcquireResponse()

	// 發送請求
	if err := fasthttp.Do(req, resp); err != nil {
		return err
	}

	// 調用回調函數處理響應
	return handler(resp)
}

func SendGetRequest(url string, handler ResponseHandler) error {
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("GET")
	req.SetRequestURI(url)
	req.Header.SetContentType("application/json;charset=utf-8")

	req.Body()
	// 創建一個新的響應對象
	resp := fasthttp.AcquireResponse()

	// 發送請求
	if err := fasthttp.Do(req, resp); err != nil {
		return err
	}

	// 調用回調函數處理響應
	return handler(resp)
}
