package member

import (
	"member_service/api/router/middleware/jwthandler"
	"member_service/internal/locals"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
	"gopkg.in/go-playground/validator.v9"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/member")
	{
		g.Get("/status", jwthandler.New(), statusHandler)             // 取得會員狀態
		g.Get("/info", jwthandler.New(), infoHandler)                 // 取得會員基本資料
		g.Post("/info", jwthandler.New(), updateInfoHandler)          // 更新會員基本資料
		g.Get("/transRecords", jwthandler.New(), transRecordsHandler) // 取得會員交易紀錄
		g.Get("/betRecords", jwthandler.New(), betRecordsHandler)     // 取得會員投注紀錄
	}
}

// 取得會員基本狀態
func statusHandler(c *fiber.Ctx) error {
	localUser, err := locals.GetUserInfo(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}

	res, err := GetMemberStatus(localUser.MemberId)
	if err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.JSON(res)
}

// 取得會員基本資料
func infoHandler(c *fiber.Ctx) error {
	localUser, err := locals.GetUserInfo(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}

	res, err := GetMemberInfo(localUser.MemberId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

// 更新會員基本資料
func updateInfoHandler(c *fiber.Ctx) error {
	req := BasicInfo{}
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

	localUser, err := locals.GetUserInfo(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}

	err = EditMemberInfo(localUser.MemberId, req)
	if err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "更新成功"})
}

// 取得會員交易紀錄
func transRecordsHandler(c *fiber.Ctx) error {
	// 從查詢參數中提取資料
	startDate := c.Query("start_date") // 必填字段
	endDate := c.Query("end_date")     // 必填字段
	typeParam := c.QueryInt("type")    // 可選字段，默認值設為 0

	// 初始化請求結構體
	req := TransRecordReq{
		StartDate: startDate, // 從查詢參數中提取的 start_date
		EndDate:   endDate,   // 從查詢參數中提取的 end_date
	}
	// 如果 type 不為 0，設置 type 指標
	if typeParam != 0 {
		req.Type = &typeParam
	}
	// req validate
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数验证失败"})
	}

	localUser, err := locals.GetUserInfo(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}

	res, err := GetTransRecords(localUser.MemberId, req)
	if err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

// 取得會員投注紀錄
func betRecordsHandler(c *fiber.Ctx) error {
	gameType := c.QueryInt("game_type", 0) // GameType 是指標類型，可以為空
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	// 構建請求結構體
	req := BetRecordReq{
		StartDate: startDate, // 從查詢參數獲取的開始日期
		EndDate:   endDate,   // 從查詢參數獲取的結束日期
	}

	// 如果 gameType 不為空，則設置指標
	if gameType != 0 {
		req.GameType = &gameType
	}

	// req validate
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "请求参数验证失败"})
	}

	localUser, err := locals.GetUserInfo(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}

	res, err := GetBetRecords(localUser.MemberId, req)
	if err != nil {
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}
