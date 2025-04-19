package apicaller

import (
	"encoding/json"
	"fmt"

	mlog "github.com/mike504110403/goutils/log"

	"testing"

	"github.com/valyala/fasthttp"
)

func TestSendPostRequest(t *testing.T) {
	// 定義API路由
	url := "https://api.example.com/data"

	// 定義頭部
	headers := map[string]string{
		"Content-Type": "application/json",
		// 這邊應該是要放遊戲商的key
		"Authorization": "Bearer your_token_here",
	}

	// 定義正文
	bodyData := map[string]interface{}{
		"key": "value",
	}
	// 也可以使用結構體
	// ex. type BodyData struct {
	// 	Test string `json:"test"`
	// }
	// bodyData := BodyData{
	// 	Key: "value",
	// }

	body, err := json.Marshal(bodyData)
	if err != nil {
		mlog.Fatal(err.Error())
	}

	// 定義回調函數
	handler := func(resp *fasthttp.Response) error {
		// 取得響應狀態碼
		statusCode := resp.StatusCode()
		mlog.Info(fmt.Sprintf("Status Code: %d\n", statusCode))

		// 取得響應正文
		responseBody := resp.Body()
		mlog.Info(fmt.Sprintf("Response Body: %s\n", responseBody))
		return nil
	}

	// 發送POST請求並處理響應
	if err := SendPostRequest(url, headers, body, handler); err != nil {
		mlog.Fatal(err.Error())
	}
}
