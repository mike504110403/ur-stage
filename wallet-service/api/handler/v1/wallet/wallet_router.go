package wallet

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
	"gopkg.in/go-playground/validator.v9"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/wallet")
	{
		g.Post("/transfer", transferHandler)                // 錢包轉點
		g.Get("/withdrawGate", withdrawGateHandler)         // 提款限制
		g.Get("/withdrawGateFlow", withdrawGateFlowHandler) // 提款限制流水
	}
}

// 提款限制
func withdrawGateHandler(c *fiber.Ctx) error {
	// 檢查請求中是否有 mid
	midStr := string(c.Request().Header.Peek("mid"))
	if midStr == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	mid, err := strconv.Atoi(midStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	res, err := getWithdrawGate(mid)
	if err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	return c.JSON(res)
}

// 提款流水限制
func withdrawGateFlowHandler(c *fiber.Ctx) error {
	// 檢查請求中是否有 mid
	midStr := string(c.Request().Header.Peek("mid"))
	if midStr == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	mid, err := strconv.Atoi(midStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	data, err := getWithdrawGateFlow(mid)
	if err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	res := WithdrawGateFlowRes{
		WithdrawEvents: data,
	}

	return c.JSON(res)
}

// 錢包轉點
func transferHandler(c *fiber.Ctx) error {
	req := TransferReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数解析失败"})
	}
	// req validate
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数验证失败"})
	}
	// sql 轉帳紀錄和
	if err := transferAndRecord(req); err != nil {
		if err.Error() == "餘額不足" {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"message": "余额不足"})
		}
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "转点成功"})
}
