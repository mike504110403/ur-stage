package jwthandler

import (
	"errors"
	"fmt"
	"game_service/internal/access"
	"game_service/internal/locals"

	"github.com/golang-jwt/jwt/v4"
	mlog "github.com/mike504110403/goutils/log"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// Init : 初始化
func Init(initConfig Config) error {
	cfg = initConfig
	return checkConfig()
}

// checkConfig : 檢查所需的參數是否已設定
func checkConfig() error {
	if len(cfg.Secret) < 5 {
		return errors.New("jwt密鑰長度過短")
	}
	if cfg.Expires < 1 {
		mlog.Info(fmt.Sprintf("jwt expires: %d", cfg.Expires))
		return errors.New("maxAgeSecond過短")
	}
	if len(cfg.LocalsTokenKey) == 0 {
		return errors.New("未設定LocalsTokenKey")
	}

	return nil
}

// New : 回傳jwtHandler
func New() fiber.Handler {
	if err := checkConfig(); err != nil {
		mlog.Fatal(err.Error())
	}
	return jwtware.New(jwtware.Config{
		SigningKey:  []byte(cfg.Secret),
		TokenLookup: cfg.TokenLookupKey,
		ContextKey:  cfg.LocalsTokenKey,
		SuccessHandler: func(c *fiber.Ctx) error {
			return cfg.OnSuccess(c)
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			mlog.Error(fmt.Sprintf("jwt error: %v", err))
			return cfg.OnJWTError(c, err)
		},
	})
}

func WSNew() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := &access.Claims{}
		tokenStr := c.Query("token")
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return []byte(cfg.Secret), nil
		})
		if err != nil {
			mlog.Error(fmt.Sprintf("jwt error: %v", err))
			return err
		}

		// 驗證token是否有效
		if !token.Valid {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// 解析token
		if claims, ok := token.Claims.(*access.Claims); ok && token.Valid {
			locals.SetUserInfo(c, locals.UserInfo{
				MemberId: claims.MemberId,
				UserName: claims.UserName,
			})
		}
		mlog.Info(fmt.Sprintf("claims: %v", claims))

		return c.Next()
	}
}
