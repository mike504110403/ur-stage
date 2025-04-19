package rsg_elec

type API_URI string

const (
	CreateUserUri    API_URI = "/Player/CreatePlayer"
	PointInUri       API_URI = "/Player/Deposit"
	PointOutUri      API_URI = "/Player/Withdraw"
	LoginUri         API_URI = "/Player/GetURLToken"
	GetBalanceUri    API_URI = "/Player/GetBalance"
	LoginLobbyUri    API_URI = "/Player/GetLobbyURLToken"
	KickOutUri       API_URI = "/Player/Kickout"
	GetGameDetailUri API_URI = "/History/GetGameDetail"
)
