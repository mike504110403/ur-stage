package rsg_elec

import (
	encoder "game_service/pkg/encoder"
	"time"
)

type RsgElecSecrect_Info struct {
	Ischeck     bool
	RefreshTime time.Time
	CE          encoder.CashEncryption
}

type Response struct {
	ErrorCode    int    `json:"ErrorCode"`    // 錯誤代碼
	ErrorMessage string `json:"ErrorMessage"` // 錯誤訊息
	Timestamp    int    `json:"Timestamp"`    // 時間戳記
}

type (
	CreatePlayerReq struct {
		SystemCode string `json:"SystemCode"` // 系統代碼(只限英數)
		WebId      string `json:"WebId"`      // 站台代碼(只限英數)
		UserId     string `json:"UserId"`     // 會員惟一識別碼(只限英數)
		Currency   string `json:"Currency"`   // 幣別代碼(請參照代碼表)
	}

	CreatePlayerRes struct {
		Response Response
		Data     struct {
			SystemCode string `json:"SystemCode"` // 系統代碼
			WebId      string `json:"WebId"`      // 站台代碼
			UserId     string `json:"UserId"`     // 會員的唯一識別碼
		} `json:"Data"`
	}
)

type (
	DepositReq struct {
		SystemCode    string  `json:"SystemCode"`
		WebId         string  `json:"WebId"`
		UserId        string  `json:"UserId"`
		TransactionID string  `json:"TransactionID"`
		Currency      string  `json:"Currency"`
		Balance       float64 `json:"Balance"`
	}

	DepositRes struct {
		Response Response
		Data     struct {
			TransactionID        string  `json:"TransactionID"`
			TransactionTime      string  `json:"TransactionTime"`
			UserId               string  `json:"UserId"`
			PointID              string  `json:"PointID"`
			Balance              float64 `json:"Balance"`
			CurrentPlayerBalance float64 `json:"CurrentPlayerBalance"`
		} `json:"Data"`
	}
)

type (
	WithdrawReq struct {
		SystemCode    string  `json:"SystemCode"`
		WebId         string  `json:"WebId"`
		UserId        string  `json:"UserId"`
		TransactionID string  `json:"TransactionID"`
		Currency      string  `json:"Currency"`
		Balance       float64 `json:"Balance"`
	}

	WithdrawRes struct {
		Response Response
		Data     struct {
			TransactionID        string  `json:"TransactionID"`
			TransactionTime      string  `json:"TransactionTime"`
			UserId               string  `json:"UserId"`
			PointID              string  `json:"PointID"`
			Balance              float64 `json:"Balance"`
			CurrentPlayerBalance float64 `json:"CurrentPlayerBalance"`
		} `json:"Data"`
	}
)

type (
	// 請求結構體
	EnterGameReq struct {
		SystemCode string `json:"SystemCode"` // 系統代碼(只限英數)
		WebId      string `json:"WebId"`      // 站台代碼(只限英數)
		UserId     string `json:"UserId"`     // 會員惟一識別碼(只限英數)
		UserName   string `json:"UserName"`   // 會員暱稱
		GameId     int    `json:"GameId"`     // 遊戲代碼
		Currency   string `json:"Currency"`   // 幣別代碼(請參照s代碼表)
		Language   string `json:"Language"`   // 語系代碼(請參照代碼表)
		ExitAction string `json:"ExitAction"` // 離開遊戲時導向特定網址
	}

	// 回應結構體
	EnterGameRes struct {
		Response Response
		Data     struct {
			URL string `json:"URL"` // 進入遊戲的網址
		} `json:"Data"`
	}
)

type (
	// 請求結構體
	GetBalanceReq struct {
		SystemCode string `json:"SystemCode"` // 系統代碼(只限英數)
		WebId      string `json:"WebId"`      // 站台代碼(只限英數)
		UserId     string `json:"UserId"`     // 會員惟一識別碼(只限英數)
		Currency   string `json:"Currency"`   // 幣別代碼(請參照代碼表)
	}

	// 回應結構體
	GetBalanceRes struct {
		Response Response
		Data     struct {
			UserId               string  `json:"UserId"`               // 會員惟一識別碼
			CurrentPlayerBalance float64 `json:"CurrentPlayerBalance"` // 會員當前點數
		} `json:"Data"`
	}
)

type (
	// 請求結構體
	EnterLobbyGameReq struct {
		SystemCode string `json:"SystemCode"` // 系統代碼(只限英數)
		WebId      string `json:"WebId"`      // 站台代碼(只限英數)
		UserId     string `json:"UserId"`     // 會員惟一識別碼(只限英數)
		UserName   string `json:"UserName"`   // 會員暱稱
		Currency   string `json:"Currency"`   // 幣別代碼(請參照s代碼表)
		Language   string `json:"Language"`   // 語系代碼(請參照代碼表)
		//ExitAction string `json:"ExitAction"` // 離開遊戲時導向特定網址
	}

	// 回應結構體
	EnterLobbyGameRes struct {
		Response Response
		Data     struct {
			URL string `json:"URL"` // 進入遊戲的網址
		} `json:"Data"`
	}
)

type (
	// 請求結構體
	KickOutReq struct {
		KickType   int    `json:"KickType"`   // 剔除模式
		SystemCode string `json:"SystemCode"` // 系統代碼(只限英數)
		WebId      string `json:"WebId"`      // 站台代碼(只限英數)
		UserId     string `json:"UserId"`     // 會員惟一識別碼(只限英數)
		GameId     int    `json:"GameId"`     // 若 KickType 不為 3,則填 0
	}

	// 回應結構體
	KickOutRes struct {
		Response Response
		Data     struct {
			UserCount int `json:"UserCount"` // 被剔除的會員數量
		} `json:"Data"`
	}
)

type (
	// 請求結構體
	GetGameDetailReq struct {
		SystemCode string `json:"SystemCode"` // 系統代碼(只限英數)
		WebId      string `json:"WebId"`      // 站台代碼(只限英數)
		GameType   int    `json:"GameType"`   // 遊戲類型(1.老虎機 2.捕魚機)
		TimeStart  string `json:"TimeStart"`  // 開始時間(yyyy-MM-dd HH:mm)
		TimeEnd    string `json:"TimeEnd"`    // 結束時間(yyyy-MM-dd HH:mm)
	}

	// 回應結構體
	GetGameDetailRes struct {
		Response Response
		Data     struct {
			GameDetail []struct {
				Currency            string  `json:"Currency"`            // 幣別代碼
				WebId               string  `json:"WebId"`               // 站台代碼
				UserId              string  `json:"UserId"`              // 會員惟一識別碼
				SequenNumber        int64   `json:"SequenNumber"`        // 遊戲紀錄惟一編號
				GameId              int     `json:"GameId"`              // 遊戲代碼(請參照代碼表)
				SubGameType         int     `json:"SubGameType"`         // 子遊戲代碼(請參照代碼表)
				BetAmt              float64 `json:"BetAmt"`              // 下注(小數點兩位)
				WinAmt              float64 `json:"WinAmt"`              // 贏分(小數點兩位)
				PlayTime            string  `json:"PlayTime"`            // 遊戲時間
				JackpotContribution float64 `json:"JackpotContribution"` // Jackpot 貢獻值(小數點五位)
			} `json:"GameDetail"`
		} `json:"Data"`
	}
)
