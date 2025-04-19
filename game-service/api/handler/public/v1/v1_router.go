package v1

import (
	game "game_service/api/handler/public/v1/game"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/v1")
	{
		game.SetRouter(g)
	}
}
