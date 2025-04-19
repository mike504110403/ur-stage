package lottery_game

import (
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
	user := CreateUserApiReq{
		ApiId: "f300101",
		//Username:   "F3_JEEETT07",
		Username:   "jerry01",
		Proxyname:  "f3_pr00101",
		Experience: "n",
	}

	data, err := createUserAPI(user, string("https://apiwinneradm.dv2test.com"+LotteryCreateUserUri))
	assert.NoError(t, err)
	assert.Equal(t, "JEEETTS", data.Username)
}

// func TestCheckUserAPI(t *testing.T) {
// 	logInit()

// 	// 測試 createUserAPI 函數
// 	user := CheckUserReq{
// 		ApiId: "f300101",
// 		//Username: "F3_JEEETT07",
// 		Username: "F3_Mike05",
// 	}

// 	data, err := checkUserAPI(user, string("https://apiwinneradm.dv2test.com"+LotteryCheckUserUri))
// 	assert.NoError(t, err)
// 	assert.Equal(t, "JEEETTS", data.ErrorCode)
// }

// func TestForwardGameAPI(t *testing.T) {
// 	logInit()

// 	// 測試 createUserAPI 函數
// 	user := ForwardGameReq{
// 		ApiId:    "f300101",
// 		Username: "F3_JEEETT07",
// 	}

// 	data, err := forwardGameAPI(user, string("https://apiwinneradm.dv2test.com"+LotteryForwardGameUri))
// 	assert.NoError(t, err)
// 	assert.Equal(t, "JEEETTS", data.ErrorCode)
// }

// func TestKickUserAPI(t *testing.T) {
// 	logInit()

// 	// 測試 createUserAPI 函數
// 	user := KickUserReq{
// 		ApiId:    "f300101",
// 		Username: "F3_JEEETT07",
// 	}

// 	data, err := kickoutAPI(user, string("https://apiwinneradm.dv2test.com"+LotteryKickUserUri))
// 	assert.NoError(t, err)
// 	assert.Equal(t, "JEEETTS", data.ErrorCode)
// }

// func TestTransferInAPI(t *testing.T) {
// 	logInit()

// 	// 測試 createUserAPI 函數
// 	user := TransferInReq{
// 		ApiId:    "f300101",
// 		Username: "F3_JEEETT07",
// 		BPoints:  "0",
// 		Points:   "100000",
// 		APoints:  "100000",
// 	}

// 	data, err := transferInAPI(user, string("https://apiwinneradm.dv2test.com"+LotteryTransferInUri))
// 	assert.NoError(t, err)
// 	assert.Equal(t, "JEEETTS", data.ErrorCode)
// }

// func TestBuyListGetAPI(t *testing.T) {
// 	logInit()

// 	// 測試 createUserAPI 函數
// 	user := BuyListGetApiReq{
// 		ApiId:     "f300101",
// 		ProxyName: "f3_pr00101",
// 		PriType:   "1",
// 	}

// 	data, err := buyListGetAPI(user, string("https://apiwinneradm.dv2test.com"+LotteryBuyListGetUri))
// 	assert.NoError(t, err)
// 	assert.Equal(t, "JEEETTS", data.ErrorCode)
// }
