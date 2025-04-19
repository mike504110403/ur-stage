package v1

import (
	bonus "wallet_service/api/handler/v1/bonus"
	order "wallet_service/api/handler/v1/order"
	payment "wallet_service/api/handler/v1/payment"
	wallet "wallet_service/api/handler/v1/wallet"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/v1")
	{
		wallet.SetRouter(g)
		order.SetRouter(g)
		payment.SetRouter(g)
		bonus.SetRouter(g)
	}
}
