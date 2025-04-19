package mt_live

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
	user := service.MemberGameAccountInfo{}

	data, err := createUserAPI(user)
	fmt.Println(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data.Result)
}

func TestEditUserAPI(t *testing.T) {
	logInit()

	user := EditUserReq{
		SystemCode: "THGKE9JKS9",
		WebId:      "GKTEST01",
		UserId:     "Jwtjk03",
		UserName:   "酷哥jj",
	}

	data, err := editUserAPI(user, string("https://zone10.ofa16899.net/api/sapphire/"+EditUserUri))
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data.Result)
}

func TestCheckUserAPI(t *testing.T) {
	logInit()
	user := "test0010"
	data, err := checkUserAPI(user)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data.Result)
}

func TestGetURLTokenAPI(t *testing.T) {
	logInit()
	user := "test0010"

	data, err := getURLTokenAPI(user)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestDepositAPI(t *testing.T) {
	logInit()

	user := "jerrytest0001"
	point := 100000.50

	data, err := depositAPI(user, point)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestWithdrawAPI(t *testing.T) {
	logInit()

	user := "Jwtjk01"
	point := 10000.00

	data, err := withdrawAPI(user, point)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestGetPlayerOnlineListAPI(t *testing.T) {
	logInit()

	req := PlayerOnlineListReq{
		SystemCode: "THGKE9JKS9",
		WebId:      "GKTEST01",
	}

	data, err := getPlayerOnlineListAPI(req, string("https://zone10.ofa16899.net/api/sapphire/"+PlayerOnlineListUri))
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestGetBalanceAPI(t *testing.T) {
	logInit()

	username := "Jwtjk01"

	data, err := getBalanceAPI(username)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestKickoutAPI(t *testing.T) {
	logInit()

	UserId := "Jwtjk01"

	data, err := kickoutAPI(UserId)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestGetBetRecordAPI(t *testing.T) {
	logInit()

	req := GetBetRecordReq{
		SystemCode: "THGKE9JKS9",
		WebId:      "GKTEST01",
		StartTime:  "2024-09-18 12:00:00",
		EndTime:    "2024-09-18 23:59:59",
	}

	data, err := getBetRecordAPI(req, string("https://zone10.ofa16899.net/api/sapphire/"+GetBetRecordUri))
	fmt.Println(data)
	assert.NoError(t, err)
	// assert.Equal(t, 1, data.Data)
}

func TestGetTransactionRecordAPI(t *testing.T) {
	logInit()

	req := GetTransationRecordReq{
		SystemCode: "THGKE9JKS9",
		WebId:      "GKTEST01",
		StartTime:  "2024-09-05 12:00:00",
		EndTime:    "2024-09-05 23:59:59",
		Page:       "1",
	}

	data, err := getTransactionRecordAPI(req, string("https://zone10.ofa16899.net/api/sapphire/"+GetTransationRecordUri))
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestFindTransactionRecordAPI(t *testing.T) {
	logInit()

	req := FindTransationRecordReq{
		SystemCode: "THGKE9JKS9",
		WebId:      "GKTEST01",
		TransferSN: "2k9pcb97d8",
	}

	data, err := findTransationRecordAPI(req, string("https://zone10.ofa16899.net/api/sapphire/"+FindTransationRecordUri))
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestGetDonateRecordAPI(t *testing.T) {
	logInit()

	req := DonateRecordReq{
		SystemCode: "THGKE9JKS9",
		WebID:      "GKTEST01",
		StartTime:  "2024-09-05 12:00:00",
		EndTime:    "2024-09-05 21:59:59",
	}

	data, err := getDonateRecordAPI(req, string("https://zone10.ofa16899.net/api/sapphire/"+GetDonateRecordUri))
	fmt.Println(data)
	assert.NoError(t, err)
	// assert.Equal(t, 1, data.Data)
}

func TestGetUpdateBetRecordAPI(t *testing.T) {
	logInit()

	req := GetUpdateBetRecordReq{
		SystemCode: "THGKE9JKS9",
		WebId:      "GKTEST01",
		StartTime:  "2024-08-30 12:00:00",
		EndTime:    "2024-08-30 21:59:59",
		Page:       1,
		PageSize:   100,
	}

	data, err := getUpdateBetRecordAPI(req, string("https://zone10.ofa16899.net/api/sapphire/"+GetUpdateBetRecordUri))
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}
