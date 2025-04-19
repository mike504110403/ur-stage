package atg_elect

type GetTokenRes struct {
	Status string `json:"status"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data"`
}

type (
	RegisterReq struct {
		Username string `json:"username"`
	}
	RegisterRes struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
)

type (
	BalanceReq struct {
		Username string `json:"username"`
	}
	BalanceRes struct {
		Status string `json:"status"`
		Data   struct {
			Balance string `json:"balance"`
			Status  string `json:"status"`
			Playing string `json:"playing"`
		} `json:"data"`
	}
)

type (
	GameLobbyReq struct {
		Username string `json:"username"`
	}
	GameLobbyRes struct {
		Status string `json:"status"`
		Data   struct {
			Url string `json:"url"`
		} `json:"data"`
	}
)

type (
	TransferReq struct {
		Username   string  `json:"username"`
		Balance    float64 `json:"balance"`
		Action     string  `json:"action"`
		TransferId string  `json:"transferId"`
	}
	TransferRes struct {
		Status string `json:"status"`
		Data   struct {
			Balance string `json:"balance"`
		} `json:"data"`
	}
)

type (
	KickoutReq struct {
		Username string `json:"username"`
	}
	KickoutRes struct {
		Status string `json:"status"`
	}
)

// BetQueryReq 結構體
type (
	BetRecordReq struct {
		Operator string `json:"Operator"`
		Key      string `json:"Key"`
		SDate    string `json:"SDate"`
		EDate    string `json:"EDate"`
	}
	BetRecordRes []struct {
		GameProvider  string  `json:"gameprovider"`
		MemberName    string  `json:"membername"`
		GameName      string  `json:"gamename"`
		BettingID     string  `json:"bettingId"`
		BettingCode   int     `json:"bettingcode"`
		BettingDate   string  `json:"bettingdate"`
		GameID        int     `json:"gameid"`
		GameCode      string  `json:"gamecode"`
		RoundNo       *string `json:"roundno"`
		Result        *string `json:"result"`
		Bet           *string `json:"bet"`
		WinLoseResult string  `json:"winloseresult"`
		TotalStake    string  `json:"totalstake"`
		BettingAmount float64 `json:"bettingamount"`
		ValidBet      string  `json:"validbet"`
		WinLoseAmount float64 `json:"winloseamount"`
		Balance       *string `json:"balance"`
		Currency      string  `json:"currency"`
		IsFree        string  `json:"isfree"`
		IsFeature     string  `json:"isfeature"`
		Handicap      *string `json:"handicap"`
		Status        string  `json:"status"`
		GameCategory  string  `json:"gamecategory"`
		SettleDate    string  `json:"settledate"`
		Remark        *string `json:"remark"`
		BetInfo       *string `json:"betinfo"`
		BetData       string  `json:"betdata"`
		ReplayURL     string  `json:"replayurl"`
	}
)
