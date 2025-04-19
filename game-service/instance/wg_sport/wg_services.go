package wg_sport

import (
	"encoding/json"
	"errors"
	"fmt"
	"game_service/internal/apicaller"
	"game_service/internal/cachedata"
	"game_service/internal/database"
	"game_service/internal/service"
	"game_service/pkg/encoder"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

var WG_SPORT_SECRECT_INFO = WgSportSecrect_Info{}
var WSSI = WG_SPORT_SECRECT_INFO.WSSI
var AGENT_ID int

func Init() {
	cacheSecrect := cachedata.AgentSecretMap()
	cacheId := cachedata.AgentIdMap()
	wgsportSecret, ok := cacheSecrect["wg_sport"]
	if !ok {
		mlog.Error("wg_sport secret not found")
		return
	}
	wgsportId, ok := cacheId["wg_sport"]
	if !ok {
		mlog.Error("wg_sport id not found")
		return
	}
	// 取得代理商ID
	agentId, err := strconv.Atoi(wgsportId)
	if err != nil {
		mlog.Error(err.Error())
		return
	}
	AGENT_ID = agentId
	// 取得代理商資訊
	err = json.Unmarshal([]byte(wgsportSecret), &WG_SPORT_SECRECT_INFO.WSSI)
	if err != nil {
		return
	}

	WG_SPORT_SECRECT_INFO.RefreshTime = time.Now().Add(time.Minute * 30)
	WSSI = WG_SPORT_SECRECT_INFO.WSSI
}

func getMD5string(req interface{}) string {

	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return ""
	}
	var userMap map[string]string
	if err := json.Unmarshal(body, &userMap); err != nil {
		mlog.Error(fmt.Sprintf("json.Unmarshal error %s", err.Error()))
	}
	val := reflect.ValueOf(req)
	typeOfReq := val.Type()

	var orderedKeys []string
	for i := 0; i < val.NumField(); i++ {
		field := typeOfReq.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			orderedKeys = append(orderedKeys, jsonTag)
		}
	}

	// 按照結構字段的順序排序
	var orderedParams []KV
	for _, key := range orderedKeys {
		if value, exists := userMap[key]; exists && value != "" {
			orderedParams = append(orderedParams, KV{key, value})
		}
	}

	// 打印過濾後的鍵
	keys := make([]string, len(orderedParams))
	for i, kv := range orderedParams {
		keys[i] = kv.Key
	}

	// 生成鍵值對的陣列
	keyValuePairs := make([]string, len(orderedParams))
	for i, kv := range orderedParams {
		keyValuePairs[i] = fmt.Sprintf("%s=%s", kv.Key, kv.Value)
	}
	queryString := strings.Join(keyValuePairs, "&")

	return queryString
}

// 創建遊戲商帳號
func createUserAPI(member *service.MemberGameAccountInfo) (CreateUserRes, error) {
	if !WG_SPORT_SECRECT_INFO.RefreshTime.After(time.Now()) || !WG_SPORT_SECRECT_INFO.Ischeck {
		Init()
	}
	resp := CreateUserRes{}
	req := CreateUserReq{
		Prefix:     WSSI.Prefix,
		Username:   member.UserName,
		Nickname:   member.NickName,
		Upusername: WSSI.Upusername,
	}
	url := WSSI.Httpdomain + string(CreateUserUri)
	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("json.UnMarshal error %s", err.Error()))
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 檢查遊戲商帳號是否存在
func checkUserAPI(username string) (CheckUserRes, error) {
	resp := CheckUserRes{}
	req := CheckUserReq{
		Prefix:   WSSI.Prefix,
		Username: username,
	}

	url := WSSI.Httpdomain + string(CheckUserUri)
	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("json.UnMarshal error %s", err.Error()))
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 修改遊戲商帳號資訊
func editUserAPI(req EditUserReq, url string) (EditUserRes, error) {
	resp := EditUserRes{}

	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("json.UnMarshal error %s", err.Error()))
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 進入遊戲
func forwardGameAPI(username string) (string, error) {
	if !WG_SPORT_SECRECT_INFO.RefreshTime.After(time.Now()) || !WG_SPORT_SECRECT_INFO.Ischeck {
		Init()
	}
	req := ForwardGameReq{
		Prefix:   WSSI.Prefix,
		Username: username,
	}
	resp := ForwardGameRes{}
	url := WSSI.Httpdomain + string(ForwardGameUri)

	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return "", err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("json.UnMarshal error %s", err.Error()))
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return "", err
	}
	gameUrl := "https:" + resp.Url
	return gameUrl, nil
}

