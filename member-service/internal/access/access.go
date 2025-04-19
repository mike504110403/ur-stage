package access

import (
	"fmt"
	"member_service/internal/locals"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var cfg Config

func Init(initConfig Config) {
	cfg = initConfig
}

func JWTOnSuccess(c *fiber.Ctx) error {
	jwtData := getJwtData(c)
	if jwtData == nil {
		CookieKeyJWT.Clear(c)
		return c.SendStatus(http.StatusUnauthorized)
	}

	if !jwtData.Valid {
		CookieKeyJWT.Clear(c)
		return c.SendStatus(http.StatusUnauthorized)
	}

	claims, err := parseTokenToClaims(jwtData.Raw)
	if err != nil {
		CookieKeyJWT.Clear(c)
		return c.SendStatus(http.StatusUnauthorized)
	}

	// 沒有帳號的話直接回傳失敗
	if claims.UserName == "" || claims.MemberId == 0 {
		CookieKeyJWT.Clear(c)
		return c.SendStatus(http.StatusUnauthorized)
	} else {
		locals.SetUserInfo(c, locals.UserInfo{
			MemberId: claims.MemberId,
			UserName: claims.UserName,
		})
	}

	// 如果剩下不到一半的存活時間就給新的token
	if claims.ExpiresAt.Before(time.Now().Add(-1 * time.Second * time.Duration(cfg.JWTExpires/2))) {
		// 更新token exp
		jwtToken, err := generateToken(claims.MemberId, claims.UserName, time.Now().Add(time.Second*time.Duration(cfg.JWTExpires)))
		if err != nil {
			CookieKeyJWT.Clear(c)
			return c.SendStatus(http.StatusUnauthorized)
		}
		CookieKeyJWT.Set(c, jwtToken)
	}

	return c.Next()
}

func JWTOnError(c *fiber.Ctx, err error) error {
	CookieKeyJWT.Clear(c)
	return c.SendStatus(http.StatusUnauthorized)
}

// LoginSuccess : 登入成功後的行為，會以uid產生token後存入cookies
func LoginSuccess(c *fiber.Ctx, id int, username string) (string, error) {
	// 產生token
	jwtToken, err := generateToken(id, username, time.Now().Add(time.Second*time.Duration(cfg.JWTExpires)))
	fmt.Printf("jwtToken: %s\n", jwtToken)
	if err != nil {
		return "", err
	}

	// 設定token至cookies
	CookieKeyJWT.Set(c, jwtToken)

	return jwtToken, nil
}

// Logout : 登出
func Logout(c *fiber.Ctx) {
	// 清除token的cookies
	CookieKeyJWT.Clear(c)
	locals.ClearLocal(c)
}

// generateToken : 產生金鑰
func generateToken(id int, username string, exp time.Time) (jwtToken string, err error) {
	claims := Claims{
		MemberId: id,
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtToken, err = token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// parseTokenToClaims : 從JWT中取得使用者資訊
func parseTokenToClaims(token string) (*Claims, error) {
	claims := &Claims{}

	if _, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return []byte(cfg.JWTSecret), nil
	}); err != nil {
		return nil, err
	}

	return claims, nil
}

// getJwtData : 透過Locals取得jwt
func getJwtData(ctx *fiber.Ctx) *jwt.Token {
	if ctx.Locals(locals.KeyJWTToken) != nil {
		switch ctx.Locals(locals.KeyJWTToken).(type) {
		case *jwt.Token:
			return ctx.Locals(locals.KeyJWTToken).(*jwt.Token)
		}
	}
	return nil
}
