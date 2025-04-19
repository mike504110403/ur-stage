package game

import (
	"github.com/gofiber/fiber/v2"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/game")
	{
		g.Get("/transfer", transferHandler) // 錢包轉點
	}
}

func transferHandler(c *fiber.Ctx) error {
	return nil
}
