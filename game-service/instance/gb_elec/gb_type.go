package gb_elec

import (
	"game_service/pkg/encoder"
	"time"
)

type GbElecSecrect_Info struct {
	Ischeck     bool
	RefreshTime time.Time
	SC          encoder.SecretConfig
}

type (
	Resquest struct {
		Action    string `json:"action"`    // register
		AppID     string `json:"appID"`     // 應用ID
		AppSecret string `json:"appSecret"` // 應用密鑰
	}
	Response struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}
)

type (
	RegisterReq struct {
		Action    string `json:"action"`    // register
		AppID     string `json:"appID"`     // 應用ID
		AppSecret string `json:"appSecret"` // 應用密鑰
		//Parent    *string `json:"parent"`    // 子代理ID
		Uid string `json:"uid"` // 用户 ID，最大 50 个字符，支持-和_符号
		//Name      *string `json:"name"`      // 用户暱稱，最大 50 个字符，支持-和_符号
		SignKey string `json:"sign_key"` // 签名密钥
		//Debug     *int    `json:"debug"`     // 0=关闭，1=开启，默认为关闭
	}
	RegisterRes struct {
		ReturnCode string   `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string   `json:"returnMsg"`  // 返回消息
		Data       []string `json:"data"`       // 返回数据
	}
)

type (
	LoginReq struct {
		Action    interface{} `json:"action"`    // register
		AppID     string      `json:"appID"`     // 應用ID
		AppSecret string      `json:"appSecret"` // 應用密鑰
		//Parent    *string `json:"parent"`    // 子代理ID
		//Lang      *string `json:"lang"`      // 语言
		Uid      string `json:"uid"`      // 用户 ID，最大 50 个字符，支持-和_符号
		GameCode string `json:"gameCode"` // 游戏代码
		SignKey  string `json:"sign_key"` // 签名密钥
		//Debug     *int    `json:"debug"`     // 0=关闭，1=开启，默认为关闭
		//IsHttps   *int    `json:"is_https"`  // 是否返回 https 游戏地址 1=是 0=否
	}
	LoginRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 关闭类型（点击退出游戏按钮时的操作）；1：跳转到{host}/api/Index/close 地址，APP 内可监听webview地址栏变化并关闭，prod host: api.jav8889.com，stagehost: api.gjylgames.com
	}
	LoginDTO struct {
		Path string `json:"path"` // 登入地址
	}
)

type (
	DemoLoginReq struct {
		Action    interface{} `json:"action"`    // int/string 15/demoLogin
		AppID     string      `json:"appID"`     // 應用ID
		AppSecret string      `json:"appSecret"` // 應用密鑰
		//Lang      *string     `json:"lang"`      // 游戏语言版本 默认为 cn 支持【cn】简体中文、【en】英文、【vn】越南语、【th】泰语、【ko】韩语、【tw】繁体中文、【ja】日语、【pt】葡萄牙语、【esp】西班牙语
		GameCode string `json:"gameCode"` // 游戏代码
		SignKey  string `json:"sign_key"` // 签名密钥
		//Debug     *int        `json:"debug"`     // 0=关闭，1=开启，默认为关闭
		//IsHttps   *int        `json:"is_https"`  // 是否返回 https 游戏地址 1=是 0=否
	}

	DemoLoginRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}

	DemoLoginDTO struct {
		Path string `json:"path"` // 带 token 的登入地址
	}
)

type (
	LobbyLoginReq struct {
		Action    interface{} `json:"action"`    // int/string 25/lobbyLogin
		AppID     string      `json:"appID"`     // 應用ID
		AppSecret string      `json:"appSecret"` // 應用密鑰
		//Lang      *string     `json:"lang"`      // 游戏语言版本 默认为 cn 支持【cn】中文【en】英文【vn】越南语【th】泰语【ko】韩语、【ja】日语、【pt】葡萄牙语、【esp】西班牙语
		Uid     string `json:"uid"`      // 用户 ID
		SignKey string `json:"sign_key"` // 签名密钥
		//Debug     *int        `json:"debug"`     // 0=关闭，1=开启，默认为关闭
		//IsHttps   *int        `json:"is_https"`  // 是否返回 https 地址 1=是 0=否
		//CloseType *int        `json:"closeType"` // 关闭类型
	}

	LobbyLoginRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}

	LobbyLoginDTO struct {
		Path string `json:"path"` // 带 token 的游戏大厅地址
	}
)

type (
	GetMoneyReq struct {
		Action    interface{} `json:"action"`    // int/string 13/getMoney
		AppID     string      `json:"appID"`     // 應用ID
		AppSecret string      `json:"appSecret"` // 應用密鑰
		Uid       string      `json:"uid"`       // 用户 ID
		//Debug     *int        `json:"debug"`     // 0=关闭，1=开启，默认为关闭
		SignKey string `json:"sign_key"` // 签名密钥
	}

	GetMoneyRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}

	GetMoneyDTO struct {
		Amount string `json:"amount"` // 用户余额
	}
)

type (
	TransferReq struct {
		Action    interface{} `json:"action"`    // int/string 14/transfer
		AppID     string      `json:"appID"`     // 應用ID
		AppSecret string      `json:"appSecret"` // 應用密鑰
		Uid       string      `json:"uid"`       // 用户 ID
		OrderNo   string      `json:"orderNo"`   // 商家自定义订单号(建议不超过 36 个字符)
		Amount    string      `json:"amount"`    // 充值金额，正数为充值，负数为回收
		//Debug     *int        `json:"debug"`   // 0=关闭，1=开启，默认为关闭
		//IsForce   *int        `json:"isForce"` // 是否强制 0=否，1=是
		SignKey string `json:"sign_key"` // 签名密钥
	}

	TransferRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}

	TransferDTO struct {
		Uid     string `json:"uid"`     // 用户 ID
		Amount  string `json:"amount"`  // 交易金额
		Balance string `json:"balance"` // 用户余额
		OrderNo string `json:"orderNo"` // 商家自定义订单号
	}
	TransferOutDTO struct {
		Uid     string `json:"uid"`     // 用户 ID
		Amount  string `json:"amount"`  // 交易金额
		Balance string `json:"balance"` // 用户余额
		OrderNo string `json:"orderNo"` // 商家自定义订单号
	}
)

type (
	GetTransferStateReq struct {
		Action    interface{} `json:"action"`    // int/string 24/getTransferState
		AppID     string      `json:"appID"`     // 應用ID
		AppSecret string      `json:"appSecret"` // 應用密鑰
		Uid       string      `json:"uid"`       // 用户 ID
		OrderNo   string      `json:"orderNo"`   // 商家自定义订单号
		SignKey   string      `json:"sign_key"`  // 签名密钥
		//Debug     *int        `json:"debug"`     // 0=关闭，1=开启，默认为关闭
	}

	GetTransferStateRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}

	GetTransferStateDTO struct {
		Uid         string `json:"uid"`         // 用户 ID
		Amount      string `json:"amount"`      // 交易金额
		OrderNo     string `json:"orderNo"`     // 商家自定义订单号
		State       string `json:"state"`       // 订单状态 1=成功，0=失败
		OrderTime   string `json:"orderTime"`   // 订单时间
		PreAmount   string `json:"preAmount"`   // 交易前余额
		AfterAmount string `json:"afterAmount"` // 交易后余额
	}
)

type (
	GetDetailedReq struct {
		Action    interface{} `json:"action"`    // int/string 17/getDetailed
		AppID     string      `json:"appID"`     // 應用ID
		AppSecret string      `json:"appSecret"` // 應用密鑰
		StartTime string      `json:"starttime"` // 开始时间("Y-m-d H:i:s")
		EndTime   string      `json:"endtime"`   // 结束时间("Y-m-d H:i:s")
		//Page      int         `json:"page"`      // 页码
		//PerNumber int         `json:"pernumber"` // 每页分页数量，默认为 100，最大支持5000
		//Debug     *int        `json:"debug"`     // 0=关闭，1=开启，默认为关闭
		SignKey string `json:"sign_key"` // 签名密钥
	}

	GetDetailedRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
		Page       string      `json:"page"`       // 当前页码
		PerNumber  string      `json:"pernumber"`  // 每页显示数量
		Count      string      `json:"count"`      // 总数量
		TotalPage  string      `json:"totalPage"`  // 总页数
	}

	GetDetailedDataDto struct {
		Count     string               `json:"count"`
		List      []GetDetailedListDTO `json:"list"`
		Page      string               `json:"page"`
		Pernumber string               `json:"pernumber"`
		TotalPage string               `json:"totalPage"`
	}

	GetDetailedListDTO struct {
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
	GetGameListReq struct {
		Action    interface{} `json:"action"`    // int/string 20/getGameList
		AppID     string      `json:"appID"`     // 應用ID
		AppSecret string      `json:"appSecret"` // 應用密鑰
		SignKey   string      `json:"sign_key"`  // 签名密钥
		//Lang      *string     `json:"lang"`      // 语言版本 默认为 cn 支持【cn】中文、【en】英文、【vn】越南语、【th】泰语、【ko】韩语、【ja】日语、【pt】葡萄牙语、【esp】西班牙语
		// Debug     *int        `json:"debug"`     // 0=关闭，1=开启，默认为关闭
	}

	GetGameListRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}

	GetGameListDTO []struct {
		GameCode string `json:"gameCode"` // 游戏代码（游戏唯一标识）
		GameID   string `json:"gameId"`   // 游戏 ID（非唯一值，新游戏可能会使用下架游戏的 gameId）
		GameType string `json:"gameType"` // 游戏类型 ID：1 街机游戏，2 捕鱼游戏，3 牌类游戏，4 电子游戏
		Name     string `json:"name"`     // 游戏名称
		Image    string `json:"image"`    // 游戏图片地址（新版接口不再返回信息，字段将保留）
	}
)

type (
	GetHistoryTransferReq struct {
		Action    interface{} `json:"action"`    // int/string 21/getHistoryTransfer
		AppID     string      `json:"appID"`     // 應用ID
		AppSecret string      `json:"appSecret"` // 應用密鑰
		//Uid       *string     `json:"uid"`       // 用户 ID
		StartTime string `json:"starttime"` // 开始时间("Y-m-d H:i:s")
		EndTime   string `json:"endtime"`   // 结束时间("Y-m-d H:i:s")
		SignKey   string `json:"sign_key"`  // 签名密钥
		//Debug     *int        `json:"debug"`     // 0=关闭，1=开启，默认为关闭
	}

	GetHistoryTransferRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}

	GetHistoryTransferDTO []struct {
		No      int    `json:"No"`      // 订单号
		Type    int    `json:"type"`    // 1=充值，2=回收，3=奖金
		Amount  string `json:"amount"`  // 金额
		Balance string `json:"balance"` // 余额
		LogTime string `json:"LogTime"` // 转账时间
		Uid     string `json:"uid"`     // 用户 ID
		OrderNo string `json:"orderNo"` // 商家自定义的订单号
	}
)

type (
	GetOrderStatReq struct {
		Action    interface{} `json:"action"`    // int/string 31/getOrderStat
		AppID     string      `json:"appID"`     // 應用ID
		AppSecret string      `json:"appSecret"` // 應用密鑰
		StartTime string      `json:"starttime"` // 开始时间("Y-m-d H:i:s")
		EndTime   string      `json:"endtime"`   // 结束时间("Y-m-d H:i:s")
		SignKey   string      `json:"sign_key"`  // 签名密钥
		//GameCode  *string     `json:"gameCode"`  // 按游戏代码查询单个游戏数据
		//GameType  *int        `json:"gameType"`  // 按游戏类型查询数据 目前支持2:捕鱼类
	}

	GetOrderStatRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}

	GetOrderStatDTO struct {
		BetOnGold     float64 `json:"BetonGold"`     // 注单金额
		WinGold       float64 `json:"WinGold"`       // 输赢金额
		PaymentAmount float64 `json:"paymentAmount"` // 派彩金额
		OrderSum      int     `json:"orderSum"`      // 注单数量
		UserSum       int     `json:"userSum"`       // 下注用户数
	}
)

type (
	ActivityListsReq struct {
		Action    interface{} `json:"action"`    // int/string 33/activityLists
		AppID     string      `json:"appID"`     // 應用ID
		AppSecret string      `json:"appSecret"` // 應用密鑰
		Uid       string      `json:"uid"`       // 用户 ID
		SignKey   string      `json:"sign_key"`  // 签名密钥
	}

	ActivityListsRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}

	ActivityListsDTO []struct {
		Name    string `json:"name"`    // 活动名称
		LinkUrl string `json:"linkUrl"` // 活动链接
		Icon    string `json:"icon"`    // 活动图标地址
		Type    string `json:"type"`    // 活动类型 1 龙虎榜 2 红包 3 足球
	}
)

type (
	ActivityWinnerListReq struct {
		Action       string `json:"action"`        // int/string 34/activityWinnerList
		AppID        string `json:"appID"`         // 應用ID
		AppSecret    string `json:"appSecret"`     // 應用密鑰
		StartTime    string `json:"start_time"`    // 开始时间
		EndTime      string `json:"end_time"`      // 结束时间
		ActivityType string `json:"activity_type"` // 活动类型 0 全部 2 排行榜 4 红包 5 免费旋转
		// Page         int         `json:"page"`          // 页数 默认 1
		// PageSize     int         `json:"page_size"`     // 每页条数 默认 100
		SignKey string `json:"sign_key"` // 签名密钥
	}

	ActivityWinnerListRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}

	ActivityWinnerListData struct {
		Count     int                     `json:"count"`     // 总数量
		TotalPage int                     `json:"totalPage"` // 总页数
		Page      int                     `json:"page"`      // 当前页数
		PageSize  int                     `json:"pageSize"`  // 每页条数
		List      []ActivityWinnerListDTO `json:"list"`      // 中奖名单
	}

	ActivityWinnerListDTO struct {
		Prize    string `json:"prize"`    // 中奖金额
		Time     string `json:"time"`     // 时间
		Accounts string `json:"Accounts"` // 用户账号
		ID       string `json:"id"`       // 唯一 id
		ActID    string `json:"act_id"`   // 活动 id
		Type     int    `json:"type"`     // 活动类型 2 排行榜 4 红包 5 免费旋转
		No       string `json:"No"`       // 订单号(除免费旋转外其他皆为空字符串)
	}
)

type (
	UserGameStateReq struct {
		Action    string `json:"action"`    // int 35
		AppID     string `json:"appID"`     // 應用ID
		AppSecret string `json:"appSecret"` // 應用密鑰
		Uid       string `json:"uid"`       // 用户账号
		SignKey   string `json:"sign_key"`  // 签名密钥
	}

	UserGameStateRes struct {
		ReturnCode string      `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string      `json:"returnMsg"`  // 返回消息
		Data       interface{} `json:"data"`       // 返回数据
	}

	UserGameStateData struct {
		IsUserHasUnfinishedGames bool `json:"isUserHasUnfinishedGames"` // 用户是否有未完成的游戏
	}
)

