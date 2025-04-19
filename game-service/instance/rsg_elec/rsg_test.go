package rsg_elec

import (
	"fmt"
	"testing"
	"time"

	mlog "github.com/mike504110403/goutils/log"

	"github.com/stretchr/testify/assert"
)

func logInit() {
	// 初始化log
	mlog.Init(mlog.Config{
		EnvMode: mlog.EnvMode("dev"),
		LogType: mlog.LogType("console"),
	})
}

func TestCreateUserAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := "protestjerry002"

	data, err := createUserAPI(user)
	fmt.Println(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestPointInAPI(t *testing.T) {
	logInit()

	user := "testjerry0002"
	amount := 100.00

	data, err := PointIn(user, amount)
	fmt.Println(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestPointOutAPI(t *testing.T) {
	logInit()

	user := "testjerry0002"
	amount := 100.00

	data, err := PointOut(user, amount)
	fmt.Println(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestLoginAPI(t *testing.T) {
	logInit()

	user := "testjerry0002"

	data, err := loginAPI(user)
	fmt.Println(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestGetBalanceAPI(t *testing.T) {
	logInit()

	user := "testjerry0002"

	data, err := GetBalance(user)
	fmt.Println(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestLoginLobbyAPI(t *testing.T) {
	logInit()

	user := "protestjerry0002"

	data, err := LoginLobby(user)
	fmt.Println(data)
	assert.NoError(t, err)
	assert.Equal(t, "", data)
}

func TestKickOutAPI(t *testing.T) {
	logInit()

	user := "testjerry0002"

	data, err := kickOutAPI(user)
	fmt.Println(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestGetGameDetailAPI(t *testing.T) {
	logInit()

	//user := "testjerry0002"

	data, err := getGameDetailAPI(time.Minute)
	fmt.Println(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}
