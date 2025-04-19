package public

import (
	v1 "game_service/api/handler/public/v1"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/public")
	{
		v1.SetRouter(g)
	}
}
