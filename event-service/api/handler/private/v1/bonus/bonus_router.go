package bonus

import "github.com/gofiber/fiber/v2"

func SetRouter(router fiber.Router) {
	g := router.Group("/bonus")
	{
		g.Post("/test", testHandler)
	}
}

func testHandler(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
