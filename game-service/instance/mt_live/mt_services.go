package mt_live

import (
	"errors"

	"game_service/internal/apicaller"
	"game_service/internal/cachedata"
	"game_service/internal/database"
	"game_service/internal/service"
	"strconv"
	"time"

	"encoding/json"

	"fmt"

	"github.com/google/uuid"

	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

var MT_LIVE_SECRECT_INFO = MtLiveSecrect_Info{}
var ce = MT_LIVE_SECRECT_INFO.CE
var AGENT_ID int

func Init() {
	cacheSecrect := cachedata.AgentSecretMap()
	cacheId := cachedata.AgentIdMap()
	mtliveSecret, ok := cacheSecrect["mt_live"]
	if !ok {
		mlog.Error("mt_live secret not found")
		return
	}
	mtliveId, ok := cacheId["mt_live"]
	if !ok {
		mlog.Error("mt_live id not found")
		return
	}
	// 取得代理商ID
	agentId, err := strconv.Atoi(mtliveId)
	if err != nil {
		mlog.Error(fmt.Sprintf("MT_Live取得代理商ID失敗: %s", err.Error()))
		return
	}
	AGENT_ID = agentId
	// 取得代理商資訊
	err = json.Unmarshal([]byte(mtliveSecret), &MT_LIVE_SECRECT_INFO.CE)
	if err != nil {
		mlog.Error(fmt.Sprintf("MT_Live取得Secret_Info失敗: %s", err.Error()))
		return
	}

	MT_LIVE_SECRECT_INFO.RefreshTime = time.Now().Add(time.Minute * 30)
	ce = MT_LIVE_SECRECT_INFO.CE
}

// RegisterValid 註冊檢查 應檢查註冊資訊格式是否正確 並非檢查帳號是否存在
func AccountRegister(mid int, member *service.MemberGameAccountInfo) error {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.Before(time.Now()) {
		Init()
	}
	tx, err := database.MEMBER.TX()
	if err != nil {
		mlog.Error(fmt.Sprintf("資料庫連線錯誤: %s", err.Error()))
		return err
	}
	// 註冊會員
	member.GamePassword = uuid.New().String()
	err = service.CreateUser(tx, mid, *member)
	if err != nil {
		mlog.Info(fmt.Sprintf("會員註冊寫入失敗: %s", err.Error()))
		return err
	}

	resp, err := createUserAPI(*member)
	if err != nil {
		mlog.Info(fmt.Sprintf("遊戲會員註冊失敗: %s", err.Error()))
		tx.Rollback()
		return err
	}
	if resp.Data.Result != 1 {
		mlog.Info(fmt.Sprintf("遊戲會員註冊失敗: %s", resp.Message))
		tx.Rollback()
		return errors.New("註冊失敗")
	}

	tx.Commit()
	return nil
}

// 建立會員
func createUserAPI(member service.MemberGameAccountInfo) (CreateUserRes, error) {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	req := CreateUserReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     member.UserName,
		UserName:   member.NickName,
		Currency:   "USD",
	}
	url := ce.Httpdomain + string(CreateUserUri)
	resp := CreateUserRes{}

	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("createUserAPI jsonMarshal失敗: %s", err.Error()))
		return resp, err
	}
	// cbc 加密
	encrypted, err := ce.EncryptionData(body)
	if err != nil {
		mlog.Error(fmt.Sprintf("createUserAPI cbc加密錯誤: %s", err.Error()))
		return resp, err
	}

	timestamp := time.Now().Unix()
	// md5 加簽
	signature := ce.SignatureData(encrypted, timestamp)

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(fmt.Sprintf("createUserAPI httpStatusCode Not 200: %s", err.Error()))
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("createUser jsonUnmarshal錯誤 %s", err.Error()))
			return err
		} else {
			if resp.Code != "00000" {
				mlog.Error(fmt.Sprintf("createUser 請求回應碼錯誤(不等於00000): %s", resp.Message))
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encrypted, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求失敗 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 修改會員資訊
func editUserAPI(user EditUserReq, url string) (EditUserRes, error) {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	resp := EditUserRes{}

	body, err := json.Marshal(user)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	// cbc 加密
	encrypted, err := ce.EncryptionData(body)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return resp, err
	}

	timestamp := time.Now().Unix()
	// md5 加簽
	signature := ce.SignatureData(encrypted, timestamp)

	handler := func(r *fasthttp.Response) error {
		resp := CreateUserRes{}
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			if resp.Code != "00000" {
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encrypted, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

func GameAccountExist(username string) (*CheckUserRes, error) {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	// 確認遊戲商帳號是否存在
	resp, err := checkUserAPI(username)
	if err != nil {
		mlog.Error(fmt.Sprintf("checkUserAPI Error: %s", err.Error()))
		return nil, err
	}
	if resp.Data.Result != 1 {
		mlog.Info("遊戲商帳號不存在")
		return &resp, err
	}

	return &resp, nil
}

// 查詢會員
func checkUserAPI(username string) (CheckUserRes, error) {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	req := CheckUserReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     username,
	}
	url := ce.Httpdomain + string(CheckUserUri)
	resp := CheckUserRes{}

	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("JsonMarshal 失敗: %s", err.Error()))
		return resp, err
	}
	// cbc 加密
	encrypted, err := ce.EncryptionData(body)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗: %s", err.Error()))
		return resp, err
	}

	timestamp := time.Now().Unix()
	// md5 加簽
	signature := ce.SignatureData(encrypted, timestamp)

	handler := func(r *fasthttp.Response) error {
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("jsonUnmarshal 失敗: %s", err.Error()))
			return err
		} else {
			if resp.Code != "00000" {
				mlog.Error(fmt.Sprintf("請求回應碼錯誤(不等於00000): %s", resp.Message))
				return errors.New(resp.Message)
			}
		}
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(fmt.Sprintf("httpStatusCode Not 200: %s", resp.Message))
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}

		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encrypted, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求失敗: %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 取得URL Token
