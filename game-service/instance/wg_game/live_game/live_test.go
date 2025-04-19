package live_game

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

func TestForwardGameAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := ForwardGameReq{
		ApiId:    "f300101",
		Username: "F3_JEEETT07",
	}

	data, err := forwardGameAPI(user, string("https://apiwinneradm.dv2test.com"+LiveForwardGameUri))
	assert.NoError(t, err)
	assert.Equal(t, "JEEETTS", data.ErrorCode)
}

func TestTransferInAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := TransferInReq{
		ApiId:    "f300101",
		Username: "F3_JEEETT07",
		BPoints:  "0",
		Points:   "100000",
		APoints:  "100000",
	}

	data, err := transferInAPI(user, string("https://apiwinneradm.dv2test.com"+LiveTransferInUri))
	assert.NoError(t, err)
	assert.Equal(t, "JEEETTS", data.ErrorCode)
}

func TestBuyListGetAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := BuyListGetApiReq{
		ApiId:     "f300101",
		ProxyName: "f3_pr00101",
		PriType:   "1",
	}

	data, err := buyListGetAPI(user, string("https://apiwinneradm.dv2test.com"+LiveBuyListGetUri))
	assert.NoError(t, err)
	assert.Equal(t, "JEEETTS", data.ErrorCode)
}

// func TestProxyWinloseGetAPI(t *testing.T) {
// 	logInit()

// 	// 測試 createUserAPI 函數
// 	user := ProxyWinloseGetReq{
// 		ApiId:   "f300101",
// 		Date:    "2024-08-29",
// 		PriType: "1",
// 	}

// 	data, err := proxyWinloseGetAPI(user, string("https://apiwinneradm.dv2test.com"+LiveProxyWinloseGetUri))
// 	assert.NoError(t, err)
// 	assert.Equal(t, "JEEETTS", data.ErrorCode)
// }
