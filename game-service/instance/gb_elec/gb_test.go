package gb_elec

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
	user := "test0004"
	data, err := createUserAPI(user)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestLoginAPI(t *testing.T) {
	logInit()
	user := "test0004"
	data, err := loginAPI(user)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestDemoLoginAPI(t *testing.T) {
	logInit()
	data, err := demoLoginAPI()
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestLobbyLoginAPI(t *testing.T) {
	logInit()
	user := "test0002"
	data, err := lobbyLoginAPI(user)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestGetMoneyAPI(t *testing.T) {
	logInit()
	user := "test0006"
	data, err := getMoneyAPI(user)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestTransferAPI(t *testing.T) {
	logInit()
	user := "test0031"
	amount := "10000.00"
	data, err := transferAPI(user, amount)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestGetTransferStateAPI(t *testing.T) {
	logInit()
	req := GetTransferStateReq{
		Action:    "getTransferState",
		AppID:     "lVqSt9m5vJlvQBqzwX",
		AppSecret: "OqrDz9091oW1f8VO4k1NA8W9YZp72JpV",
		Uid:       "test0002",
		OrderNo:   "5b961941-7de0-4a67-9285-22e79020a24f",
		SignKey:   "357d4cf555d6b4a18dd1617487bf6bad",
	}
	data, err := getTransferStateAPI(req)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestGetDetailed(t *testing.T) {
	logInit()
	req := GetDetailedReq{
		Action:    "getDetailed",
		AppID:     "lVqSt9m5vJlvQBqzwX",
		AppSecret: "OqrDz9091oW1f8VO4k1NA8W9YZp72JpV",
		StartTime: "2024-09-18 18:00:00",
		EndTime:   "2024-09-18 21:59:59",
		SignKey:   "357d4cf555d6b4a18dd1617487bf6bad",
	}
	data, err := getDetailedAPI(req)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestGetGameList(t *testing.T) {
	logInit()
	// req := GetGameListReq{
	// 	Action:    "getGameList",
	// 	AppID:     "lVqSt9m5vJlvQBqzwX",
	// 	AppSecret: "OqrDz9091oW1f8VO4k1NA8W9YZp72JpV",
	// 	SignKey:   "357d4cf555d6b4a18dd1617487bf6bad",
	// }
	data, err := getGameListAPI()
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestGetHistoryTransfer(t *testing.T) {
	logInit()
	req := GetHistoryTransferReq{
		Action:    "getHistoryTransfer",
		AppID:     "lVqSt9m5vJlvQBqzwX",
		AppSecret: "OqrDz9091oW1f8VO4k1NA8W9YZp72JpV",
		StartTime: "2024-09-18 00:00:00",
		EndTime:   "2024-09-18 23:59:59",
		SignKey:   "357d4cf555d6b4a18dd1617487bf6bad",
	}
	data, err := getHistoryTransferAPI(req)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestGetOrderStat(t *testing.T) {
	logInit()
	req := GetOrderStatReq{
		Action:    "getOrderStat",
		AppID:     "lVqSt9m5vJlvQBqzwX",
		AppSecret: "OqrDz9091oW1f8VO4k1NA8W9YZp72JpV",
		StartTime: "2024-09-18 00:00:00",
		EndTime:   "2024-09-20 23:59:59",
		SignKey:   "357d4cf555d6b4a18dd1617487bf6bad",
	}
	data, err := getOrderStatAPI(req)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestActivityLists(t *testing.T) {
	logInit()
	req := ActivityListsReq{
		Action:    "activityLists",
		AppID:     "lVqSt9m5vJlvQBqzwX",
		AppSecret: "OqrDz9091oW1f8VO4k1NA8W9YZp72JpV",
		Uid:       "test0002",
		SignKey:   "357d4cf555d6b4a18dd1617487bf6bad",
	}
	data, err := getActivityListsAPI(req)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestActivityWinnerList(t *testing.T) {
	logInit()
	req := ActivityWinnerListReq{
		Action:       "34",
		AppID:        "lVqSt9m5vJlvQBqzwX",
		AppSecret:    "OqrDz9091oW1f8VO4k1NA8W9YZp72JpV",
		StartTime:    "2024-09-18 00:00:00",
		EndTime:      "2024-09-19 10:59:59",
		ActivityType: "0",
		SignKey:      "357d4cf555d6b4a18dd1617487bf6bad",
	}
	data, err := getActivityWinnerListAPI(req)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestGetUserGameStateAPI(t *testing.T) {
	logInit()
	user := "test0002"
	data, err := getUserGameStateAPI(user)
	assert.NoError(t, err)
	assert.Equal(t, 1, data)
}

func TestGetUserGameNoteAPI(t *testing.T) {
	logInit()
	uid := "test0002"
	data, err := getUserGameNoteAPI(uid)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestKickOutAPI(t *testing.T) {
	logInit()
	uid := "test0002"
	data, err := kickOutAPI(uid)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}

func TestUserCompensationAmountRecordAPI(t *testing.T) {
	logInit()
	req := GetUserCompensationAmountRecordReq{
		Action:    string(UserCompensationAmountRecordUri),
		AppID:     "lVqSt9m5vJlvQBqzwX",
		AppSecret: "OqrDz9091oW1f8VO4k1NA8W9YZp72JpV",
		StartTime: "2024-09-18 00:00:00",
		EndTime:   "2024-09-20 10:59:59",
		Page:      "1",
		PageSize:  "100",
		SignKey:   "357d4cf555d6b4a18dd1617487bf6bad",
	}
	data, err := getUserCompensationAmountRecordAPI(req)
	assert.NoError(t, err)
	assert.Equal(t, 1, data.Data)
}
