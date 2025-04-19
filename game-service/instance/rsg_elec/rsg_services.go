package rsg_elec

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"game_service/internal/apicaller"
	"game_service/internal/cachedata"
	encoder "game_service/pkg/encoder"
	"math/big"

	"strconv"
	"time"

	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

type reqData struct {
	encrypted    string
	signature    string
	timestampStr string
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var RSG_ELEC_SECRECT_INFO = RsgElecSecrect_Info{}
var ce = RSG_ELEC_SECRECT_INFO.CE
var AGENT_ID int

func Init() {
	cacheSecrect := cachedata.AgentSecretMap()
	cacheId := cachedata.AgentIdMap()
	rsgelecSecret, ok := cacheSecrect["rsg_elec"]
	if !ok {
		mlog.Error("rsg_elec secret not found")
		return
	}
	rsgelecId, ok := cacheId["rsg_elec"]
	if !ok {
		mlog.Error("rsg_elec id not found")
		return
	}
	// 取得代理商ID
	agentId, err := strconv.Atoi(rsgelecId)
	if err != nil {
		mlog.Error(err.Error())
		return
	}
	AGENT_ID = agentId
	// 取得代理商資訊
	err = json.Unmarshal([]byte(rsgelecSecret), &RSG_ELEC_SECRECT_INFO.CE)
	if err != nil {
		mlog.Error(err.Error())
		return
	}

	RSG_ELEC_SECRECT_INFO.RefreshTime = time.Now().Add(time.Minute * 30)
	ce = RSG_ELEC_SECRECT_INFO.CE
}

// 產交易序號ID
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			panic(err)
		}
		b[i] = letterBytes[randomIndex.Int64()]
	}
	return string(b)
}

// 請求數據處理
func reqDataProcessing(req interface{}) (reqData, error) {
	if !RSG_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !RSG_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	reqData := reqData{}
	body, err := json.Marshal(req)
	if err != nil {
		return reqData, err
	}

	// DES-CBC 加密
	encrypted, err := encoder.DESEncryptionDataCBC(body, ce.DesKey, ce.DesIV)
	if err != nil {
		return reqData, err
	}
	reqData.encrypted = encrypted
	timestamp := time.Now().Unix()
	reqData.timestampStr = strconv.FormatInt(timestamp, 10)

	// 生成 MD5 簽章
	data := ce.ClientID + ce.ClientSecret + reqData.timestampStr + encrypted
	reqData.signature = encoder.WGSportSignData(data)
	return reqData, nil
}

// 建立會員
func createUserAPI(username string) (CreatePlayerRes, error) {
	if !RSG_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !RSG_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	req := CreatePlayerReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     username,
		Currency:   "USDT",
	}
	url := ce.Httpdomain + string(CreateUserUri)
	resp := CreatePlayerRes{}

	reqData, err := reqDataProcessing(req)
	if err != nil {
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return errors.New(resp.Response.ErrorMessage)
		}

		// 解密回應數據
		data, err := encoder.DESDecryptCBC(r.Body(), ce.DesKey, ce.DesIV)
		if err != nil {
			return err
		}
		// 解析 JSON
		err = json.Unmarshal([]byte(data), &resp)
		if err != nil {
			return err
		} else {
			if resp.Response.ErrorCode != 0 {
				return errors.New(resp.Response.ErrorMessage)
			}
		}
		return nil
	}

	var header = map[string]string{
		"X-API-ClientID":  ce.ClientID,
		"X-API-Signature": reqData.signature,
		"X-API-Timestamp": reqData.timestampStr,
	}

	formData := "Msg=" + reqData.encrypted
	err = apicaller.RSGSendPostRequest(url, header, []byte(formData), handler)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// 加入指定遊戲
func loginAPI(username string) (EnterGameRes, error) {
	if !RSG_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !RSG_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	req := EnterGameReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     username,
		UserName:   username,
		GameId:     117,
		Currency:   "USDT",
		Language:   "zh-TW",
		ExitAction: "",
	}
	url := ce.Httpdomain + string(LoginUri)
	resp := EnterGameRes{}

	reqData, err := reqDataProcessing(req)
	if err != nil {
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return errors.New(resp.Response.ErrorMessage)
		}

		// 解密回應數據
		data, err := encoder.DESDecryptCBC(r.Body(), ce.DesKey, ce.DesIV)
		if err != nil {
			return err
		}

		// 解析 JSON
		err = json.Unmarshal([]byte(data), &resp)
		if err != nil {
			return err
		} else {
			if resp.Response.ErrorCode != 0 {
				return errors.New(resp.Response.ErrorMessage)
			}
		}
		return nil
	}

	var header = map[string]string{
		"X-API-ClientID":  ce.ClientID,
		"X-API-Signature": reqData.signature,
		"X-API-Timestamp": reqData.timestampStr,
	}

	formData := "Msg=" + reqData.encrypted
	err = apicaller.RSGSendPostRequest(url, header, []byte(formData), handler)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// TODO:WithBalance/Player/GetTransactionResult

// 踢人
func kickOutAPI(username string) (KickOutRes, error) {
	if !RSG_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !RSG_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	req := KickOutReq{
		KickType:   4,
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     username,
		GameId:     0,
	}
	url := ce.Httpdomain + string(KickOutUri)
	resp := KickOutRes{}

	reqData, err := reqDataProcessing(req)
	if err != nil {
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return errors.New(resp.Response.ErrorMessage)
		}

		// 解密回應數據
		data, err := encoder.DESDecryptCBC(r.Body(), ce.DesKey, ce.DesIV)
		if err != nil {
			return err
		}

		// 解析 JSON
		err = json.Unmarshal([]byte(data), &resp)
		if err != nil {
			return err
		} else {
			if resp.Response.ErrorCode != 0 {
				return errors.New(resp.Response.ErrorMessage)
			}
		}
		return nil
	}

	var header = map[string]string{
		"X-API-ClientID":  ce.ClientID,
		"X-API-Signature": reqData.signature,
		"X-API-Timestamp": reqData.timestampStr,
	}

	formData := "Msg=" + reqData.encrypted
	err = apicaller.RSGSendPostRequest(url, header, []byte(formData), handler)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// 取得遊戲詳細資訊
func getGameDetailAPI(frequency time.Duration) (GetGameDetailRes, error) {
	if !RSG_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !RSG_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	req := GetGameDetailReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		GameType:   1,
		TimeStart:  time.Now().Add(((-frequency) * 2) - (2 * time.Minute)).Format("2006-01-02 15:04"),
		TimeEnd:    time.Now().Add(((-frequency) * 1) - time.Minute).Format("2006-01-02 15:04"),
	}
	url := ce.Httpdomain + string(GetGameDetailUri)
	resp := GetGameDetailRes{}

	reqData, err := reqDataProcessing(req)
	if err != nil {
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return errors.New(resp.Response.ErrorMessage)
		}

		// 解密回應數據
		data, err := encoder.DESDecryptCBC(r.Body(), ce.DesKey, ce.DesIV)
		if err != nil {
			return err
		}

		// 解析 JSON
		err = json.Unmarshal([]byte(data), &resp)
		if err != nil {
			return err
		} else {
			if resp.Response.ErrorCode != 0 {
				return errors.New(resp.Response.ErrorMessage)
			}
		}
		return nil
	}

	var header = map[string]string{
		"X-API-ClientID":  ce.ClientID,
		"X-API-Signature": reqData.signature,
		"X-API-Timestamp": reqData.timestampStr,
	}

	formData := "Msg=" + reqData.encrypted
	err = apicaller.RSGSendPostRequest(url, header, []byte(formData), handler)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