func getURLTokenAPI(username string) (string, error) {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	req := GetURLTokenReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     username,
		Language:   "zh-TW",
	}
	url := ce.Httpdomain + string(GetURLTokenUri)
	resp := GetURLTokenRes{}

	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("JsonMarshal 失敗: %s", err.Error()))
		return "資料解析異常", err
	}

	encryptedData, err := ce.EncryptionData(body)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗: %s", err.Error()))
		return "資料加密異常", err
	}

	timestamp := time.Now().Unix()
	signature := ce.SignatureData(encryptedData, timestamp)

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(fmt.Sprintf("httpStatusCode Not 200: %s", resp.Message))
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("jsonUnmarshal 失敗: %s", err.Error()))
			return err
		} else {
			if resp.Code != "00000" {
				mlog.Error(fmt.Sprintf("請求回應碼錯誤(不等於00000): %s", resp.Message))
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encryptedData, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求失敗: %s", err.Error()))
		return "伺服器服務異常", err
	}
	return resp.Data.Url, nil
}

func TransferIn(username string, point float64) (DepositRes, error) {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	// 確認遊戲商帳號是否存在
	resp, err := depositAPI(username, point)
	if err != nil {
		return resp, err
	}
	if resp.Code != "00000" {
		return resp, errors.New("存款失敗")
	}

	return resp, nil
}

// 充入點數
func depositAPI(username string, point float64) (DepositRes, error) {
	resp := DepositRes{}
	req := DepositReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     username,
		Balance:    point,
	}
	url := ce.Httpdomain + string(DepositUri)
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("JsonMarshal 失敗: %s", err.Error()))
		return resp, err
	}

	encryptedData, err := ce.EncryptionData(body)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗: %s", err.Error()))
		return resp, err
	}

	timestamp := time.Now().Unix()
	signature := ce.SignatureData(encryptedData, timestamp)

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(fmt.Sprintf("httpStatusCode Not 200: %s", resp.Message))
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("jsonUnmarshal 失敗: %s", err.Error()))
			return err
		} else {
			if resp.Code != "00000" {
				mlog.Error(fmt.Sprintf("請求回應碼錯誤(不等於00000): %s", resp.Message))
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encryptedData, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求失敗: %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 遊戲商提出點數
func TransferOut(username string, point float64) (WithdrawRes, error) {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	// 確認遊戲商帳號是否存在
	resp, err := withdrawAPI(username, point)

	if err != nil {
		return resp, err
	}
	if resp.Code != "00000" {
		return resp, errors.New("提款失敗")
	}

	return resp, nil
}

// 取出點數
func withdrawAPI(username string, point float64) (WithdrawRes, error) {
	resp := WithdrawRes{}
	req := WithdrawReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     username,
		Balance:    point,
	}
	url := ce.Httpdomain + string(WithdrawUri)
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}

	encryptedData, err := ce.EncryptionData(body)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return resp, err
	}

	timestamp := time.Now().Unix()
	signature := ce.SignatureData(encryptedData, timestamp)

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			if resp.Code != "00000" {
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encryptedData, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

func CheckBalance(username string) (GetBalanceRes, error) {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	// 確認遊戲商帳號是否存在
	resp, err := getBalanceAPI(username)
	if err != nil {
		return resp, err
	}
	if resp.Code != "00000" {
		return resp, errors.New("查詢點數失敗")
	}

	return resp, nil
}

// 查詢點數
func getBalanceAPI(username string) (GetBalanceRes, error) {
	resp := GetBalanceRes{}
	req := GetBalanceReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     username,
	}
	url := ce.Httpdomain + string(GetBalanceUri)
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("JsonMarshal 失敗: %s", err.Error()))
		return resp, err
	}

	encryptedData, err := ce.EncryptionData(body)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗: %s", err.Error()))
		return resp, err
	}

	timestamp := time.Now().Unix()
	signature := ce.SignatureData(encryptedData, timestamp)

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			if resp.Code != "00000" {
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encryptedData, handler)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func Kickout(member service.MemberGameAccountInfo) (int, error) {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	// 確認遊戲商帳號是否存在
	_, err := kickoutAPI(member.UserName)
	if err != nil {
		mlog.Error(fmt.Sprintf("踢人錯誤: %s", err.Error()))
		return 0, err
	}
	return 1, nil
}

