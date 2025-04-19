package bonus

import (
	"member_service/api/router/middleware/jwthandler"
	"member_service/internal/proxy"
	"os"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/bonus")
	{
		g.Get("/firstDeposit", jwthandler.New(), firstDepositHandler) // 首存優惠
	}
}

// 首存優惠
func firstDepositHandler(c *fiber.Ctx) error {
	uri := os.Getenv("WALLET_SERVER") + "/v1/bonus/firstDeposit"
	return proxy.PassingApi(c, uri)
}
