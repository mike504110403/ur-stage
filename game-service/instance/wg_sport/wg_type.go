package wg_sport

import (
	"game_service/pkg/encoder"
	"time"
)

type MemberGameAccount struct {
	Username string `db:"username"`
	Nickname string `db:"nickname"`
	Password string `db:"password"`
}

type KV struct {
	Key   string
	Value string
}

type WgSportSecrect_Info struct {
	Ischeck     bool
	RefreshTime time.Time
	WSSI        encoder.WgSportSecretInfo
}

type (
	CreateUserReq struct {
		Prefix     string `json:"prefix"`     // 我方提供之代理碼
		Username   string `json:"username"`   // 會員帳號
		Nickname   string `json:"nickname"`   // 會員暱稱
		Upusername string `json:"upusername"` // 上層帳號
		Sign       string `json:"sign"`       // 於簽名規則說明
	}
	CreateUserRes struct {
		Error_code string `json:"error_code"` // 於共同參數error_code表查詢
		Username   string `json:"username"`   // error_code為OK時，會回傳會員帳號
	}
)

type (
	CheckUserReq struct {
		Prefix   string `json:"prefix"` // 我方提供之代理碼
		Username string `json:"username"`
		Sign     string `json:"sign"` // 於簽名規則說明
	}
	CheckUserRes struct {
		Error_code string `json:"error_code"` // 於共同參數error_code表查詢
	}
)

type (
	EditUserReq struct {
		Prefix   string `json:"prefix"`   // 我方提供之代理碼
		Username string `json:"username"` // 帳號
		// Nickname  string `json:"nickname"`  // 暱稱變更
		// Enabled   int    `json:"enabled"`   // 狀態 1 = 啟用， 2 = 停押(無法投注)， 3 = 停用(無法登入)
		// Limitname string `json:"limitname"` // 限紅名稱，限英數字，例 : limit2
		Sign string `json:"sign"` // 於簽名規則說明
	}

	EditUserRes struct {
		ErrorCode string `json:"error_code"` // 於共同參數error_code表查詢
	}
)

type (
	ForwardGameReq struct {
		Prefix   string  `json:"prefix"`   // 我方提供之代理碼
		Username string  `json:"username"` // 會員帳號
		Mobile   *string `json:"mobile"`   // 手機 : y
		Lang     *string `json:"lang"`     // 繁中 : zh-tw (預設)/簡中 : zh-cn/英文 : en-us/越南文 : vi-vn
		Sign     string  `json:"sign"`     // 於簽名規則說明
	}
	ForwardGameRes struct {
		Error_code string `json:"error_code"` // 於共同參數error_code表查詢
		Url        string `json:"url"`        // 進入遊戲連結，請用get方式打開此連結。同一連結只能使用乙次 ※不支援iframe方式開啟連結※
	}
)

type (
	PointUserReq struct {
		Prefix   string `json:"prefix"` // 我方提供之代理碼
		Username string `json:"username"`
		Sign     string `json:"sign"` // 於簽名規則說明
	}
	PointUserRes struct {
		Error_code string  `json:"error_code"` // 於共同參數error_code表查詢
		Points     float64 `json:"points"`     // 現有點數，error_code為OK時，會回傳點數
	}
)

type (
	TransferInReq struct {
		Prefix   string `json:"prefix"` // 我方提供之代理碼
		Username string `json:"username"`
		Bpoints  string `json:"bpoints"` // 轉前點數，取得點數(PointUser)後，所回傳之點數(points)
		Points   string `json:"points"`  // 轉入點數，限正整數
		Apoints  string `json:"apoints"` // 轉後點數，小數點後三位之後無條件捨去
		Sign     string `json:"sign"`    // 於簽名規則說明
	}
	TransferInRes struct {
		Error_code string `json:"error_code"` // 於共同參數error_code表查詢
		Transferid string `json:"transferid"` // 交易序號:error_code為OK時，會回傳交易單號，用於確認該筆單有無正確入帳

	}
)

type (
	TransferOutReq struct {
		Prefix   string `json:"prefix"` // 我方提供之代理碼
		Username string `json:"username"`
		Bpoints  string `json:"bpoints"` // 轉前點數，取得點數(PointUser)後，所回傳之點數(points)
		Points   string `json:"points"`  // 轉出點數，限正整數
		Apoints  string `json:"apoints"` // 轉後點數，小數點後三位之後無條件捨去
		Sign     string `json:"sign"`    // 於簽名規則說明
	}
	TransferOutRes struct {
		Error_code string `json:"error_code"` // 於共同參數error_code表查詢
		Transferid string `json:"transferid"` // 交易序號:error_code為OK時，會回傳交易單號，用於確認該筆單有無正確入帳

	}
)

type (
	TransferCheckReq struct {
		Prefix     string `json:"prefix"`     // 我方提供之代理碼
		Username   string `json:"username"`   // 會員帳號
		Transferid string `json:"transferid"` // 交易單號，轉點成功後，我方所提供之交易單號
		Sign       string `json:"sign"`       // 於簽名規則說明
	}
	TransferCheckRes struct {
		Error_code string `json:"error_code"` // 於共同參數error_code表查詢
		Support    string `json:"support"`    // 轉點方式，in為轉入點數，out為轉出點數
		Log        string `json:"log"`        // 轉點內容，轉點內容記錄
	}
)

