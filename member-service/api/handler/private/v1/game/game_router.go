package game

import (
	gameoperatorhandler "member_service/api/router/middleware/gameOperatorHandler"
	"member_service/api/router/middleware/jwthandler"
	"member_service/internal/proxy"
	"os"

	mlog "github.com/mike504110403/goutils/log"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/go-playground/validator.v9"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/game")
	{
		g.Post("/joinGame", jwthandler.New(), gameoperatorhandler.CheckGameOperator(), joinGameHandler)    // 加入遊戲
		g.Post("/gamelist", gameListHandler)                                                               // 遊戲列表
		g.Get("/recycle", jwthandler.New(), gameoperatorhandler.CheckGameOperator(), recycleHandler)       // 回收點數
		g.Post("/bringin", jwthandler.New(), gameoperatorhandler.CheckGameOperator(), bringinHandler)      // 轉入點數s
		g.Get("/kickoutall", jwthandler.New(), gameoperatorhandler.CheckGameOperator(), kickoutAllHandler) // 踢出遊戲
	}
}

// 加入遊戲
func joinGameHandler(c *fiber.Ctx) error {
	uri := os.Getenv("GAME_SERVER") + "/private/v1/game/joingame"
	return proxy.PassingApi(c, uri)
}

// 遊戲列表
func gameListHandler(c *fiber.Ctx) error {
	req := GameListReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求格式错误"})
	}

	// req validate
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数验证失败"})
	}

	res, err := getGameList(req)
	if err != nil {
		mlog.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(res)
}

// 回收點數
func recycleHandler(c *fiber.Ctx) error {
	uri := os.Getenv("GAME_SERVER") + "/private/v1/game/recycle"
	return proxy.PassingApi(c, uri)
}

// 轉入點數
func bringinHandler(c *fiber.Ctx) error {
	uri := os.Getenv("GAME_SERVER") + "/private/v1/game/bringin"
	return proxy.PassingApi(c, uri)
}

// 踢出遊戲
func kickoutAllHandler(c *fiber.Ctx) error {
	uri := os.Getenv("GAME_SERVER") + "/private/v1/game/kickoutall"
	return proxy.PassingApi(c, uri)
}
