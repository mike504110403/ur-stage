package main

import (
	"bytes"
	"database/sql"
	"event_service/internal/cachedata"
	"event_service/internal/config"
	"event_service/internal/database"
	"event_service/internal/point"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"event_service/internal/cookies"

	routers "event_service/api/router"

	"github.com/gofiber/fiber/v2"
	mlog "github.com/mike504110403/goutils/log"
	"gitlab.com/gogogo2712128/common_moduals/apiprotocol"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

var (
	version     = "0.0.0"
	commitID    = "dev"
	environment = os.Getenv("environment")
	port        = os.Getenv("port")
)

// ServiceName : 服務名稱 | 不可修改，會影響資料一致性
const ServiceName = "event_service"

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
		RefreshDuration: time.Hour * 12,
		RetryDuration:   time.Second * 3,
	})
	// 設定參數初始化
	if err := typeparam.Init(typeparam.Config{
		FuncGetDB: func() (*sql.DB, error) {
			return database.SETTING.DB()
		},
	}); err != nil {
		mlog.Fatal(fmt.Sprintf("參數初始化錯誤: %S", err))
	}

	//初始化redis
	// redisConfig := redis.Config{
	// 	RedisServer:   env.SystemENV.RedisServer,
	// 	RedisPassword: env.SystemENV.RedisPassword,
	// 	RedisDB:       0,
	// }
	// redis.Init(redisConfig)

	// 初始化cookies
	cookies.Init(cookies.Config{
		MaxAge:   env.MaxAge,
		Secure:   false,
		SameSite: "Lax",
		HTTPOnly: false,
	})
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

	if err := routers.Set(app); err != nil {
		mlog.Fatal(fmt.Sprintf("routers設定失敗, err: %v", err))
	}

	// 啟動Server
	port := fmt.Sprintf(":%v", sysENV.Port)

	// 流水分發
	go point.CronGetBetFlow()

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