type (
	UserGameNoteReq struct {
		Action    string `json:"action"`    // int 36
		AppID     string `json:"appID"`     // 應用ID
		AppSecret string `json:"appSecret"` // 應用密鑰
		Uid       string `json:"uid"`       // 用户账号
		SignKey   string `json:"sign_key"`  // 签名密钥
	}

	UserGameNoteRes struct {
		ReturnCode string          `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string          `json:"returnMsg"`  // 返回消息
		Data       UserGameNoteDTO `json:"data"`       // 返回数据
	}

	UserGameNoteDTO struct {
		GameCode   string `json:"gameCode"`   // 游戏代码
		FreeTimes  string `json:"FreeTimes"`  // 免费游戏次数
		BonusTimes string `json:"BonusTimes"` // 小游戏次数
		Special    string `json:"Special"`    // 是否有特殊模式 (1=有 0=无)
	}
)

type (
	KickOutReq struct {
		Action    string `json:"action"`    // int 37
		AppID     string `json:"appID"`     // 應用ID
		AppSecret string `json:"appSecret"` // 應用密鑰
		Uid       string `json:"uid"`       // 用户账号
		SignKey   string `json:"sign_key"`  // 签名密钥
	}

	KickOutRes struct {
		ReturnCode string   `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string   `json:"returnMsg"`  // 返回消息
		Data       []string `json:"data"`       // 返回数据
	}
)

