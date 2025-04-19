package apiret

type ApiRet struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
