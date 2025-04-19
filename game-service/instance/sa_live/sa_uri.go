package sa_live

type API_URI string

const (
	RegUserInfoUri                       API_URI = "RegUserInfo"
	VerifyUsernameUri                    API_URI = "VerifyUsername"
	GetUserStatusDVUri                   API_URI = "GetUserStatusDV"
	LoginRequestUri                      API_URI = "LoginRequest"
	KickUserUri                          API_URI = "KickUser"
	CreditBalanceDVUri                   API_URI = "CreditBalanceDV"
	DebitBalanceDVUri                    API_URI = "DebitBalanceDV"
	GetAllBetDetailsForTimeIntervalDVUri API_URI = "GetAllBetDetailsForTimeIntervalDV"
)
