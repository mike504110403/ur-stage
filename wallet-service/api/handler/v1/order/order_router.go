package order

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
	"gopkg.in/go-playground/validator.v9"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/order")
	{
		g.Post("/init", initOrderHandler)            // 訂單起單
		g.Post("/payment", paymentHandler)           // 訂單收款 - 給第三方支付系統用
		g.Post("/htCallback", htCallbackHandler)     // 訂單收款 - 給 ht_pay 用
		g.Post("/adminPayment", adminPaymentHandler) // 訂單收款 - 給管理員用
	}
}

// initOrderHandler 訂單起單
func initOrderHandler(c *fiber.Ctx) error {
	req, res, url := OrderInitReq{}, OrderInitRes{}, ""
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数解析失败"})
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数验证失败"})
	}
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
	cip := c.IP()
	if data, err := orderInit(req, mid, cip); err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	} else {
		url = data
	}
	res.Url = url
	return c.JSON(res)
}

// adminPaymentHandler 管理員後台上分
func adminPaymentHandler(c *fiber.Ctx) error {
	req := AdminPaymentReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数解析失败"})
	}
	// 檢查請求中是否有 mid
	midStr := string(c.Request().Header.Peek("mid"))
	mid, err := strconv.Atoi(midStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	// 檢查權限
	if auth := checkAuth(mid); !auth {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	} else {
		if err := adminPayment(req, mid); err != nil {
			if err == errors.New("無此會員") {
				return c.Status(fiber.StatusBadRequest).
					JSON(fiber.Map{"message": "未登录"})
			}
			mlog.Error(err.Error())
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"message": "服务器错误"})
		}
	}
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "转点设置成功"})
}

// paymentOrderHandler 訂單收款 - 三方收款用
func paymentHandler(c *fiber.Ctx) error {
	req := PayCallbackReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数解析失败"})
	}

	if err := paymentOrder(req); err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "订单收款成功"})
}

// htCallbackHandler 訂單收款 - ht_pay 用
func htCallbackHandler(c *fiber.Ctx) error {
	mlog.Info("請求收款")
	req := PayCallbackReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数解析失败"})
	}

	if err := paymentOrder(req); err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.SendString("success")
}
