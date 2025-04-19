package atg_elect

import (
	"fmt"
	"testing"

	mlog "github.com/mike504110403/goutils/log"
)

func logInit() {
	// 初始化log
	mlog.Init(mlog.Config{
		EnvMode: mlog.EnvMode("dev"),
		LogType: mlog.LogType("console"),
	})
}

func TestGetTokenAPI(t *testing.T) {
	logInit()

	data, err := getToken()
	fmt.Println(data, err)
}

func TestRegisterAPI(t *testing.T) {
	logInit()
	user := "test002"
	err := Register(user)
	fmt.Println(err)
}

func TestGetBalanceAPI(t *testing.T) {
	logInit()
	user := "test002"
	data, err := getBalance(user)
	fmt.Println(data, err)
}

func TestGetGameLobbyAPI(t *testing.T) {
	logInit()
	user := "test002"
	data, err := getGameLobby(user)
	fmt.Println(data, err)
}

func TestTransferIn(t *testing.T) {
	logInit()
	user := "test002"
	amount := 100.01
	data, err := transferIn(user, amount)
	fmt.Println(data, err)
}

func TestTransferOut(t *testing.T) {
	logInit()
	user := "test002"
	amount := 100.01
	data, err := transferOut(user, amount)
	fmt.Println(data, err)
}

func TestKickOut(t *testing.T) {
	logInit()
	user := "test002"
	err := kickOut(user)
	fmt.Println(err)
}

func TestBetRecord(t *testing.T) {
	req := BetRecordReq{
		Operator: "Ur_USDT_beta",
		Key:      "ae6bc6d729894f088615aa1e772cdef5",
		SDate:    "2024-10-23 19:03:00",
		EDate:    "2024-10-23 19:08:59",
	}
	data, err := betRecord(req)
	fmt.Println(data, err)
}
