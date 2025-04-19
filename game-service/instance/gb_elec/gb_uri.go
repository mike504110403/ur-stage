package gb_elec

type API_PARAMS string

const (
	RegisterUri                     API_PARAMS = "register"           // 註冊
	LoginUri                        API_PARAMS = "login"              // 登入
	DemoLoginUri                    API_PARAMS = "demoLogin"          // 試玩帳號登入
	LobbyLoginUri                   API_PARAMS = "lobbyLogin"         // 大廳登入
	GetMoneyUri                     API_PARAMS = "getMoney"           // 取得用戶當前餘額
	TransferUri                     API_PARAMS = "transfer"           // 攜入/攜出
	GetTransferStateUri             API_PARAMS = "getTransferState"   // 儲值訂單查詢
	GetDetailedUri                  API_PARAMS = "getDetailed"        // 查詢遊戲詳細資料
	GetGameListUri                  API_PARAMS = "getGameList"        // 取得遊戲列表
	GetHistoryTransferUri           API_PARAMS = "getHistoryTransfer" // 取得玩家現金轉帳記錄
	GetOrderStatUri                 API_PARAMS = "getOrderStat"       // 注單統計
	ActivityListsUri                API_PARAMS = "activityLists"      // 玩家活動列表
	ActivityWinnerListUri           API_PARAMS = "activityWinnerList" // 玩家活動中獎名單
	UserGameStateUri                API_PARAMS = "35"                 // 查詢玩家是否有未完成遊戲
	UserGameNoteUri                 API_PARAMS = "36"                 // 玩家未完成的遊戲訊息
	KickOutUri                      API_PARAMS = "37"                 // 踢出玩家
	UserCompensationAmountRecordUri API_PARAMS = "38"                 // 玩家補償金額記錄
)
