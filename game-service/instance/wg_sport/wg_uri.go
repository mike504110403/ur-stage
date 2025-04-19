package wg_sport

type API_URI string

const (
	CreateUserUri    API_URI = "/spbet_Api/CreateUser"    // 系統會員建立
	CheckUserUri     API_URI = "/spbet_Api/CheckUser"     // 確認會員創建
	EditUserUri      API_URI = "/spbet_Api/EditUser"      // 系統會員修改
	ForwardGameUri   API_URI = "/spbet_Api/ForwardGame"   // 進入遊戲
	GetPointUserUri  API_URI = "/spbet_Api/PointUser"     // 取得點數
	TransferInUri    API_URI = "/spbet_Api/TransferIn"    // 轉入點數，此為增加體育平台點數。
	TransferOutUri   API_URI = "/spbet_Api/TransferOut"   // 轉出點數，此為減少體育平台點數。
	TransferCheckUri API_URI = "/spbet_Api/TransferCheck" // 確認轉點交易單號
	BuyListGetUri    API_URI = "/spbet_Api/BuyListGet"    // 查詢注單，查詢要求間隔為1分鍾。資料保留為本月、上月。
	BuyDetailGetUri  API_URI = "/spbet_Api/BuyDetailGet"  // 查詢注單明細
	KickUserUri      API_URI = "/spbet_Api/KickUser"      // 踢出會員
)
