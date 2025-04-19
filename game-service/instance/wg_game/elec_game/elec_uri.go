package elec_game

type API_URI string

// elec_uri
const (
	CreateUserUri          API_URI = "/lot_api/CreateUserApi"    // 建立玩家帳號(電子)
	CheckUserUri           API_URI = "/lot_api/CheckUser"        // 查詢玩家帳號(電子)
	EditUserUri            API_URI = "/lot_api/EditUser"         // 修改玩家帳號(電子)
	ElecTransferCheckUri   API_URI = "/elec_api/TransferCheck"   // 確認交易單號(電子)
	ElecForwardGameUri     API_URI = "/elec_api/ForwardGame"     // 進入遊戲(電子)
	ElecTransferInUri      API_URI = "/elec_api/TransferIn"      // 轉入點數(電子)，此為增加電子平台點數。
	ElecTransferOutUri     API_URI = "/elec_api/TransferOut"     // 轉出點數(電子)，此為減少體育平台點數。
	ElecGetPointUser       API_URI = "/elec_api/PointUser"       // 取得點數(電子)
	ElecBuyListGetUri      API_URI = "/elec_api/BuyListGetApi"   // 查詢注單(電子)，查詢要求間隔為1分鍾。資料保留為本月、上月。
	ElecProxyWinloseGetUri API_URI = "/elec_api/ProxyWinloseGet" // 查詢代理線輸贏(電子)
	ElecBuySingleGetApiUri API_URI = "/elec_api/BuySingleGetApi" // 查詢注單開獎號(電子)
)
