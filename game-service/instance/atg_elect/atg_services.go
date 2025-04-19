package atg_elect

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

func getToken() (string, error) {
	// 創建響應對象和請求對象
	res := GetTokenRes{}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 確保釋放資源

	// 設置請求方法和 URL
	req.Header.SetMethod("GET")
	req.SetRequestURI("https://api.godeebxp.com/" + string(TokenUri))

	// 添加必要的 Header
	req.Header.Set("X-Operator", "Ur_USDT_beta")
	req.Header.Set("X-Key", "ae6bc6d729894f088615aa1e772cdef5")

	// 配置 HTTP 客戶端
	apiclient := fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		ReadTimeout:              30 * time.Second,
		WriteTimeout:             30 * time.Second,
	}

	// 創建一個響應對象
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 確保釋放資源

	// 發送請求並檢查錯誤
	if err := apiclient.Do(req, resp); err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	// 檢查 HTTP 狀態碼
	if resp.StatusCode() != fasthttp.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	// 解析 JSON 響應
	if err := json.Unmarshal(resp.Body(), &res); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return res.Data.Token, nil
}

func apiCallerPost(uri string, reqBody *[]byte, handler func(r *fasthttp.Response) error) error {
	// 創建一個新的請求對象
	token, err := getToken()
	if err != nil {
		return err
	}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI("https://api.godeebxp.com/" + uri)

	req.Header.Set("X-Token", token)

	// 設置請求正文
	req.Header.SetContentType("application/json")
	if reqBody != nil {
		req.SetBody(*reqBody)
	}

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
	return handler(resp)
}

func apiCallerGet(uri string, urlParamsStr string, handler func(r *fasthttp.Response) error) error {
	// 創建一個新的請求對象
	token, err := getToken()
	if err != nil {
		return err
	}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("GET")
	req.SetRequestURI("https://api.godeebxp.com/" + uri + "?" + urlParamsStr)

	req.Header.Set("X-Token", token)

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
	return handler(resp)
}

func StructToQueryString(data interface{}) string {
	v := reflect.ValueOf(data)
	t := v.Type()
	var params []string

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		value := fmt.Sprintf("%v", v.Field(i).Interface())
		if tag != "" {
			params = append(params, fmt.Sprintf("%s=%s", tag, url.QueryEscape(value)))
		}
	}

	return strings.Join(params, "&")
}

func apiCallerDelete(uri string, handler func(r *fasthttp.Response) error) error {
	// 創建一個新的請求對象
	token, err := getToken()
	if err != nil {
		return err
	}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("DELETE")
	req.SetRequestURI("https://api.godeebxp.com/" + uri)

	req.Header.Set("X-Token", token)

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
	return handler(resp)
}

func apiCallerFormDataPost(uri string, formData map[string]string, handler func(r *fasthttp.Response) error) error {

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req) // 釋放資源

	// 設置請求方法和URL
	req.Header.SetMethod("POST")
	req.SetRequestURI("https://api.godeebxp.com/" + uri)

	// 設置 Header（標頭）
	req.Header.SetContentType("application/x-www-form-urlencoded") // 或 multipart/form-data

	// 構建 form-data
	form := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(form) // 釋放資源
	for key, value := range formData {
		form.Set(key, value)
	}
	req.SetBody(form.QueryString()) // 將表單資料設為請求的 Body

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
	return handler(resp)
}

// 將結構體轉換為 map[string]string 的函數
func structToMap(req interface{}) (map[string]string, error) {
	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	formData := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i).Interface()
		if value != nil {
			formData[field.Tag.Get("json")] = fmt.Sprintf("%v", value)
		}
	}
	return formData, nil
}
