package payment

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
	"gopkg.in/go-playground/validator.v9"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/payment")
	{
		g.Post("/withdraw", withdrawHandler)     // 提款
		g.Post("/htCallback", htCallbackHandler) // 提款完成 - 給第三方支付系統用
	}
}

// withdrawHandler 提款
func withdrawHandler(c *fiber.Ctx) error {
	mid := string(c.Request().Header.Peek("mid"))
	if mid == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}

	req := withdrawReq{}
	if err := c.BodyParser(&req); err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数解析失败"})
	}

	// req validate
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数验证失败"})
	}
	midInt, err := strconv.Atoi(mid)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	req.MemberId = midInt
	if err := withdraw(req); err != nil {
		if err.Error() == "未滿足提款條件" {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"message": "未满足提款条件或手续费不足"})
		}
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "成功"})
}

// htCallbackHandler 提款完成 - 給第三方支付系統用
func htCallbackHandler(c *fiber.Ctx) error {
	req := PayoutResult{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数解析失败"})
	}

	if err := doneWithdraw(req); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"message": "订单不存在"})
		}
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.SendString("success")
}
