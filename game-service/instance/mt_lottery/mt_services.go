package mt_lottery

import (
	"encoding/json"
	"errors"
	"fmt"
	"game_service/internal/apicaller"
	"game_service/internal/cachedata"
	"game_service/internal/database"
	"game_service/internal/service"
	"strconv"
	"time"

	"github.com/google/uuid"
	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

var MT_LOTTERY_SECRECT_INFO = MtLotterySecrect_Info{}
var hc = MT_LOTTERY_SECRECT_INFO.HC
var AGENT_ID int

func Init() {
	cacheSecrect := cachedata.AgentSecretMap()
	cacheId := cachedata.AgentIdMap()
	mtlotterySecret, ok := cacheSecrect["mt_lottery"]
	if !ok {
		mlog.Error("mt_lottery secret not found")
		return
	}
	mtlotteryId, ok := cacheId["mt_lottery"]
	if !ok {
		mlog.Error("mt_lottery id not found")
		return
	}
	// 取得代理商ID
	agentId, err := strconv.Atoi(mtlotteryId)
	if err != nil {
		mlog.Error(fmt.Sprintf("MT_Lottery取得代理商ID失敗: %s", err.Error()))
		return
	}
	AGENT_ID = agentId
	// 取得代理商資訊
	if err = json.Unmarshal([]byte(mtlotterySecret), &MT_LOTTERY_SECRECT_INFO.HC); err != nil {
		mlog.Error(fmt.Sprintf("MT_Lottery取得Secret_Info失敗: %s", err.Error()))
		return
	}

	MT_LOTTERY_SECRECT_INFO.RefreshTime = time.Now().Add(time.Minute * 30)
	hc = MT_LOTTERY_SECRECT_INFO.HC
}

// RegisterValid 註冊檢查 應檢查註冊資訊格式是否正確 並非檢查帳號是否存在
func AccountRegister(member *service.MemberGameAccountInfo) error {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}
	tx, err := database.MEMBER.TX()
	if err != nil {
		return err
	}

	member.GamePassword = uuid.New().String()[:20]
	// 註冊會員
	if err = service.CreateGameUser(tx, member.MemberId, *member); err != nil {
		return err
	}

	resp, err := createUserAPI(member.UserName, member.GamePassword)
	if err != nil {
		tx.Rollback()
		return err
	}
	if !resp.Rows.Enable {
		tx.Rollback()
		return errors.New("註冊失敗")
	}

	return tx.Commit()
}

// 建立會員
func createUserAPI(userName string, newPassword string) (CreateUserRes, error) {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}
	req := CreateUserReq{
		Account:  userName,
		Password: newPassword,
	}
	url := hc.Httpdomain + string(CreateUserUri)
	resp := Response{}
	result := CreateUserRes{}
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("jsonMarshal解析失敗 %s", err.Error()))
		return result, err
	}

	// cbc 加密
	encrypted, err := hc.EncryptionData([]byte(body))
	if err != nil {
		mlog.Error(fmt.Sprintf("資料加密失敗 %s", err.Error()))
		return result, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(fmt.Sprintf("請求回應不等於200 %s", err.Error()))
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err = json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("jsonUnmarshal錯誤、 %s", err.Error()))
			return err
		} else {
			if resp.Code != 200 {
				mlog.Error(fmt.Sprintf("遊戲商代碼錯誤 %s", err.Error()))
				return errors.New(resp.ErrMsg)
			}

			if resp.Data != "" {
				decryptedData, err := hc.DecryptData(resp.Data)
				if err != nil {
					mlog.Error(fmt.Sprintf("解密資料異常 %s", err.Error()))
					return err
				}
				err = json.Unmarshal([]byte(decryptedData), &result)
				if err != nil {
					mlog.Error(fmt.Sprintf("資料解析錯誤: %s", err.Error()))
					return err
				}
			}
		}
		return nil
	}
	err = apicaller.MTSendPostRequest(url, nil, string(encrypted), handler, hc.HashKey)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求失敗: %s", err.Error()))
		return result, err
	}
	mlog.Info(fmt.Sprintf("註冊帳號: %v\n", result))
	return result, nil
}

func CheckPoint(member service.MemberGameAccountInfo) (float64, error) {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}

	resp, err := checkPointAPI(member)
	if err != nil {
		return 0, err
	}

	return resp.Rows.MemberPoint, nil
}

