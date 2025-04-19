package main

import (
	//"member_service/internal/database"

	"bytes"
	"database/sql"
	"fmt"
	"game_service/api/handler/private/v1/game"
	routers "game_service/api/router"
	"game_service/api/router/middleware/jwthandler"
	gws "game_service/internal/ws/gameclient"

	"game_service/internal/access"
	"game_service/internal/cachedata"
	"game_service/internal/config"
	"game_service/internal/database"
	"game_service/internal/locals"
	"game_service/internal/redis"

	"gitlab.com/gogogo2712128/common_moduals/apiprotocol"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"

	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
)

var (
	version     = "0.0.0"
	commitID    = "dev"
	environment = os.Getenv("environment")
	port        = os.Getenv("port")
)

// ServiceName : 服務名稱 | 不可修改，會影響資料一致性
const ServiceName = "gamr_service"

func Init() {
	if version == "0.0.0" {
		cmd := exec.Command("git", "rev-parse", "HEAD")
		out := bytes.Buffer{}
		cmd.Stdout = &out
		err := cmd.Run()
		if err == nil {
			version = strings.Replace(out.String(), "\n", "", -1)
		}
	} else {
		commitID = strings.Replace(commitID, "\"", "", -1)
		version = strings.Replace(version, "\"", "", -1)
	}
	// 環境變數初始化
	config.Init(version)
	config.EnvInit()
	env := config.GetENV()
	// 初始化log
	mlog.Init(mlog.Config{
		EnvMode: mlog.EnvMode(env.SystemENV.EnvMod),
		LogType: mlog.LogType(env.SystemENV.LogType),
	})

	// 初始化資料庫
	database.Init()
	// 初始化快取
	cachedata.Init(cachedata.Config{
		RefreshDuration: time.Minute * 3,
		RetryDuration:   time.Second * 3,
	})

	// 初始化參數
	if err := typeparam.Init(typeparam.Config{
		FuncGetDB: func() (*sql.DB, error) {
			return database.SETTING.DB()
		},
	}); err != nil {
		mlog.Fatal(fmt.Sprintf("參數初始化錯誤: %S", err))
	}
	// 初始化redis
	redisConfig := redis.Config{
		RedisServer:   env.SystemENV.RedisServer,
		RedisPassword: env.SystemENV.RedisPassword,
		RedisDB:       0,
	}
	redis.Init(redisConfig)

	// 初始化jwtHandler，設定金鑰以及其他相關設定
	if err := jwthandler.Init(jwthandler.Config{
		TokenLookupKey: fmt.Sprintf("cookie:%s", access.CookieKeyJWT),
		Secret:         env.JwtSecretkey,
		Expires:        env.JwtExpires,
		LocalsTokenKey: locals.KeyJWTToken,
		OnSuccess:      access.JWTOnSuccess,
		OnJWTError:     access.JWTOnError,
	}); err != nil {
		mlog.Fatal(fmt.Sprintf("[func->jwthandler.Init] %s", err))
	}
}

func main() {
	Init()
	sysENV := config.GetSystemENV()
	// 設定fiber的config
	app := fiber.New(fiber.Config{
		ReadTimeout:    sysENV.ReadTimeout,
		WriteTimeout:   sysENV.WriteTimeout,
		IdleTimeout:    sysENV.IdleTimeout,
		ReadBufferSize: int(sysENV.FiberHeaderSizeLimitMb) * 1024,
		BodyLimit:      int(sysENV.FiberBodyLimitMb) * 1024 * 1024,
		ErrorHandler:   apiprotocol.ErrorHandler(),
	})

	app.Get("/version", func(c *fiber.Ctx) error {
		statusData := fiber.Map{
			"environment":        environment,
			"Service":            ServiceName,
			"Version":            version,
			"commitID":           commitID,
			"OpenConnections":    c.App().Server().GetOpenConnectionsCount(),
			"CurrentConcurrency": c.App().Server().GetCurrentConcurrency(),
			"goVersion":          runtime.Version(),
			"config":             app.Config(),
		}
		return c.Status(fiber.StatusOK).JSON(statusData)
	})

	if err := routers.Set(app); err != nil {
		mlog.Fatal(fmt.Sprintf("routers設定失敗, err: %v", err))
	}

	go game.BetRecordCron()

	// 初始化WebSocketManager
	gws.NewWebSocketManager()
	// 啟動Server
	port := fmt.Sprintf(":%v", sysENV.Port)
	// 關閉服務
	go func() {
		if err := app.Listen(port); err != nil {
			mlog.Error(err.Error())
		}
	}()
	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	app.Shutdown()
	time.Sleep(5 * time.Second)
}
