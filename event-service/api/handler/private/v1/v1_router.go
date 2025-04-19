package v1

import (
	bonus "event_service/api/handler/private/v1/bonus"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/v1")
	{
		bonus.SetRouter(g)
	}
}
