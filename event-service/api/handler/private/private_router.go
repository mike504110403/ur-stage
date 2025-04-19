package private

import (
	v1 "event_service/api/handler/private/v1"

	"github.com/gofiber/fiber/v2"
)

func SetRouter(router fiber.Router) {
	g := router.Group("private")
	{
		v1.SetRouter(g)
	}
}
