package reshandler

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ResponseMiddleware(c *fiber.Ctx) error {
	// 執行下一個中間件或處理器
	if err := c.Next(); err != nil {
		return err
	}

	// 如果需要跳過格式化，直接返回
	if skipFormat, ok := c.Locals("skip_format").(bool); ok && skipFormat {
		return nil
	}

	res := Res{}
	statusCode := c.Response().StatusCode() // 取得狀態碼
	body := c.Response().Body()             // 取得原始回應內容

	// 檢查回應是否為有效 JSON 並解析
	var parsedBody map[string]interface{}
	isJSON := json.Valid(body)
	if isJSON {
		_ = json.Unmarshal(body, &parsedBody)
	}

	// 錯誤處理 (狀態碼非 200 的情況)
	switch statusCode {
	case fiber.StatusUnauthorized:
		res.Message = "未登录"
	case fiber.StatusInternalServerError:
		res.Message = "服务器错误"
		statusCode = fiber.StatusBadRequest
	case fiber.StatusOK:
		// 嘗試從 Locals 中提取 message，否則使用預設的 "success"
		if message, ok := c.Locals("message").(string); ok {
			res.Message = message
		} else {
			res.Message = "success"
		}

		// 如果 body 是有效 JSON，則直接使用 body 作為 Data，否則轉換為文本
		if isJSON {
			if message, ok := parsedBody["message"].(string); ok {
				res.Message = message
			} else {
				res.Data = json.RawMessage(body)
			}
		} else {
			res.Data = json.RawMessage(fmt.Sprintf(`"%s"`, body))
		}
	default:
		// 嘗試提取 message
		if isJSON {
			if message, ok := parsedBody["message"].(string); ok {
				res.Message = message
			}
		} else {
			// 非 JSON 的錯誤回應，直接使用原始內容作為 message
			if message, ok := c.Locals("message").(string); ok {
				res.Message = message
			} else {
				res.Message = string(body)
			}
		}
		res.Data = nil // 錯誤時不需要返回 Data
		statusCode = fiber.StatusBadRequest
	}
	if statusCode == fiber.StatusOK {
		res.Code = 0
	} else {
		res.Code = 1
	}

	return c.Status(statusCode).JSON(res)
}
