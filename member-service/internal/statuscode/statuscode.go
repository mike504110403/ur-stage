package statuscode

import (
	"gitlab.com/gogogo2712128/common_moduals/apiprotocol"
)

// ServiceNo = 100表示該服務尚未正常地初始化
var cfg = Config{
	ServiceNo: 100,
}

type Config struct {
	ServiceNo int
}

func Init(initCfg Config) {
	cfg = initCfg
	apiprotocol.Init(apiprotocol.Config{
		ServiceNo: cfg.ServiceNo,
	})

	retStatusList := map[apiprotocol.Code]apiprotocol.RetStatusContent{
		// API 1000~1099 一般API狀態
		apiprotocol.PadServiceNo(apiprotocol.Code(cfg.ServiceNo), Success1000):         {Msg: "Success"},
		apiprotocol.PadServiceNo(apiprotocol.Code(cfg.ServiceNo), API_BodyParseFail):   {Msg: "Body解析錯誤"},
		apiprotocol.PadServiceNo(apiprotocol.Code(cfg.ServiceNo), API_ReqValidateFail): {Msg: "Req 驗證錯誤"},
		// Auth 1100~1199 會員權、登入相關狀態
		apiprotocol.PadServiceNo(apiprotocol.Code(cfg.ServiceNo), Auth_LoginFail):     {Msg: "登入失敗"},
		apiprotocol.PadServiceNo(apiprotocol.Code(cfg.ServiceNo), Auth_RegisterFail):  {Msg: "註冊失敗"},
		apiprotocol.PadServiceNo(apiprotocol.Code(cfg.ServiceNo), Auth_NoAccount):     {Msg: "無此帳號"},
		apiprotocol.PadServiceNo(apiprotocol.Code(cfg.ServiceNo), Auth_GetMemberFail): {Msg: "取得會員資料失敗"},
		// jwt 1200~1299 jwt相關狀態
		apiprotocol.PadServiceNo(apiprotocol.Code(cfg.ServiceNo), JwtgenerateFail): {Msg: "token 生成失敗"},
	}

	apiprotocol.Append(retStatusList)
}

// API 1000~1099 一般API狀態
const (
	// Success1000 : 正常回應
	Success1000 apiprotocol.Code = iota + 1000
	// Body解析錯誤
	API_BodyParseFail
	// Req 驗證錯誤
	API_ReqValidateFail
)

// Auth 1100~1199 會員權、登入相關狀態
const (
	Auth apiprotocol.Code = iota + 1100
	// 登入失敗
	Auth_LoginFail
	// 註冊失敗
	Auth_RegisterFail
	// 無此帳號
	Auth_NoAccount
	// 取得會員資料失敗
	Auth_GetMemberFail
)

// jwt 1200~1299 jwt相關狀態
const (
	Jwt apiprotocol.Code = iota + 1200
	// token 生成失敗
	JwtgenerateFail
)
