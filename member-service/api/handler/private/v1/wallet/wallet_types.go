package wallet

// DetailRes 錢包明細回傳
type (
	DetailRes struct {
		Details []WalletDetail `json:"details"`
	}
	WalletDetail struct {
		Game    string  `json:"game"`
		Balance float64 `json:"balance"`
	}
)

// CenterWalletRes 中心錢包回傳
type (
	CenterWalletRes struct {
		Balance     float64 `json:"balance"`
		LockBalance float64 `json:"lock_balance"`
	}
)

// BringingReq 帶入點數請求
type (
	BringingReq struct {
		GameId string `json:"game_id"`
	}
)
