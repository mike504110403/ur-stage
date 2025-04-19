package wallet

import (
	"member_service/api/router/middleware/jwthandler"
	"member_service/internal/locals"
	"member_service/internal/proxy"
	"os"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/wallet")
	{
		g.Get("/detail", jwthandler.New(), detailHandler)                    // 錢包明細
		g.Get("/centerWallet", jwthandler.New(), centerWalletHandler)        // 中心錢包
		g.Get("/withdrawGate", jwthandler.New(), withdrawGateHandler)        // 提款限制
		g.Get("/withdrawGateFlow", jwthandler.New(), witdrawGateFlowHandler) // 提款流水限制
		g.Post("/initOrder", jwthandler.New(), initOrderHandler)             // 起單
		g.Post("/adminPayment", jwthandler.New(), adminPaymentHandler)       // 管理員支付
		g.Post("pay/htCallback", payHtCallbackHandler)                       // 回調
		g.Post("/payout/htCallback", payoutHtCallbackHandler)                // 回調
		g.Post("pay/asCallback", payAsCallbackHandler)                       // 回調
		g.Post("/payout/asCallback", payoutAsCallbackHandler)                // 回調
		g.Post("/withdraw", jwthandler.New(), withdrawHandler)               // 提現
	}
}

// 錢包明細
func detailHandler(c *fiber.Ctx) error {
	res := DetailRes{}
	localUser, err := locals.GetUserInfo(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}

	data, err := getWalletDetail(localUser.MemberId)
	if err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	res.Details = data
	return c.Status(fiber.StatusOK).JSON(res)
}

// 中心錢包
func centerWalletHandler(c *fiber.Ctx) error {
	localUser, err := locals.GetUserInfo(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}

	res, err := getCenterWallet(localUser.MemberId)
	if err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// 取得提款限制
func withdrawGateHandler(c *fiber.Ctx) error {
	uri := os.Getenv("WALLET_SERVER") + "/v1/wallet/withdrawGate"
	return proxy.PassingApi(c, uri)
}

// 提款流水限制
func witdrawGateFlowHandler(c *fiber.Ctx) error {
	uri := os.Getenv("WALLET_SERVER") + "/v1/wallet/withdrawGateFlow"
	return proxy.PassingApi(c, uri)
}

// 起單
func initOrderHandler(c *fiber.Ctx) error {
	uri := os.Getenv("WALLET_SERVER") + "/v1/order/init"
	return proxy.PassingApi(c, uri)
}

// 管理員後台上分
func adminPaymentHandler(c *fiber.Ctx) error {
	uri := os.Getenv("WALLET_SERVER") + "/v1/order/adminPayment"
	return proxy.PassingApi(c, uri)
}

// 支付回調
func payHtCallbackHandler(c *fiber.Ctx) error {
	uri := os.Getenv("WALLET_SERVER") + "/v1/order/htCallback"
	return proxy.PassingApiWithLocal(c, uri)
}

// 提款回調
func payoutHtCallbackHandler(c *fiber.Ctx) error {
	uri := os.Getenv("WALLET_SERVER") + "/v1/payment/htCallback"
	return proxy.PassingApiWithLocal(c, uri)
}

// 支付回調
func payAsCallbackHandler(c *fiber.Ctx) error {
	uri := os.Getenv("WALLET_SERVER") + "/v1/order/asCallback"
	return proxy.PassingApiWithLocal(c, uri)
}

// 提款回調
func payoutAsCallbackHandler(c *fiber.Ctx) error {
	uri := os.Getenv("WALLET_SERVER") + "/v1/payment/asCallback"
	return proxy.PassingApiWithLocal(c, uri)
}

// 提款
func withdrawHandler(c *fiber.Ctx) error {
	uri := os.Getenv("WALLET_SERVER") + "/v1/payment/withdraw"
	return proxy.PassingApi(c, uri)
}
