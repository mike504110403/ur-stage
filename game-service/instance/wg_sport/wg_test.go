package wg_sport

import (
	"fmt"
	"game_service/internal/service"
	"testing"

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
	user := service.MemberGameAccountInfo{
		UserName: "jerrytest0010",
		NickName: "陳大帥",
	}

	data, err := createUserAPI(&user)
	assert.NoError(t, err)
	assert.Equal(t, "", data.Username)
}

func TestCheckUserAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := "jerrytest0001"

	data, err := checkUserAPI(user)
	assert.NoError(t, err)
	assert.Equal(t, "OK", data.Error_code)
}

func TestEditUserAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := EditUserReq{
		Prefix:   "KAIT",
		Username: "jerrytest0001",
	}
	data, err := editUserAPI(user, string("https://apisport.dv2test.com"+EditUserUri))
	assert.NoError(t, err)
	assert.Equal(t, "OK", data.ErrorCode)
}

func TestForwardGameAPI(t *testing.T) {
	logInit()

	user := "jerrytest0001"
	data, err := forwardGameAPI(user)
	fmt.Printf(data, err)
}

func TestPointUserAPI(t *testing.T) {
	logInit()

	user := "jerrytest0001"

	data, err := PointUserAPI(user)
	assert.NoError(t, err)
	assert.Equal(t, "OK", data)
}

func TestTransferInAPI(t *testing.T) {
	logInit()

	user := "jerrytest0001"
	point := 55.2
	bpoint, err := PointUserAPI(user)
	data, err := transferInAPI(user, bpoint, point)
	assert.NoError(t, err)
	assert.Equal(t, "OK", data.Error_code)
}

func TestTransferOutAPI(t *testing.T) {
	logInit()

	user := "jerrytest0001"
	bpoint, err := PointUserAPI(user)
	data, err := transferOutAPI(user, bpoint)
	assert.NoError(t, err)
	assert.Equal(t, "OK", data)
}

func TestTransferCheckAPI(t *testing.T) {
	logInit()

	user := TransferCheckReq{
		Prefix:     "KAIT",
		Username:   "jerrytest0001",
		Transferid: "KAITWYlnUkOSE71050052542",
	}
	data, err := transferCheckAPI(user, string("https://apisport.dv2test.com"+TransferCheckUri))
	assert.NoError(t, err)
	assert.Equal(t, "OK", data.Error_code)
}

func TestBuyListGetAPI(t *testing.T) {
	logInit()

	user := BuyListGetReq{
		Prefix:   "KAIT",
		Type:     "3",
		Sdate:    "2024-09-25 00:00:01",
		Edate:    "2024-09-25 23:59:00",
		Pri_type: "0",
	}
	data, err := buyListGetAPI(user)
	assert.NoError(t, err)
	assert.Equal(t, "OK", data)
}

// func TestBuyDetailGetAPI(t *testing.T) {
// 	logInit()

// 	user := BuyDetailReq{
// 		Prefix: "KAIT",
// 		BuyID:  "BS00002932",
// 	}
// 	data, err := buyDetailGetAPI(user, string("https://apisport.dv2test.com"+BuyDetailGetUri))
// 	assert.NoError(t, err)
// 	assert.Equal(t, "OK", data)
// }

func TestKickOutAPI(t *testing.T) {
	logInit()

	user := "jerrytest0001"
	data, err := kickOutAPI(user)
	assert.NoError(t, err)
	assert.Equal(t, "OK", data.Error_code)
}
