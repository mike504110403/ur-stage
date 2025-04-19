package mt_live

type API_URI string

const (
	CreateUserUri           API_URI = "Player/CreateUser"
	EditUserUri             API_URI = "Player/EditUser"
	CheckUserUri            API_URI = "Player/CheckUser"
	GetURLTokenUri          API_URI = "Player/GetURLToken"
	DepositUri              API_URI = "Player/Deposit"
	WithdrawUri             API_URI = "Player/Withdraw"
	GetBalanceUri           API_URI = "Player/GetBalance"
	KickoutUri              API_URI = "Player/Kickout"
	PlayerOnlineListUri     API_URI = "Player/PlayerOnlineList"
	GetBetRecordUri         API_URI = "Report/GetBetRecord"
	GetTransationRecordUri  API_URI = "Report/GetTransactionRecord"
	GetDonateRecordUri      API_URI = "Report/GetDonateRecord"
	FindTransationRecordUri API_URI = "Report/FindTransactionRecord"
	GetUpdateBetRecordUri   API_URI = "Report/GetUpdateBetRecord"
)
