package game

import (
	"game_service/instance/apollo"
	"game_service/internal/service"
	"time"
)

// 判斷遊戲商 - 流程用
type GameAgent string

const (
	MT_LIVE    GameAgent = "MT_LIVE"
	MT_LOTTERY GameAgent = "MT_LOTTERY"
)

// type GameAgent Agent
type Agent struct {
	Id         int
	Name       string
	SecretInfo interface{}
}

type MyMemberGameAccountInfo struct {
	service.MemberGameAccountInfo
}

// 遊戲商參數
var (
	MT_LIVE_Agent = Agent{
		Id:         1,
		Name:       "MT_LIVE",
		SecretInfo: "MT_LIVE",
	}
	MT_LOTTERY_Agent = Agent{
		Id:         7,
		Name:       "MT_LOTTERY",
		SecretInfo: "MT_LOTTERY",
	}
	GB_ELEC_Agent = Agent{
		Id:         9,
		Name:       "GB_ELEC",
		SecretInfo: "GB_ELEC",
	}
	WG_SPORT_Agent = Agent{
		Id:         11,
		Name:       "WG_SPORT",
		SecretInfo: "WG_SPORT",
	}
	// DG_LIVE_Agent = Agent{
	// 	Id:         10,
	// 	Name:       "DG_LIVE",
	// 	SecretInfo: "DG_LIVE",
	// }
)

type (
	JoinGameReq struct {
		AgentId int `json:"agent_id" validate:"required"`
	}
	JoinGameRes struct {
		Url string `json:"url"`
	}
	BetRecordReq struct {
		Username  string    `json:"mid"`
		GameID    string    `json:"game_id"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	}
	BetRecordRes struct {
		Reports []apollo.BetListResult `json:"reports"`
	}
	TransferInReq struct {
		AgentId int     `json:"agent_id" validate:"required"`
		Point   float64 `json:"point" `
	}
	TransferInRes struct {
		Status int `json:"status"` // 0: 失敗, 1: 成功
	}
	TransferOutReq struct {
		AgentId int `json:"agent_id" validate:"required"`
		Point   int `json:"point" validate:"required"`
	}
	TransferOutRes struct {
		Status int `json:"status"`
	}
	DonateReq struct {
		AgentId int `json:"agent_id" validate:"required"`
	}
	DonateRes struct {
		Status string `json:"status"`
	}
	LiveRecordReq struct {
		AgentId int `json:"agent_id" validate:"required"`
	}
	LiveRecordRes struct {
		Status string `json:"status"`
	}
	LotteryRecordRer struct {
		AgentId int `json:"agent_id" validate:"required"`
	}
	LotteryRecordRes struct {
		Status string `json:"status"`
	}
	TransferReq struct {
		MemberId      int                   `validate:"required" json:"member_id"`
		Amount        float64               `validate:"required" json:"amount"`
		Type          string                `validate:"required" json:"type"`
		AgentId       int                   `validate:"required" json:"agent_id"`
		GameWalletMap map[int]GameWalletMap `json:"game_wallet_map"`
	}
	GameWalletMap struct {
		Balance      float64 `json:"balance"`
		TransBalance float64 `json:"trans_balance"`
	}

	TransferRes struct {
	}
)
