package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
	mlog "github.com/mike504110403/goutils/log"
)

// config 管理所有通用參數，包含各種來源：env, json file等

var version = "v0.0.0"

// startTime : startTime
var startTime time.Time = time.Now()

var _env = &Env{}

func Init(initVersion string) {
	version = initVersion
}

func EnvInit() {
	if err := env.Parse(_env); err != nil {
		mlog.Fatal(fmt.Sprintf("env.Parse: %s", err))
	}

	if err := validator.New().Struct(_env); err != nil {
		mlog.Fatal(fmt.Sprintf("env.validator: %v", err))
	}
}

func GetVersion() string {
	return version
}

// GetStartTime : 取得程式運行開始時間點
func GetStartTime() time.Time {
	return startTime
}

func GetENV() *Env {
	return _env
}
func GetSystemENV() SystemENV {
	return _env.SystemENV
}

func GetFunctionalENV() FunctionalENV {
	return _env.FunctionalENV
}

func GetProtocolENV() ProtocolENV {
	return _env.ProtocolENV
}