func CheckPoint(member service.MemberGameAccountInfo) (float64, error) {
	if !WG_SPORT_SECRECT_INFO.RefreshTime.After(time.Now()) || !WG_SPORT_SECRECT_INFO.Ischeck {
		Init()
	}
	resp, err := PointUserAPI(member.UserName)
	if err != nil {
		mlog.Info(fmt.Sprintf("確認點數失敗: %s", err.Error()))
		return 0, err
	}

	return resp, nil
}

// 取得點數
func PointUserAPI(username string) (float64, error) {
	if !WG_SPORT_SECRECT_INFO.RefreshTime.After(time.Now()) || !WG_SPORT_SECRECT_INFO.Ischeck {
		Init()
	}
	resp := PointUserRes{}
	req := PointUserReq{
		Prefix:   WSSI.Prefix,
		Username: username,
	}
	url := WSSI.Httpdomain + string(GetPointUserUri)
	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return 0, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("json.UnMarshal error %s", err.Error()))
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return 0, err
	}
	return resp.Points, nil
}

// 遊戲商提出點數
func TransferIn(member service.MemberGameAccountInfo, bpoint float64, point float64) (TransferInRes, error) {
	if !WG_SPORT_SECRECT_INFO.RefreshTime.After(time.Now()) || !WG_SPORT_SECRECT_INFO.Ischeck {
		Init()
	}
	// 確認遊戲商帳號是否存在
	resp, err := transferInAPI(member.UserName, bpoint, point)

	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	if resp.Error_code != "OK" {
		mlog.Error("提款失敗")
		return resp, errors.New("提款失敗")
	}

	return resp, nil
}

// 轉入點數
func transferInAPI(username string, bamount float64, point float64) (TransferInRes, error) {
	resp := TransferInRes{}
	amountInt := int(point)
	apoints := float64(amountInt) + bamount
	req := TransferInReq{
		Prefix:   WSSI.Prefix,
		Username: username,
		Bpoints:  strconv.FormatFloat(bamount, 'f', 3, 64),
		Points:   strconv.Itoa(amountInt),
		Apoints:  strconv.FormatFloat(apoints, 'f', 3, 64),
	}

	url := WSSI.Httpdomain + string(TransferInUri)
	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("json.UnMarshal error %s", err.Error()))
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 遊戲商提出點數
func TransferOut(member service.MemberGameAccountInfo, bpoint float64) (float64, error) {
	if !WG_SPORT_SECRECT_INFO.RefreshTime.After(time.Now()) || !WG_SPORT_SECRECT_INFO.Ischeck {
		Init()
	}
	// 確認遊戲商帳號是否存在
	resp, err := transferOutAPI(member.UserName, bpoint)

	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}

	return resp, nil
}

// 轉出點數
func transferOutAPI(username string, bpoint float64) (float64, error) {
	resp := TransferOutRes{}
	bpoints := strconv.FormatFloat(bpoint, 'f', 3, 64)
	intPoint := int(bpoint)
	points := strconv.Itoa(intPoint)
	apoint := bpoint - float64(intPoint)
	apoints := strconv.FormatFloat(apoint, 'f', 3, 64)
	req := TransferOutReq{
		Prefix:   WSSI.Prefix,
		Username: username,
		Bpoints:  bpoints,
		Points:   points,
		Apoints:  apoints,
	}
	url := WSSI.Httpdomain + string(TransferOutUri)
	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return 0, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("json.UnMarshal error %s", err.Error()))
			return err
		}
		if resp.Error_code != "OK" {
			mlog.Error("提款失敗")
			return errors.New("提款失敗")
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return 0, err
	}
	return float64(intPoint), nil
}

// 確認轉點交易單號
func transferCheckAPI(req TransferCheckReq, url string) (TransferCheckRes, error) {
	resp := TransferCheckRes{}

	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("json.UnMarshal error %s", err.Error()))
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 查詢注單
func buyListGetAPI(req BuyListGetReq) (BuyListGetRes, error) {
	resp := BuyListGetRes{}
	url := WSSI.Httpdomain + string(BuyListGetUri)
	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		fmt.Println("resp:", string(r.Body()))
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("json.UnMarshal error %s", err.Error()))
			return err
		}

		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 查詢注單內容
