package apollo

type API_URI string

const (
	AccountExistUri          API_URI = "/user/exist"
	AccountRegisterUri       API_URI = "/user/register"
	DepositUri               API_URI = "/user/deposit"
	GetQuotaUri              API_URI = "/user/getquota"
	AccountLoginUri          API_URI = "/user/login"
	GetBetReportUri          API_URI = "/report/betreport"
	AccountChangePasswordUri API_URI = "/user/changepassword"
	GetAccountBookUri        API_URI = "/report/accountbook"
	GetPeriodResultUri       API_URI = "/game/result"
	GetRealReportUri         API_URI = "/report/realreport"
	CheckRefnoUri            API_URI = "/user/refno"
	GetTotalRealGoldUri      API_URI = "/report/totalrealgold"
	GetFhDetailsUri          API_URI = "/report/fhdetails"
	GetRecReportUri          API_URI = "/report/recreport"
	AccountLogoutUri         API_URI = "/user/logout"
	CheckOnlineStatusUri     API_URI = "/user/onlinestatus"
	GetLineStatusUri         API_URI = "/user/getline"
	PreDepositUri            API_URI = "/user/pre_deposit"
	CheckDepositUri          API_URI = "/user/ck_deposit"
)
