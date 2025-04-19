package game

type (
	TransferReq struct {
		MemberId int     `json:"member_id" validate:"required"`
		Amount   float64 `json:"amount" validate:"required"`
		Type     string  `json:"type" validate:"required"`
		GameId   int     `json:"game_id" validate:"required"`
	}
)