func buyDetailGetAPI(req BuyDetailReq, url string) (string, error) {
	var resp string

	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		resp = string(r.Body())
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

func KickOut(member service.MemberGameAccountInfo) error {
	if !WG_SPORT_SECRECT_INFO.RefreshTime.After(time.Now()) || !WG_SPORT_SECRECT_INFO.Ischeck {
		Init()
	}
	_, err := kickOutAPI(member.UserName)
	if err != nil {
		mlog.Error(fmt.Sprintf("踢出失敗: %s", err.Error()))
		return err
	}
	return nil
}

// 踢出會員
func kickOutAPI(username string) (KickUserRes, error) {
	resp := KickUserRes{}

	req := KickUserReq{
		Prefix:   WSSI.Prefix,
		Username: username,
	}
	url := WSSI.Httpdomain + string(KickUserUri)

	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("json.UnMarshal error %s", err.Error()))
			return err
		}
		return nil
	}

	var header = map[string]string{}

	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// RegisterValid 註冊檢查 應檢查註冊資訊格式是否正確 並非檢查帳號是否存在
func AccountRegister(member *service.MemberGameAccountInfo) error {
	if !WG_SPORT_SECRECT_INFO.RefreshTime.After(time.Now()) || !WG_SPORT_SECRECT_INFO.Ischeck {
		Init()
	}
	// 註冊會員
	tx, err := database.MEMBER.TX()
	if err != nil {
		mlog.Error(fmt.Sprintf("資料庫連線錯誤: %s", err.Error()))
		return err
	}
	// 註冊會員
	member.GamePassword = uuid.New().String()[:20]
	err = service.CreateGameUser(tx, member.MemberId, *member)
	if err != nil {
		mlog.Error(fmt.Sprintf("CreateUser失敗: %s", err.Error()))
		return err
	}

	resp, err := createUserAPI(member)
	if err != nil {
		mlog.Error(fmt.Sprintf("遊戲商CreateUser失敗: %s", err.Error()))
		tx.Rollback()
		return err
	}
	if resp.Error_code != "OK" {
		mlog.Error("註冊失敗")
		tx.Rollback()
		return errors.New("註冊失敗")
	}

	tx.Commit()
	return nil
}

// 進入遊戲
func AccountLogin(login ForwardGameReq) (string, error) {
	if !WG_SPORT_SECRECT_INFO.RefreshTime.After(time.Now()) || !WG_SPORT_SECRECT_INFO.Ischeck {
		Init()
	}
	resp, err := getURLTokenAPI(
		login,
		WSSI.Httpdomain+string(ForwardGameUri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return "", err
	}
	if resp.Error_code == "OK" {
		return resp.Url, nil
	} else {
		return "", errors.New("進入遊戲失敗")
	}
}

// 取得URL Token
func getURLTokenAPI(req ForwardGameReq, url string) (ForwardGameRes, error) {
	resp := ForwardGameRes{}

	queryString := getMD5string(req)
	queryString += "&md5key=" + WSSI.Md5key

	// MD5加密
	signature := encoder.WGSportSignData(queryString)

	// 將簽名加入到 user 的 Sign 欄位中
	req.Sign = signature

	// 將更新後的 user 轉換為 JSON 格式
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("json.Marshal error %s", err.Error()))
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("json.UnMarshal error %s", err.Error()))
			return err
		}
		return nil
	}

	var header = map[string]string{}
	err = apicaller.WgGameSendPostRequest(url, header, body, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求發生錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

func GameAccountExist(username string) (*CheckUserRes, error) {
	if !WG_SPORT_SECRECT_INFO.RefreshTime.After(time.Now()) || !WG_SPORT_SECRECT_INFO.Ischeck {
		Init()
	}
	// 確認遊戲商帳號是否存在
	resp, err := checkUserAPI(username)
	if err != nil {
		mlog.Error(err.Error())
		return nil, err
	}
	if resp.Error_code == "CLIENT_NOT_EXIST" {
		mlog.Info("遊戲商帳號不存在")
		return &resp, errors.New("遊戲商帳號不存在")
	}

	return &resp, nil
}
