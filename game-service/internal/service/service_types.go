package service

import "time"

// db dto
type (
	MemberGameAccountInfo struct {
		MemberId     int    `db:"member_id"`
		UserName     string `db:"username"`
		NickName     string `db:"nick_name"`
		GamePassword string `db:"game_password"`
		GameAgentId  int    `db:"game_agent_id"`
	}
)

type (
	BetRecord struct {
		Id         int       `db:"id"`
		MemberId   int       `db:"member_id"`
		AgentId    int       `db:"agent_id"`
		BetUnique  string    `db:"bet_unique"`
		GameTypeId int       `db:"game_type_id"`
		BetAt      time.Time `db:"bet_at"`
		Bet        float64   `db:"bet"`
		EffectBet  float64   `db:"effect_bet"`
		WinLose    float64   `db:"win_lose"`
		BetInfo    any       `db:"bet_info"`
		IsConfirm  bool      `db:"is_confirm"`
		CreateAt   time.Time `db:"create_at"`
		UpdateAt   time.Time `db:"update_at"`
	}
)

// live response dto
type (
	GetBetRecordReq struct {
		SystemCode string `json:"system_code"` // 系統代碼
		WebId      string `json:"web_id"`      // 站台代碼，即代理唯一識別碼ID
		UserId     string `json:"user_id"`     // 玩家的唯一識別碼
		StartTime  string `json:"start_time"`  // 開始時間,格式:”2018-01-01 02:00:00”
		EndTime    string `json:"end_time"`    // 結束時間,格式:”2018-01-01 02:00:00”，目前結束時間與起始時間不能相差超過 24 小時
		Page       int    `json:"page"`        // 起始頁，正整數，從 1 開始計數
		PageSize   *int   `json:"page_size"`   // 每頁筆數，正整數，最大值100
	}
	GetBetRecordRes struct {
		Code      string              `json:"code"`      // 00000即為成功，其它代碼皆為失敗
		Message   string              `json:"message"`   // 錯誤訊息
		Timestamp int64               `json:"timestamp"` // 時間戳記
		Data      GetBetRecordDataDTO `json:"data"`      // 回傳資料
	}
	GetBetRecordDataDTO struct {
		SystemCode  string    `json:"system_code"`  // 系統代碼
		WebId       string    `json:"web_id"`       // 站台代碼
		CurrentPage int       `json:"current_page"` // 當前頁碼
		TotalPage   int       `json:"total_page"`   // 總頁數
		TotalCount  int       `json:"total_count"`  // 總筆數
		List        []ListDTO `json:"list"`         // 資料列表
	}
	ListDTO struct {
		UserId     string `json:"user_id"`     // 玩家的唯一識別碼
		Sn         string `json:"sn"`          // 注單編號
		GameCode   string `json:"game_code"`   // 遊戲代碼
		TableCode  string `json:"table_code"`  // 桌號
		GameName   string `json:"game_name"`   // 遊戲局號
		PlayCode   string `json:"play_code"`   // 玩法代碼
		PlayName   string `json:"play_name"`   // 玩法名稱
		Odds       string `json:"odds"`        // 賠率
		OrderMoney string `json:"order_money"` // 投注金額
		ValidMoney string `json:"valid_money"` // 有效投注金額
		WinMoney   string `json:"win_money"`   // 中獎金額 (不含本金)
		Profit     string `json:"profit"`      // 輸贏金額
		OrderTime  string `json:"order_time"`  // 下注時間
		SettleTime string `json:"settle_time"` // 結算時間
		SettleDate string `json:"settle_date"` // 歸帳日期
		Status     int    `json:"status"`      // 注單狀態 (0:待結算 ,1:取消 ,2:未中獎,3: 中獎 ,4:和局)
		Currency   string `json:"currency"`    // 幣別
		IP         string `json:"ip"`          // IP
	}
)

