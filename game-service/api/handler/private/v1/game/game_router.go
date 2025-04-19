package game

import (
	"fmt"
	"game_service/internal/ws/gameclient"
	"strconv"

	"gopkg.in/go-playground/validator.v9"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/game")
	{
		g.Post("/joingame", joinGameHandler)    // 加入遊戲
		g.Get("/recycle", recycleHandler)       // 回收
		g.Get("/kickoutall", kickOutAllHandler) // 踢出所有玩家
		//g.Post("/bringin", bringInHandler)   // 點數吸入
	}
}

// 進入遊戲
func joinGameHandler(c *fiber.Ctx) error {
	// 解析請求中的 JSON 資料
	req := JoinGameReq{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数解析失败"})
	}

	// req validate
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数验证失败"})
	}

	// 檢查請求中是否有 mid
	mid := string(c.Request().Header.Peek("mid"))
	if mid == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	midInt, err := strconv.Atoi(mid)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	url, err := JoinGame(midInt, req.AgentId)
	if err != nil {
		mlog.Error(fmt.Sprintf("join game failed: %v", err))
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	// 假設加入成功，回應前端
	return c.Status(fiber.StatusOK).JSON(url)
}

// 回收遊戲內所有點數
func recycleHandler(c *fiber.Ctx) error {
	// 檢查請求中是否有 mid
	mid := string(c.Request().Header.Peek("mid"))
	if mid == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	midInt, err := strconv.Atoi(mid)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	if err := Recycle(midInt); err != nil {
		if err.Error() == "请先离开游戏" {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"message": "请先离开游戏"})
		}
		mlog.Info(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "点数回收成功"})
}

// 踢出所有玩家
func kickOutAllHandler(c *fiber.Ctx) error {
	// 檢查請求中是否有 mid
	mid := string(c.Request().Header.Peek("mid"))
	if mid == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	midInt, err := strconv.Atoi(mid)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}

	if err := KickOutAll(midInt); err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	} else {
		gameclient.GameWebSocket.RemoveClient(midInt)
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "游戏登出成功"})
}

// func getDonateReportHandler(c *fiber.Ctx) error {
// 	// 解析請求中的 JSON 資料
// 	req, res := DonateReq{}, DonateRes{}
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).
// 			JSON(fiber.Map{"message": "请求参数解析失败"})
// 	}

// 	// req validate
// 	validate := validator.New()
// 	if err := validate.Struct(req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).
// 			JSON(fiber.Map{"message": "请求参数验证失败"})
// 	}

// 	status, err := mt_live.GetBetRecord(time.Duration(30 * time.Second))
// 	if err != nil {
// 		mlog.Error(err.Error())
// 		return c.Status(fiber.StatusInternalServerError).
// 			JSON(fiber.Map{"message": "服务器错误"})
// 	} else {
// 		res.Status = status
// 	}

// 	// 假設加入成功，回應前端
// 	return c.Status(fiber.StatusOK).JSON(res)
// }

// func getLiveRecordHandler(c *fiber.Ctx) error {
// 	// 解析請求中的 JSON 資料
// 	req, res := LiveRecordReq{}, LiveRecordRes{}
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).
// 			JSON(fiber.Map{"message": "请求参数解析失败"})
// 	}

// 	// req validate
// 	validate := validator.New()
// 	if err := validate.Struct(req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).
// 			JSON(fiber.Map{"message": "请求参数验证失败"})
// 	}

// 	status, err := mt_live.GetBetRecord(time.Duration(30 * time.Second))
// 	if err != nil {
// 		mlog.Error(err.Error())
// 		return c.Status(fiber.StatusInternalServerError).
// 			JSON(fiber.Map{"message": "服务器错误"})
// 	} else {
// 		res.Status = status
// 	}

// 	// 假設加入成功，回應前端
// 	return c.Status(fiber.StatusOK).JSON(res)
// }

// func getLotteryRecorddHandler(c *fiber.Ctx) error {
// 	// 解析請求中的 JSON 資料
// 	req, res := LiveRecordReq{}, LiveRecordRes{}
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).
// 			JSON(fiber.Map{"message": "请求参数解析失败"})
// 	}

// 	// req validate
// 	validate := validator.New()
// 	if err := validate.Struct(req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).
// 			JSON(fiber.Map{"message": "请求参数验证失败"})
// 	}

// 	status, err := mt_lottery.GetBetOrder(time.Duration(30 * time.Second))
// 	if err != nil {
// 		mlog.Error(err.Error())
// 		return c.Status(fiber.StatusInternalServerError).
// 			JSON(fiber.Map{"message": "服务器错误"})
// 	} else {
// 		res.Status = status
// 	}

// 	// 假設加入成功，回應前端
// 	return c.Status(fiber.StatusOK).JSON(res)
// }