// 踢除在線玩家
func kickoutAPI(username string) (KickoutRes, error) {
	resp := KickoutRes{}
	req := KickoutReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     username,
	}
	url := ce.Httpdomain + string(KickoutUri)

	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("jsonMarshal失敗: %s", err.Error()))
		return resp, err
	}

	encryptedData, err := ce.EncryptionData(body)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗: %s", err.Error()))
		return resp, err
	}

	timestamp := time.Now().Unix()
	signature := ce.SignatureData(encryptedData, timestamp)

	handler := func(r *fasthttp.Response) error {
		resp := KickoutRes{}
		if r.StatusCode() != fasthttp.StatusOK {
			mlog.Error(fmt.Sprintf("httpStatusCode not 200: %s", err.Error()))
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("jsonMarshal 錯誤: %s", err.Error()))
			return err
		} else {
			if resp.Code != "00000" {
				mlog.Error(fmt.Sprintf("遊戲商回應碼錯誤: %s", resp.Message))
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encryptedData, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求失敗: %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 取得線上玩家列表
func getPlayerOnlineListAPI(req PlayerOnlineListReq, url string) (PlayerOnlineListRes, error) {
	resp := PlayerOnlineListRes{}

	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}

	encryptedData, err := ce.EncryptionData(body)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return resp, err
	}

	timestamp := time.Now().Unix()
	signature := ce.SignatureData(encryptedData, timestamp)

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			if resp.Code != "00000" {
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encryptedData, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 取得遊戲商注單記錄
func getBetRecordAPI(req GetBetRecordReq, url string) (GetBetRecordRes, error) {
	resp := GetBetRecordRes{}
	body, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}
	encryptedData, err := ce.EncryptionData(body)
	if err != nil {
		return resp, err
	}
	timestamp := time.Now().Unix()
	signature := ce.SignatureData(encryptedData, timestamp)

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		if err := json.Unmarshal(r.Body(), &resp); err != nil {
			return err
		} else {
			if resp.Code != "00000" {
				return errors.New(resp.Message)
			}
		}
		return nil
	}
	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}
	if err = apicaller.SendPostRequest(url, header, encryptedData, handler); err != nil {
		return resp, err
	}
	return resp, nil
}

// 與遊戲商獲取交易(儲值)紀錄
func FindTransferRecord(sn string) (FindTransationRecordRes, error) {
	// 確認遊戲商帳號是否存在
	resp, err := findTransationRecordAPI(
		FindTransationRecordReq{
			SystemCode: ce.SystemCode,
			WebId:      ce.WebID,
			TransferSN: sn,
		},
		ce.Httpdomain+string(FindTransationRecordUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	if resp.Code != "00000" {
		mlog.Error("交易紀錄不存在")
		return resp, errors.New("交易紀錄不存在")
	}
	if resp.Data.List[0].Type != "1" {
		mlog.Error("交易紀錄不是儲值")
		return resp, errors.New("交易紀錄不是儲值")
	}

	return resp, nil
}

// 查詢遊戲商單筆交易記錄
func findTransationRecordAPI(req FindTransationRecordReq, url string) (FindTransationRecordRes, error) {
	resp := FindTransationRecordRes{}

	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}

	encryptedData, err := ce.EncryptionData(body)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return resp, err
	}

	timestamp := time.Now().Unix()
	signature := ce.SignatureData(encryptedData, timestamp)

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			if resp.Code != "00000" {
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encryptedData, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 取得遊戲商打賞記錄
func getDonateRecordAPI(req DonateRecordReq, url string) (DonateRecordRes, error) {
	resp := DonateRecordRes{}
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}

	encryptedData, err := ce.EncryptionData(body)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return resp, err
	}

	timestamp := time.Now().Unix()
	signature := ce.SignatureData(encryptedData, timestamp)

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			if resp.Code != "00000" {
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encryptedData, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 取得遊戲商異動注單記錄
func getUpdateBetRecordAPI(req GetUpdateBetRecordReq, url string) (GetUpdateBetRecordRes, error) {
	resp := GetUpdateBetRecordRes{}

	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}

	encryptedData, err := ce.EncryptionData(body)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return resp, err
	}

	timestamp := time.Now().Unix()
	signature := ce.SignatureData(encryptedData, timestamp)

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			if resp.Code != "00000" {
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encryptedData, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

// 交易紀錄
func getTransactionRecordAPI(req GetTransationRecordReq, url string) (GetTransationRecordRes, error) {
	resp := GetTransationRecordRes{}

	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}

	encryptedData, err := ce.EncryptionData(body)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return resp, err
	}

	timestamp := time.Now().Unix()
	signature := ce.SignatureData(encryptedData, timestamp)

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			return err
		} else {
			if resp.Code != "00000" {
				return errors.New(resp.Message)
			}
		}
		return nil
	}

	var header = map[string]string{
		"APICI": ce.ClientID,
		"APISI": signature,
		"APITS": strconv.Itoa(int(timestamp)),
	}

	err = apicaller.SendPostRequest(url, header, encryptedData, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}
