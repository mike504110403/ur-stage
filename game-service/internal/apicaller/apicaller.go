package apicaller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"

	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

type ResponseHandler func(*fasthttp.Response) error

func GBGameSendPostRequest(url string, headers *map[string]string, body string, handler ResponseHandler) error {
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	// 設置請求頭部
	if headers != nil {
		for key, value := range *headers {
			req.Header.Set(key, value)
		}
	}

	// 將 hashKey 和 body 組合成一個 JSON 對象
	requestBody := map[string]interface{}{
		"x": body,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// 設置請求正文
	req.Header.SetContentType("application/json")
	req.SetBody(jsonData)

	apiclient := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		ReadTimeout:              30 * time.Second,
		WriteTimeout:             30 * time.Second,
	}
	// 創建一個新的響應對象
	resp := fasthttp.AcquireResponse()

	// 發送請求
	if err := apiclient.Do(req, resp); err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return err
	}

	// 調用回調函數處理響應
	return handler(resp)
}

func RSGSendPostRequest(url string, headers map[string]string, body []byte, handler ResponseHandler) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return err
	}

	return handler(resp)
}

// SendPostRequest 發送POST請求並返回響應
func SendPostRequest(url string, headers map[string]string, body []byte, handler ResponseHandler) error {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	if err := writer.WriteField("msg", string(body)); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	// 設置請求頭部
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 設置請求正文
	req.Header.SetContentType(writer.FormDataContentType())
	req.SetBody(buffer.Bytes())

	apiclient := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		ReadTimeout:              30 * time.Second,
		WriteTimeout:             30 * time.Second,
	}
	// 創建一個新的響應對象
	resp := fasthttp.AcquireResponse()

	// 發送請求
	if err := apiclient.Do(req, resp); err != nil {
		return err
	}

	// 調用回調函數處理響應
	return handler(resp)
}

// SendPostRequest 發送POST請求並返回響應
func MTSendPostRequest(url string, headers *map[string]string, body string, handler ResponseHandler, hashKey string) error {
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	// 設置請求頭部
	if headers != nil {
		for key, value := range *headers {
			req.Header.Set(key, value)
		}
	}

	// 將 hashKey 和 body 組合成一個 JSON 對象
	requestBody := map[string]interface{}{
		"HashKey": hashKey,
		"Data":    body,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// 設置請求正文
	req.Header.SetContentType("application/json")
	req.SetBody(jsonData)

	apiclient := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		ReadTimeout:              30 * time.Second,
		WriteTimeout:             30 * time.Second,
	}
	// 創建一個新的響應對象
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 釋放資源

	// 發送請求
	if err := apiclient.Do(req, resp); err != nil {
		return err
	}

	// 調用回調函數處理響應
	return handler(resp)
}

func WgGameSendPostRequest(url string, headers map[string]string, body []byte, handler ResponseHandler) error {
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	// 設置請求頭部
	for key, value := range headers {
		req.Header.Set(key, value)
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

	// 發送請求
	if err := apiclient.Do(req, resp); err != nil {
		return err
	}

	// 調用回調函數處理響應
	return handler(resp)
}

func WgGameSendGetRequest(url string, headers map[string]string, handler ResponseHandler) error {
	// 創建一個新的請求對象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("Get")
	req.SetRequestURI(url)
	// 設置請求頭部
	for key, value := range headers {
		req.Header.Set(key, value)
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

	// 發送請求
	if err := apiclient.Do(req, resp); err != nil {
		return err
	}

	// 調用回調函數處理響應
	return handler(resp)
}