// 查詢玩家可用點數
func checkPointAPI(member service.MemberGameAccountInfo) (CheckPointRes, error) {

	resp := Response{}
	req := CheckPointReq{
		Account:  hc.Prefix + member.UserName,
		Password: member.GamePassword,
		GameCode: "lottery",
	}
	url := hc.Httpdomain + string(CheckPointUri)
	result := CheckPointRes{}
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("JsonMarshal 失敗: %s", err.Error()))
		return result, err
	}

	// cbc 加密
	encrypted, err := hc.EncryptionData([]byte(body))
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗: %s", err.Error()))
		return result, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(fmt.Sprintf("httpStatusCode Not 200: %s", err.Error()))
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err = json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("jsonUnmarshal 失敗: %s", err.Error()))
			return err
		} else {
			if resp.Code != 200 {
				mlog.Error(fmt.Sprintf("請求回應碼錯誤(不等於200): %s", resp.ErrMsg))
				return errors.New(resp.ErrMsg)
			}

			if resp.Data != "" {
				decryptedData, err := hc.DecryptData(resp.Data)
				if err != nil {
					mlog.Error(fmt.Sprintf("解密失敗: %s", err.Error()))
					return err
				}
				err = json.Unmarshal([]byte(decryptedData), &result)
				if err != nil {
					mlog.Error(fmt.Sprintf("jsonUnMarshal error: %s", resp.ErrMsg))
					return err
				}
			}
		}
		return nil
	}
	err = apicaller.MTSendPostRequest(url, nil, string(encrypted), handler, hc.HashKey)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求失敗: %s", err.Error()))
		return result, err
	}
	return result, nil
}

// 取得URL Token
func getLoginUrlAPI(member service.MemberGameAccountInfo) (string, error) {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}
	resp := Response{}
	req := LoginReq{
		Account:  hc.Prefix + member.UserName,
		Password: member.GamePassword,
		GameCode: "lottery",
	}
	url := hc.Httpdomain + string(LoginUri)
	result := LoginRes{}
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("JsonMarshal 失敗: %s", err.Error()))
		return "資料解析異常", err
	}

	// cbc 加密
	encrypted, err := hc.EncryptionData([]byte(body))
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗: %s", err.Error()))
		return "資料加密異常", err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(fmt.Sprintf("httpStatusCode Not 200: %s", err.Error()))
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err = json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("jsonUnmarshal 失敗: %s", err.Error()))
			return err
		} else {
			if resp.Code != 200 {
				mlog.Error(fmt.Sprintf("請求回應碼錯誤(不等於200): %s", resp.ErrMsg))
				return errors.New(resp.ErrMsg)
			}

			if resp.Data != "" {
				decryptedData, err := hc.DecryptData(resp.Data)
				if err != nil {
					mlog.Error(fmt.Sprintf("解密失敗: %s", err.Error()))
					return err
				}
				err = json.Unmarshal([]byte(decryptedData), &result)
				if err != nil {
					mlog.Error(fmt.Sprintf("jsonUnMarshal error: %s", resp.ErrMsg))
					return err
				}
				if result.TotalRows < 1 {
					mlog.Error(fmt.Sprintf("加入遊戲失敗 error: %s", errors.New("進入遊戲失敗")))
					return errors.New("加入遊戲失敗")
				}
			}
		}
		return nil
	}
	if err = apicaller.MTSendPostRequest(url, nil, string(encrypted), handler, hc.HashKey); err != nil {
		mlog.Error(fmt.Sprintf("請求失敗: %s", err.Error()))
		return "伺服器服務異常", err
	}
	return result.Rows[0].Url, nil
}

func TransferPoint(member service.MemberGameAccountInfo, point float64) (TransPointRes, error) {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}

	// 確認遊戲商帳號是否存在
	resp, err := transPointAPI(member, point)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// 轉入(出)玩家點數
