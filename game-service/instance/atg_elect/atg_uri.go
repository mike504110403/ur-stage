package atg_elect

type API_URI string

const (
	ProviserId       API_URI = "4/"
	TokenUri         API_URI = "token"
	RegisterUri      API_URI = "register"
	GameProvidersUri API_URI = "game-providers/"
	BalanceUri       API_URI = "balance"
	GameLobbyUri     API_URI = "lobby"
	PlayUri          API_URI = "play/"
	TransactionUri   API_URI = "transaction"
)
