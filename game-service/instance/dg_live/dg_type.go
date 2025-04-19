package dg_live

import "time"

type DGLiveSecrect_Info struct {
	Ischeck     bool
	RefreshTime time.Time
	SE          Secret
}
type Secret struct {
	URL         string `json:"url"`
	API_ACCOUNT string `json:"api_account"`
	API_SECRET  string `json:"api_secret"`
}

// 註冊回應
type (
	SignUpRes struct {
		CodeId int    `json:"codeId"`
		Msg    string `json:"msg"`
	}
)

// 踢出回應
type (
	OfflineReq struct {
		OffLineP struct {
			List []int `json:"list"`
		} `json:"offLineP"`
	}
	OfflineRes struct {
		CodeId int    `json:"codeId"`
		Msg    string `json:"msg"`
	}
)

// 帳戶餘額回應
type (
	BalanceReq struct {
		UserName string `json:"username"`
	}
	BalanceRes struct {
		CodeId   int     `json:"codeId"`
		Msg      string  `json:"msg"`
		Username string  `json:"username"`
		Balance  float64 `json:"balance"`
	}
)

// 轉帳回應
type (
	TransferRes struct {
		CodeId   int     `json:"codeId"`
		Msg      string  `json:"msg"`
		UserName string  `json:"username"`
		Amount   float64 `json:"amount"`
		Balance  float64 `json:"balance"`
		Serial   string  `json:"serial"`
	}
)

// 登入回應
type (
	LoginRes struct {
		CodeId  int      `json:"codeId"`
		Msg     string   `json:"msg"`
		Token   string   `json:"token"`
		Domains string   `json:"domains"`
		List    []string `json:"list"`
	}
)

// 注單紀錄
type (
	ReportRes struct {
		CodeId     int      `json:"codeId"`
		Msg        string   `json:"msg"`
		ReportList []Report `json:"list"`
	}

	Report struct {
		Id            int     `json:"id"`
		TableId       int     `json:"tableId"`
		ShoeId        int     `json:"shoeId"`
		PlayId        int     `json:"playId"`
		LobbyId       int     `json:"lobbyId"`
		GameType      int     `json:"gameType"`
		GameId        int     `json:"gameId"`
		BetTime       string  `json:"betTime"`
		CalTime       string  `json:"calTime"`
		WinLose       float64 `json:"winOrLoss"`
		BalanceBefore float64 `json:"balanceBefore"`
		BetPoints     float64 `json:"betPoints"`
		AvailableBet  float64 `json:"availableBet"`
		UserName      string  `json:"userName"`
		Result        string  `json:"result"`
		BetDetail     string  `json:"betDetail"`
		Ip            string  `json:"ip"`
		Ext           string  `json:"ext"`
		IsRevocation  int     `json:"isRevocation"`
		ParentBetId   int     `json:"parentBetId"`
		CurrencyId    int     `json:"currencyId"`
		DeviceType    int     `json:"deviceType"`
		TransferRes   string  `json:"transfers"`
	}
)

// 注單標記確認
type (
	MarkReportRes struct {
		CodeId int    `json:"codeId"`
		Msg    string `json:"msg"`
	}
)

// 在線會員
type (
	OnlineRes struct {
		List []OnlineMember `json:"list"`
	}

	OnlineMember struct {
		UserName string  `json:"username"`
		NickName string  `json:"nickname"`
		Ip       string  `json:"ip"`
		Device   string  `json:"device"`
		Login    string  `json:"login"`
		MemberId int     `json:"memberId"`
		Balance  float64 `json:"balance"`
	}
)
