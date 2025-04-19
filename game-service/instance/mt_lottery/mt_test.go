package mt_lottery

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
	user := "JEEETTS0"
	Password := "TEST77866"

	data, err := createUserAPI(user, Password)
	assert.NoError(t, err)
	assert.Equal(t, true, data.Rows.Enable)
}

func TestCheckPointAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := service.MemberGameAccountInfo{}
	data, err := checkPointAPI(user)
	fmt.Println(data)
	assert.NoError(t, err)
}

func TestLoginAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := service.MemberGameAccountInfo{}
	data, err := getLoginUrlAPI(user)
	fmt.Println(data)
	assert.NoError(t, err)
}

func TestTransPointAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := service.MemberGameAccountInfo{}
	point := 10000.00
	data, err := transPointAPI(user, point)
	fmt.Println(data)
	assert.NoError(t, err)
}

func TestTransactionLogAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := TransactionLogReq{
		Date: 1725376800,
	}
	data, err := getTransactionLogAPI(user, string("https://aop.ms-168.dev"+TransactionLogUri))
	fmt.Println(data)
	assert.NoError(t, err)
}

func TestLogOutAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數

	Account := "DbaU3-test0010"
	Password := "6fc2acb0-1898-451e-9"

	data, err := logoutAPI(Account, Password)
	fmt.Println(data)
	assert.NoError(t, err)
}

func TestCheckHandicap(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := HandicapCheckReq{
		Account: "DbaU3-JEEETTS0",
	}
	data, err := checkHandicapAPI(user, string("https://aop.ms-168.dev"+HandicapCheckUri))
	fmt.Println(data)
	assert.NoError(t, err)
}

func TestHandicapSetting(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := MemberHandicapSettingReq{
		Account: "DbaU3-JEEETTS0",
	}
	data, err := handicapSettingAPI(user, string("https://aop.ms-168.dev"+HandicapSettingUri))
	fmt.Println(data)
	assert.NoError(t, err)
}

func TestBetOrder(t *testing.T) {
	logInit()
	// 撈取時間戳之後的時間
	user := BetOrderReq{
		Date:       1725580800,
		GameTypeId: 4,
		GameId:     1,
	}
	data, err := getBetOrderAPI(user, string("https://aop.ms-168.dev"+BetOrderUri))
	fmt.Println(data)
	assert.NoError(t, err)
	//assert.Equal(t, true, data.Status)
}

func TestBetOrderV2(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := BetOrderV2Req{
		StartTime: 1725312000,
		EndTime:   1725398399,
		//Account:    "DbaU3-test0010",
		GameTypeId: 4,
		GameId:     1,
		DateType:   1,
	}
	data, err := getBetOrderV2API(user, string("https://aop.ms-168.dev"+BetOrderV2Uri))
	fmt.Println(data)
	assert.NoError(t, err)
	//assert.Equal(t, true, data.Status)
}
