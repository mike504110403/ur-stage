package locals

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// 統一處理API中使用的locals

// GetUserInfo : 透過Locals取得使用者資訊
func GetUserInfo(ctx *fiber.Ctx) (UserInfo, error) {
	if ctx.Locals(KeyUserInfo) != nil {
		switch ctx.Locals(KeyUserInfo).(type) {
		case UserInfo:
			userInfo := ctx.Locals(KeyUserInfo).(UserInfo)
			return userInfo, nil
		}
	}
	return UserInfo{}, errors.New("userInfo locals is not exist")
}

// SetUserInfo : 透過Locals設定使用者資訊
func SetUserInfo(ctx *fiber.Ctx, userInfo UserInfo) {
	ctx.Locals(KeyUserInfo, userInfo)
}

// GetJwt : 透過Locals取得jwt
func GetJwt(ctx *fiber.Ctx) *jwt.Token {
	if ctx.Locals(KeyJWTToken) != nil {
		switch ctx.Locals(KeyJWTToken).(type) {
		case *jwt.Token:
			return ctx.Locals(KeyJWTToken).(*jwt.Token)
		}
	}
	return nil
}

// Logout : 處理用戶登出並清除上下文中的資料
func ClearLocal(ctx *fiber.Ctx) {
	// 清除用戶信息
	ctx.Locals(KeyUserInfo, nil)

	// 清除 JWT token
	ctx.Locals(KeyJWTToken, nil)
}
