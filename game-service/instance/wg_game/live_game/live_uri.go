package live_game

type API_URI string

// live_uri
const (
	CreateUserUri           API_URI = "/lot_api/CreateUserApi"     // 建立玩家帳號(真人)
	CheckUserUri            API_URI = "/lot_api/CheckUser"         // 查詢玩家帳號(真人)
	EditUserUri             API_URI = "/lot_api/EditUser"          // 修改玩家帳號(真人)
	LiveTransferCheckUri    API_URI = "/live_api/TransferCheck"    // 確認交易單號(真人)
	LiveForwardGameUri      API_URI = "/live_api/ForwardGame"      // 進入遊戲(真人)
	LiveTransferInUri       API_URI = "/live_api/TransferIn"       // 轉入點數(真人)，此為增加真人平台點數
	LiveTransferOutUri      API_URI = "/live_api/TransferOut"      // 轉出點數(真人)，此為減少真人平台點數。
	LiveGetPointUser        API_URI = "/live_api/PointUser"        // 取得點數(真人)
	LiveBuyListGetUri       API_URI = "/live_api/BuyListGetApi"    // 查詢注單(真人)，查詢要求間隔為1分鍾。資料保留為本月、上月。
	LiveProxyWinloseGetUri  API_URI = "/live_api/ProxyWinloseGet"  // 查詢代理線輸贏(真人)
	LiveUserWinloseGetUri   API_URI = "/live_api/UserWinloseGet"   // 查詢會員輸贏(真人)
	LiveBuySingleGetApiUri  API_URI = "/live_api/BuySingleGetApi"  // 查詢注單開獎號(真人)
	LiveGetSuperLotteryUri  API_URI = "/live_api/GetSuperLottery"  // 取得超級彩票(GET)
	LiveSuperLotteryListUri API_URI = "/live_api/SuperLotteryList" // 真人超級彩金(小彩金)發放紀錄
	LiveRewardListGetApiUri API_URI = "/live_api/RewardListGetApi" // 查詢真人打賞注單
)
