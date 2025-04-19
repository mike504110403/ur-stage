package v1

import (
	auth "member_service/api/handler/private/v1/auth"
	bonus "member_service/api/handler/private/v1/bonus"
	game "member_service/api/handler/private/v1/game"
	member "member_service/api/handler/private/v1/member"
	wallet "member_service/api/handler/private/v1/wallet"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/v1")
	{
		auth.SetRouter(g)
		game.SetRouter(g)
		member.SetRouter(g)
		wallet.SetRouter(g)
		bonus.SetRouter(g)
	}
}