func transPointAPI(member service.MemberGameAccountInfo, point float64) (TransPointRes, error) {
	tradeNumber := uuid.New().String()
	resp := Response{}
	req := TransPointReq{
		Account:    hc.Prefix + member.UserName,
		Password:   member.GamePassword,
		Point:      point,
		TradeOrder: tradeNumber,
	}
	url := hc.Httpdomain + string(TransPointUri)
	result := TransPointRes{}
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("JsonMarshal 失敗: %s", err.Error()))
		return result, err
	}

	// cbc 加密
	encrypted, err := hc.EncryptionData([]byte(body))
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗: %s", err.Error()))
		return result, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(fmt.Sprintf("httpStatusCode Not 200: %s", err.Error()))
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err = json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("jsonUnmarshal 失敗: %s", err.Error()))
			return err
		} else {
			if resp.Code != 200 {
				mlog.Error(fmt.Sprintf("請求回應碼錯誤(不等於200): %s", resp.ErrMsg))
				return errors.New(resp.ErrMsg)
			}

			if resp.Data != "" {
				decryptedData, err := hc.DecryptData(resp.Data)
				if err != nil {
					mlog.Error(fmt.Sprintf("解密資料異常: %s", err.Error()))
					return err
				}
				err = json.Unmarshal([]byte(decryptedData), &result)
				if err != nil {
					mlog.Error(fmt.Sprintf("解析資料異常: %s", err.Error()))
					return err
				}
			}
		}
		return nil
	}
	if err = apicaller.MTSendPostRequest(url, nil, string(encrypted), handler, hc.HashKey); err != nil {
		mlog.Error(fmt.Sprintf("請求失敗: %s", err.Error()))
		return result, err
	}
	return result, nil
}

// 查詢玩家交易紀錄
func getTransactionLogAPI(req TransactionLogReq, url string) (TransactionLogRes, error) {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}
	resp := Response{}
	result := TransactionLogRes{}
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
		return result, err
	}

	// cbc 加密
	encrypted, err := hc.EncryptionData([]byte(body))
	if err != nil {
		return result, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err = json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			if resp.Code != 200 {
				return errors.New(resp.ErrMsg)
			}

			if resp.Data != "" {
				decryptedData, err := hc.DecryptData(resp.Data)
				if err != nil {
					mlog.Error("解密資料異常")
					return err
				}
				err = json.Unmarshal([]byte(decryptedData), &result)
				if err != nil {
					mlog.Error("解析資料異常")
					return err
				}
			}
		}
		return nil
	}
	if err = apicaller.MTSendPostRequest(url, nil, string(encrypted), handler, hc.HashKey); err != nil {
		mlog.Error(err.Error())
		return result, err
	}
	return result, nil
}

// 踢出遊戲
func KickOut(member service.MemberGameAccountInfo) error {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}
	if _, err := logoutAPI(member.UserName, member.GamePassword); err != nil {
		mlog.Error(fmt.Sprintf("踢出失敗: %s", err.Error()))
		return err
	}
	return nil
}

// 玩家登出遊戲
func logoutAPI(username string, password string) (LogoutRes, error) {
	req := LogoutReq{
		Account:  hc.Prefix + username,
		Password: password,
		GameCode: "lottery",
	}
	url := hc.Httpdomain + string(LogoutUri)
	resp := Response{}
	result := LogoutRes{}
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("jsonMarshal失敗: %s", err.Error()))
		return result, err
	}

	// cbc 加密
	encrypted, err := hc.EncryptionData([]byte(body))
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗: %s", err.Error()))
		return result, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(fmt.Sprintf("httpStatusCode not 200: %s", err.Error()))
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err = json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("jsonMarshal 錯誤: %s", err.Error()))
			return err
		} else {
			if resp.Code != 200 {
				mlog.Error(fmt.Sprintf("遊戲商回應碼錯誤: %s", resp.ErrMsg))
				return errors.New(resp.ErrMsg)
			}

			if resp.Data != "" {
				decryptedData, err := hc.DecryptData(resp.Data)
				if err != nil {
					mlog.Error(fmt.Sprintf("解密資料異常: %s", err.Error()))
					return err
				}
				err = json.Unmarshal([]byte(decryptedData), &result)
				if err != nil {
					mlog.Error(fmt.Sprintf("jsonUnMarshal 錯誤: %s", err.Error()))
					return err
				}
			}
		}
		return nil
	}
	if err = apicaller.MTSendPostRequest(url, nil, string(encrypted), handler, hc.HashKey); err != nil {
		mlog.Error(fmt.Sprintf("請求失敗: %s", err.Error()))
		return result, err
	}
	return result, nil
}

