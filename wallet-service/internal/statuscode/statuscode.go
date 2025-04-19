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
		// Wallet 1100~1199 錢包相關狀態
		apiprotocol.PadServiceNo(apiprotocol.Code(cfg.ServiceNo), Wallet_TransferFail):       {Msg: "錢包轉點失敗"},
		apiprotocol.PadServiceNo(apiprotocol.Code(cfg.ServiceNo), Wallet_TransferRecordFail): {Msg: "錢包轉帳紀錄失敗"},
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

// Wallet 1100~1199 錢包相關狀態
const (
	Wallet apiprotocol.Code = iota + 1100
	// 錢包轉點失敗
	Wallet_TransferFail
	// 錢包轉帳紀錄失敗
	Wallet_TransferRecordFail
)
