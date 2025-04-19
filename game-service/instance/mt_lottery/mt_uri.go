package mt_lottery

type API_URI string

const (
	CreateUserUri      API_URI = "/Member/Create"          // 建立玩家帳號
	CheckPointUri      API_URI = "/Member/CheckPoint"      // 查詢玩家可用點數
	TransPointUri      API_URI = "/Member/TransPoint"      // 轉入(出)玩家點數
	TransactionLogUri  API_URI = "/Member/TransactionLog"  // 查詢玩家交易記錄
	LoginUri           API_URI = "/Game/Login"             // 玩家登入遊戲
	LogoutUri          API_URI = "/Game/Logout"            // 玩家登出遊戲
	HandicapCheckUri   API_URI = "/Handicap/Check"         // 可用限紅查詢
	HandicapSettingUri API_URI = "/Member/HandicapSetting" // 會員限紅設定
	BetOrderUri        API_URI = "/Bets/BetOrder"          // 查詢指定日期投注記錄
	BetOrderV2Uri      API_URI = "/Bets/BetOrder/v2"       // 查詢指定日期投注記錄 V2
	CheckDateModifyUri API_URI = "/Bets/CheckDateModify"   // 查詢指定日期修改過投注記錄
	ModifyNickNameUri  API_URI = "/Member/ModifyNickName"  // 修改玩家暱稱
	ModifyPasswordUri  API_URI = "/Member/ModifyPassword"  // 修改玩家密碼
)
