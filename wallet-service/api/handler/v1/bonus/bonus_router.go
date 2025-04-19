package bonus

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
)

func SetRouter(router fiber.Router) {
	g := router.Group("/bonus")
	{
		g.Get("/firstDeposit", firstDepositHandler) // 首存優惠
	}
}

func firstDepositHandler(c *fiber.Ctx) error {
	// 檢查請求中是否有 mid
	midStr := string(c.Request().Header.Peek("mid"))
	if midStr == "" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	mid, err := strconv.Atoi(midStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "未登录"})
	}
	if err := firstDepositapply(mid); err != nil {
		if err == sql.ErrNoRows || err.Error() == "首存優惠資格不符" {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"message": "首存优惠资格不符"})
		}
		mlog.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"message": "服务器错误"})
	}
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"message": "首存优惠申请成功"})
}