// lottery response dto
type (
	BetOrderReq struct {
		Date       int64   `json:"Date"`       // 查詢日期 (Unix timestamp)
		Account    *string `json:"Account"`    // 會員帳號
		Limit      *int    `json:"Limit"`      // 筆數
		LastSerial *string `json:"LastSerial"` // 最後一筆流水序號
		GameTypeId int     `json:"GameTypeId"` // 遊戲類別 ID
		GameId     int     `json:"GameId"`     // 遊戲 ID
		Collect    *int    `json:"Collect"`    // 是否已對獎
	}

	BetOrderRes struct {
		TotalRows int           `json:"totalRows"` // 總筆數
		OverRows  int           `json:"overRows"`  // 剩餘筆數
		Limit     int           `json:"limit"`     // 查詢筆數
		BetTotal  int           `json:"betTotal"`  // 投注總金額
		WinTotal  int           `json:"winTotal"`  // 玩家輸贏總金額
		BetCount  int           `json:"betCount"`  // 注單總數量
		Rows      []BetOrderDTO `json:"rows"`      // 資料列
	}

	BetOrderDTO struct {
		Serial           string  `json:"Serial"`           // 流水編號
		No               string  `json:"No"`               // 注單號
		GameTypeId       string  `json:"GameTypeId"`       // 遊戲類別 ID
		GameId           string  `json:"GameId"`           // 遊戲 ID
		LotteryType      string  `json:"LotteryType"`      // 彩種代碼
		LotteryPlayGroup string  `json:"LotteryPlayGroup"` // 玩法群組編碼
		LotteryPlay      string  `json:"LotteryPlay"`      // 玩法代碼
		Issue            string  `json:"Issue"`            // 期數
		Account          string  `json:"Account"`          // 會員帳號
		Rate             string  `json:"Rate"`             // 匯率
		BetTotal         int     `json:"BetTotal"`         // 下注金額
		BetCount         string  `json:"BetCount"`         // 下注數
		Odds             string  `json:"Odds"`             // 下注時賠率
		Odds2            string  `json:"Odds2"`            // 第二組賠率
		Odds3            string  `json:"Odds3"`            // 第三組賠率
		Winnings         int     `json:"Winnings"`         // 實際中獎金額
		AfterPoints      float64 `json:"AfterPoints"`      // 投注後金額
		OpenResult       *string `json:"OpenResult"`       // 開獎號碼
		OpenResult2      *string `json:"OpenResult2"`      // 第二組開獎號碼
		Type             string  `json:"Type"`             // 下注類別 (0: 正常 / 1:追單)
		Collect          string  `json:"Collect"`          // 是否已對獎 (0:未對獎 1:已對獎)
		BetValid         int     `json:"BetValid"`         // 有效下注
		LotteryContent   string  `json:"LotteryContent"`   // 下注項目
		Status           string  `json:"Status"`           // 狀態 (1:取消 2:未開獎 3:未中獎 4:中獎 5:和局)
		Chase            *string `json:"Chase"`            // 追單 UUID
		DateCurrent      string  `json:"DateCurrent"`      // 本期所屬日期
		DateClosing      string  `json:"DateClosing"`      // 封盤時間
		DateDraw         string  `json:"DateDraw"`         // 開獎時間
		DateCreate       string  `json:"DateCreate"`       // 下注時間
		DateUpdate       string  `json:"DateUpdate"`       // 修改時間
		EffectiveBet     int     `json:"EffectiveBet"`     // 有效下注
	}
)

// live_donate response dto
type (
	DonateRecordReq struct {
		SystemCode string `json:"system_code" binding:"required"`
		WebID      string `json:"web_id" binding:"required"`
		UserID     string `json:"user_id"`
		StartTime  string `json:"start_time" binding:"required"`
		EndTime    string `json:"end_time" binding:"required"`
		Page       int    `json:"page" binding:"required"`
		PageSize   int    `json:"page_size"`
	}
	DonateRecordRes struct {
		Code      string           `json:"code"`
		Message   string           `json:"message"`
		Timestamp int              `json:"timestamp"`
		Data      DonateRecordData `json:"data"`
	}
	DonateRecordData struct {
		SystemCode  string         `json:"system_code"`
		WebID       string         `json:"web_id"`
		CurrentPage int            `json:"current_page"`
		TotalPage   int            `json:"total_page"`
		TotalCount  int            `json:"total_count"`
		List        []DonateRecord `json:"list"`
	}
	DonateRecord struct {
		UserID     string `json:"user_id"`
		SN         string `json:"sn"`
		GameCode   string `json:"game_code"`
		GameName   string `json:"game_name"`
		TableCode  string `json:"table_code"`
		Money      int    `json:"money"`
		Time       string `json:"time"`
		DealerID   string `json:"dealer_id"`
		DealerName string `json:"delaer_name"`
		GiftID     string `json:"gift_id"`
		GiftName   string `json:"gift_name"`
	}
)

type (
	GetDetailedDTO struct {
		Count     string   `json:"count"`
		List      []Record `json:"list"`
		Page      string   `json:"page"`
		Pernumber string   `json:"pernumber"`
		TotalPage string   `json:"totalPage"`
	}

	Record struct {
		AftAmount string `json:"AftAmount"`
		Feature   int    `json:"Feature"`
		MainNo    string `json:"MainNo"`
		MoneyType int    `json:"MoneyType"`
		No        string `json:"No"`
		PreAmount string `json:"PreAmount"`
		RoundEnd  int    `json:"RoundEnd"`
		SubNo     string `json:"SubNo"`
		AwardTime string `json:"awardTime"`
		Bet       string `json:"bet"`
		GameCode  string `json:"gameCode"`
		GameDate  string `json:"gameDate"`
		GameName  string `json:"gameName"`
		GameType  string `json:"gameType"`
		GameID    string `json:"gameid"`
		NetWin    string `json:"netWin"`
		State     string `json:"state"`
		UID       string `json:"uid"`
		ValidBet  string `json:"validbet"`
		Win       string `json:"win"`
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
		ResultH     string  `json:"result_h"`    // 主隊結果
		ResultC     string  `json:"result_c"`    // 客隊結果
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
