package elec_game

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

	data, err := forwardGameAPI(user, string("https://apiwinneradm.dv2test.com"+ElecForwardGameUri))
	assert.NoError(t, err)
	assert.Equal(t, "JEEETTS", data.ErrorCode)
}

func TestTransferInAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := TransferInReq{
		ApiId:    "f300101",
		Username: "F3_JEEETT07",
		Bpoints:  "0",
		Points:   "100000",
		Apoints:  "100000",
	}

	data, err := transferInAPI(user, string("https://apiwinneradm.dv2test.com"+ElecTransferInUri))
	assert.NoError(t, err)
	assert.Equal(t, "JEEETTS", data.ErrorCode)
}

func TestBuyListGetAPI(t *testing.T) {
	logInit()

	// 測試 createUserAPI 函數
	user := BuyListGetApiReq{
		ApiId:     "f300101",
		Proxyname: "f3_pr00101",
		PriType:   "1",
	}

	data, err := buyListGetAPI(user, string("https://apiwinneradm.dv2test.com"+ElecBuyListGetUri))
	assert.NoError(t, err)
	assert.Equal(t, "JEEETTS", data.ErrorCode)
}
