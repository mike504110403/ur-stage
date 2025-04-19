package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
	routers "wallet_service/api/router"
	htpay "wallet_service/instance/ht_pay"
	"wallet_service/internal/cachedata"
	"wallet_service/internal/config"
	"wallet_service/internal/database"
	"wallet_service/internal/redis"
	"wallet_service/internal/transection"

	// "wallet_service/internal/redis"
	//"wallet_service/internal/transfer"
	"gitlab.com/gogogo2712128/common_moduals/apiprotocol"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"

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
const ServiceName = "wallet_service"

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
	// 設定參數初始化
	if err := typeparam.Init(typeparam.Config{
		FuncGetDB: func() (*sql.DB, error) {
			return database.SETTING.DB()
		},
	}); err != nil {
		mlog.Fatal(fmt.Sprintf("參數初始化錯誤: %S", err))
	}

	//初始化redis
	redisConfig := redis.Config{
		RedisServer:   env.SystemENV.RedisServer,
		RedisPassword: env.SystemENV.RedisPassword,
		RedisDB:       0,
	}
	redis.Init(redisConfig)

	go transection.Init(transection.Config{
		Fequency: time.Second * 5,
	})

	// 初始化ht_pay
	if err := htpay.Init(htpay.Config{
		HT_CALL_BACK_DOMAIN: env.HTENV.HT_CALL_BACK_DOMAIN,
	}); err != nil {
		mlog.Fatal(fmt.Sprintf("ht_pay初始化錯誤: %S", err))
	}

	// go transfer.ProcessTransferRequests("walletTransfer")
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
