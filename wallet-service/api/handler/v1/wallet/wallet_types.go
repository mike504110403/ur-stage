package wallet

import "time"

// WithdrawGateDto 錢包提款門檻
type (
	WithdrawGateRes struct {
		WithdrawableTimes  int     `json:"withdrawable_times"`  // 可提款次數
		WithdrawableAmount float64 `json:"withdrawable_amount"` // 可提款金額

	}
)

type (
	WithdrawGateFlowRes struct {
		WithdrawEvents []WithdrawEvent `json:"list"`
	}
	WithdrawEvent struct {
		Name     string    `json:"name"`
		Amount   float64   `json:"amount"`
		Gate     float64   `json:"gate"`
		GateFlow float64   `json:"gate_flow"`
		CreateAt time.Time `json:"create_at"`
	}
)

// TransferReq 錢包轉點請求
type (
	TransferReq struct {
		MemberId      int                   `json:"member_id" validate:"required"`
		Amount        float64               `json:"amount"`
		Type          string                `json:"type" validate:"required"`
		AgentId       int                   `json:"agent_id"`
		GameWalletMap map[int]GameWalletMap `json:"game_wallet_map"`
	}
	GameWalletMap struct {
		Balance      float64 `json:"balance"`
		TransBalance float64 `json:"trans_balance"`
	}
)