type (
	BuyListGetReq struct {
		Prefix string `json:"prefix"` // 我方提供之代理碼
		Type   string `json:"type"`   // 取得類型，1:帳務日期、2:下注日期、3:更新日期
		Sdate  string `json:"sdate"`  // 查詢起始日，type=1(帳務日期)格式: date("Ymd")/type=2(下注日期)格式:date("Y-m-d H:i:s")/type=3(更新日期)格式:date("Y-m-d H:i:s")
		Edate  string `json:"edate"`  // 查詢結束日，type=1(帳務日期)格式: date("Ymd")/type=2(下注日期)格式:date("Y-m-d H:i:s")/type=3(更新日期)格式:date("Y-m-d H:i:s")，使用type=2或3，起訖時間以區間在24小時內為限
		//Username   *string   `json:"username"`   // 填此欄位，則查詢會員個人
		Pri_type string `json:"pri_type"` //0 全部注單，1 已結算注單
		//Upusername *string   `json:"upusername"` // 上層帳號，僅限代理商帳號，填此欄位，則只查詢該代理商帳務
		Sign string `json:"sign"` // 於簽名規則說明
	}
	BuyListGetRes struct {
		Error_code string    `json:"error_code"` // 於共同參數error_code表查詢
		Data       []DataDTO `json:"data"`       // 注單資料
	}
	DataDTO struct {
		Id          int     `json:"id"`          // 注單流水號，唯一值
		BuyId       string  `json:"buyid"`       // 注單編號，相同buyid為同一張過關單
		Wid         int     `json:"wid"`         // 過關單編號，非過關單 wid= 0，相同wid為同一張過關單
		Main        int     `json:"main"`        // 主要標記，非過關單 main = 1，過關單其中1關 main = 1
		Username    string  `json:"username"`    // 會員帳號
		OrderDate   int     `json:"orderdate"`   // 歸帳日，報表帳務日期
		Gid         int     `json:"gid"`         // 賽事編號
		Gtype       string  `json:"gtype"`       // 球種大項，參照球種大項 gtype 對照表
		Gtypes      string  `json:"gtypes"`      // 球種細項，參照球種細項 gtypes 對照表
		GameTime    string  `json:"gameTime"`    // 賽事時間
		TeamLeague  string  `json:"teamLeague"`  // 聯盟
		TeamH       string  `json:"teamH"`       // 主隊名稱
		TeamC       string  `json:"teamC"`       // 客隊名稱
		BetMsg      string  `json:"betmsg"`      // 賽事盤口
		Odds        string  `json:"odds"`        // 賠率
		ConMsg      string  `json:"conmsg"`      // 投注盤口
		Rtype       string  `json:"rtype"`       // 玩法代碼，參照玩法 rtype 對照表
		Type        string  `json:"type"`        // 投注隊伍，H-主隊、C-客隊
		Strong      int     `json:"strong"`      // 強弱隊伍，1-強隊，0-弱隊
		GtypeName   string  `json:"gtypeName"`   // 投注球種名稱
		RtypeName   string  `json:"rtypeName"`   // 投注玩法名稱
		LineName    string  `json:"lineName"`    // 投注盤別，過關單盤別為空值
		Gold        float64 `json:"gold"`        // 投注金額
		GoldOk      float64 `json:"goldok"`      // 有效投注金額，過關贏的注單，有效投注金額一律為0
		ResultH     int     `json:"result_h"`    // 主隊結果
		ResultC     int     `json:"result_c"`    // 客隊結果
		Result      float64 `json:"result"`      // 輸贏金額
		Wgold       float64 `json:"wgold"`       // 退水金額
		Wltype      string  `json:"wltype"`      // 輸贏類別
		Num         int     `json:"num"`         // 注單注數，非過關單 num = 1，過關單 num >= 2
		IP          string  `json:"IP"`          // IP
		Currency    string  `json:"currency"`    // 幣值
		CurrencyR   string  `json:"currency_r"`  // 匯差
		Status      int     `json:"status"`      // 注單狀態，1=有效注單，其他都為註銷的注單
		Pstatus     int     `json:"pstatus"`     // 派彩狀態，1-未派彩，其他都為已派彩
		Txt         string  `json:"txt"`         // 賽事備註
		InsTime     string  `json:"ins_time"`    // 下注時間
		UpdatedAt   string  `json:"updated_at"`  // 更新時間
		ResultDate  string  `json:"result_date"` // 結算時間
		AfterPoints string  `json:"afterpoints"` // 交易後金額
	}
)

type BuyDetailReq struct {
	Prefix string `json:"prefix"` // 代理碼
	BuyID  string `json:"buyid"`  // 注單編號
	//Lang   string `json:"lang"`   // 語系
	Sign string `json:"sign"` // 簽名
}

type (
	KickUserReq struct {
		Prefix   string `json:"prefix"` // 我方提供之代理碼
		Username string `json:"username"`
		Sign     string `json:"sign"` // 於簽名規則說明
	}
	KickUserRes struct {
		Error_code string `json:"error_code"` // 於共同參數error_code表查詢
	}
)
