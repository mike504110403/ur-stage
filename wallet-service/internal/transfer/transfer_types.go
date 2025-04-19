package transfer

type TransferType string

const (
	MainToSub   TransferType = "main-to-sub"
	SubToMain   TransferType = "sub-to-main"
	GameToSub   TransferType = "game-to-sub"
	BringOut    TransferType = "bring-out"
	BringIn     TransferType = "bring-in"
	WithDraw    TransferType = "withdraw"
	Transection TransferType = "transection"
)

type TransferRequest struct {
	MemberId int          `json:"member_id"`
	Amount   float64      `json:"amount"`
	Type     TransferType `json:"type"`
	GameId   int          `json:"game_id"`
}

type MemberWallet struct {
	MainBalance  float64         `json:"main_balance"`
	GameBalances map[int]float64 `json:"game_balances"`
}
