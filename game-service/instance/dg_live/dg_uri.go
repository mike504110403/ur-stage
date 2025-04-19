package dg_live

type API_URI string

const (
	SignUpUri     API_URI = "/v2/api/signup"
	LoginUri      API_URI = "/v2/api/login"
	BalanceUri    API_URI = "/v2/api/balance"
	TransferUri   API_URI = "/v2/api/transfer"
	OfflineUri    API_URI = "/v2/api/offline"
	ReportUri     API_URI = "/v2/api/report"
	MarkReportUri API_URI = "/v2/api/markReport"
	OnlineUri     API_URI = "/v2/api/online"
)