type (
	GetUserCompensationAmountRecordReq struct {
		Action    string `json:"action"`     // int 38
		AppID     string `json:"appID"`      // 應用ID
		AppSecret string `json:"appSecret"`  // 應用密鑰
		StartTime string `json:"start_time"` // 开始时间
		EndTime   string `json:"end_time"`   // 结束时间
		Page      string `json:"page"`       // 页数（默认 1）
		PageSize  string `json:"pageSize"`   // 每页条数（默认 100）
		SignKey   string `json:"sign_key"`   // 签名密钥
	}

	GetUserCompensationAmountRecordRes struct {
		ReturnCode string                           `json:"returnCode"` // 状态码 0000=成功 具体参考错误代码
		ReturnMsg  string                           `json:"returnMsg"`  // 返回消息
		Data       UserCompensationAmountRecordData `json:"data"`       // 返回数据
	}

	UserCompensationAmountRecordData struct {
		Count     int                  `json:"count"`     // 总数
		TotalPage int                  `json:"totalPage"` // 总页数
		Page      int                  `json:"page"`      // 当前页数
		PageSize  int                  `json:"pageSize"`  // 每页条数
		List      []UserTransactionDTO `json:"list"`      // 交易记录列表
	}

	UserTransactionDTO struct {
		Accounts string `json:"Accounts"` // 用户名
		LogTime  string `json:"LogTime"`  // 记录时间
		Amount   string `json:"Amount"`   // 变动金额
	}
)