// 可用限紅查詢
func checkHandicapAPI(req HandicapCheckReq, url string) (HandicapCheckRes, error) {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}
	resp := Response{}
	result := HandicapCheckRes{}
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
		return result, err
	}

	// cbc 加密
	encrypted, err := hc.EncryptionData([]byte(body))
	if err != nil {
		return result, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err = json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			if resp.Code != 200 {
				return errors.New(resp.ErrMsg)
			}

			if resp.Data != "" {
				decryptedData, err := hc.DecryptData(resp.Data)
				if err != nil {
					mlog.Error("解密資料異常")
					return err
				}
				err = json.Unmarshal([]byte(decryptedData), &result)
				if err != nil {
					mlog.Error("解析資料異常")
					return err
				}
			}
		}
		return nil
	}
	if err = apicaller.MTSendPostRequest(url, nil, string(encrypted), handler, hc.HashKey); err != nil {
		mlog.Error(err.Error())
		return result, err
	}
	return result, nil
}

// 會員限紅設定
func handicapSettingAPI(req MemberHandicapSettingReq, url string) (MemberHandicapSettingRes, error) {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}
	resp := Response{}
	result := MemberHandicapSettingRes{}
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
		return result, err
	}

	// cbc 加密
	encrypted, err := hc.EncryptionData([]byte(body))
	if err != nil {
		return result, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err = json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			if resp.Code != 200 {
				return errors.New(resp.ErrMsg)
			}

			if resp.Data != "" {
				decryptedData, err := hc.DecryptData(resp.Data)
				if err != nil {
					mlog.Error("解密資料異常")
					return err
				}
				err = json.Unmarshal([]byte(decryptedData), &result)
				if err != nil {
					mlog.Error("解析資料異常")
					return err
				}
			}
		}
		return nil
	}

	if err = apicaller.MTSendPostRequest(url, nil, string(encrypted), handler, hc.HashKey); err != nil {
		return result, err
	}
	return result, nil
}

// 查詢指定日期投注記錄
func getBetOrderAPI(req BetOrderReq, url string) (string, error) {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}
	resp, result := Response{}, BetOrderRes{}

	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
		return "", err
	}

	// cbc 加密
	encrypted, err := hc.EncryptionData([]byte(body))
	if err != nil {
		return "", err
	}

	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		if err = json.Unmarshal(r.Body(), &resp); err != nil {
			return err
		} else {
			if resp.Code != 200 {
				if resp.ErrMsg == "API.betOrder_not_exists" {
					mlog.Info("API.betOrder_not_exists")
					return nil
				} else {
					return errors.New(resp.ErrMsg)
				}
			}
			if resp.Data != "" {
				decryptedData, err := hc.DecryptData(resp.Data)
				if err != nil {
					mlog.Error("解密資料異常")
					return err
				}
				if err = json.Unmarshal([]byte(decryptedData), &result); err != nil {
					mlog.Error("解析資料異常")
					return err
				}
				data = string(decryptedData)
			}
		}
		return nil
	}
	if err = apicaller.MTSendPostRequest(url, nil, string(encrypted), handler, hc.HashKey); err != nil {
		return "", err
	}
	return data, nil
}

// 查詢指定日期投注記錄v2
func getBetOrderV2API(req BetOrderV2Req, url string) (BetOrderV2Res, error) {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}
	resp := Response{}
	result := BetOrderV2Res{}
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
		return result, err
	}

	// cbc 加密
	encrypted, err := hc.EncryptionData([]byte(body))
	if err != nil {
		return result, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		if err = json.Unmarshal(r.Body(), &resp); err != nil {
			if err.Error() == "json: cannot unmarshal array into Go struct field Response.data of type string" {
				if resp.Code != 200 {
					if resp.ErrMsg == "API.betOrder_not_exists" {
						return nil
					} else {
						return errors.New(resp.ErrMsg)
					}
				}
			}
			return err
		} else {
			if resp.Data != "" {
				decryptedData, err := hc.DecryptData(resp.Data)
				if err != nil {
					return err
				}
				if err = json.Unmarshal([]byte(decryptedData), &result); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if err = apicaller.MTSendPostRequest(url, nil, string(encrypted), handler, hc.HashKey); err != nil {
		mlog.Error(err.Error())
		return result, err
	}
	return result, nil
}
