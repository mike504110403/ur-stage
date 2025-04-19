package proxy

import (
	"fmt"
	gameoperatorhandler "member_service/api/router/middleware/gameOperatorHandler"
	"member_service/internal/locals"
	"strconv"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

// PassingApi : 透過fiber.Ctx將請求轉發到指定的uri
func PassingApi(c *fiber.Ctx, uri string) error {
	localUser, err := locals.GetUserInfo(c)
	if err != nil {
		mlog.Error(fmt.Sprintf("Error getting user info: %v", err))
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登入"})
	}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(uri)
	req.Header.SetMethod(string(c.Request().Header.Method()))
	mid := strconv.Itoa(localUser.MemberId)
	req.Header.Set("mid", mid)
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(c.Request().Body())

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	if err != nil {
		mlog.Error(fmt.Sprintf("Error forwarding request: %v", err))
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	c.Response().Header.SetStatusCode(resp.StatusCode())
	resp.Header.VisitAll(func(key, value []byte) {
		c.Response().Header.SetBytesKV(key, value) // 設置下游服務的所有標頭
	})
	c.Response().SetBody(resp.Body())
	// 釋放redis鎖
	gameoperatorhandler.ReleaseRedisLock(localUser.MemberId)

	return nil
}

func PassingApiWithLocal(c *fiber.Ctx, uri string) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(uri)
	req.Header.SetMethod(string(c.Request().Header.Method()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBody(c.Request().Body())

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		mlog.Error(fmt.Sprintf("Error forwarding request: %v", err))
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	c.Response().Header.SetStatusCode(resp.StatusCode())
	c.Response().Header.Set("Content-Type", string(resp.Header.Peek("Content-Type")))
	c.Response().SetBody(resp.Body())

	return nil
}
