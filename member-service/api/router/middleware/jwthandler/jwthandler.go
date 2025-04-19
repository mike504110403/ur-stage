package jwthandler

import (
	"errors"
	"fmt"

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
