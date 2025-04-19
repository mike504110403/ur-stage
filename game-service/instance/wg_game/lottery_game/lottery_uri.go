package lottery_game

type API_URI string

// lottery_uri
const (
	LotteryCreateUserUri       API_URI = "/lot_api/CreateUserApi"    // 系統會員建立
	LotteryCheckUserUri        API_URI = "/lot_api/CheckUser"        // 確認會員創建
	LotteryEditUserUri         API_URI = "/lot_api/EditUser"         // 系統會員修改
	LotteryTransferCheckUri    API_URI = "/lot_api/TransferCheck"    // 確認交易單號(彩票)
	LotteryForwardGameUri      API_URI = "/lot_api/ForwardGame"      // 進入遊戲(彩票)
	LotteryGetPointUser        API_URI = "/lot_api/PointUser"        // 取得點數(彩票)
	LotteryTransferInUri       API_URI = "/lot_api/TransferIn"       // 轉入點數(彩票)，此為增加彩票平台點數。
	LotteryTransferOutUri      API_URI = "/lot_api/TransferOut"      // 轉出點數(彩票)，此為減少體育平台點數。
	LotteryBuyListGetUri       API_URI = "/lot_api/BuyListGetApi"    // 查詢注單(彩票)，查詢要求間隔為1分鍾。資料保留為本月、上月。
	LotteryProxyWinloseGetUri  API_URI = "/lot_api/ProxyWinloseGet"  // 查詢代理線輸贏(彩票)
	LotteryUserWinloseGetUri   API_URI = "/lot_api/UserWinloseGet"   // 查詢會員輸贏(彩票)
	LotteryBuySingleGetApiUri  API_URI = "/lot_api/BuySingleGetApi"  // 查詢注單開獎號(彩票)
	LotteryKickUserUri         API_URI = "/lot_api/KickUser"         // 踢出會員(彩票&真人&電子)
	LotteryGiftListGetApiv2Uri API_URI = "/lot_api/GiftListGetApiv2" // 查詢彩票打賞注單(新)
)
