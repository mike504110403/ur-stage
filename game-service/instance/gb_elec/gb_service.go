package gb_elec

import (
	"encoding/json"
	"errors"
	"fmt"
	"game_service/internal/apicaller"
	"game_service/internal/cachedata"
	"game_service/internal/database"
	"game_service/internal/service"
	"game_service/pkg/encoder"
	"strconv"
	"time"

	"github.com/google/uuid"
	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

var GB_ELEC_SECRECT_INFO = GbElecSecrect_Info{}
var sc = GB_ELEC_SECRECT_INFO.SC
var AGENT_ID int

func Init() {
	cacheSecrect := cachedata.AgentSecretMap()
	cacheId := cachedata.AgentIdMap()
	gbelecSecret, ok := cacheSecrect["gb_elec"]
	if !ok {
		mlog.Error("gb_elec secret not found")
		return
	}
	gbelecId, ok := cacheId["gb_elec"]
	if !ok {
		mlog.Error("gb_elec id not found")
		return
	}
	// 取得代理商ID
	agentId, err := strconv.Atoi(gbelecId)
	if err != nil {
		mlog.Error(err.Error())
		return
	}
	AGENT_ID = agentId
	// 取得代理商資訊
	err = json.Unmarshal([]byte(gbelecSecret), &GB_ELEC_SECRECT_INFO.SC)
	if err != nil {
		return
	}

	GB_ELEC_SECRECT_INFO.RefreshTime = time.Now().Add(time.Minute * 30)
	sc = GB_ELEC_SECRECT_INFO.SC
}

// 註冊遊戲商帳號
func createUserAPI(username string) (RegisterRes, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	req := RegisterReq{
		Action:    string(RegisterUri),
		AppID:     sc.AppID,
		AppSecret: sc.AppSecret,
		Uid:       username,
		SignKey:   sc.Sign_key,
	}
	resp := RegisterRes{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}

	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}

		err := json.Unmarshal(r.Body(), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 帳號註冊
func AccountRegister(member *service.MemberGameAccountInfo, agentId int) error {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	tx, err := database.MEMBER.TX()
	if err != nil {
		return err
	}
	member.GamePassword = uuid.New().String()[:20]
	// 註冊帳號
	err = service.CreateDBUser(tx, member.MemberId, *member, agentId)
	if err != nil {
		return err
	}
	if resp, err := createUserAPI(member.UserName); err != nil {
		tx.Rollback()
		return err
	} else {
		if resp.ReturnCode != "0000" {
			tx.Rollback()
			return errors.New("遊戲商創建帳號失敗")
		}
	}

	return tx.Commit()
}

// 進入指定遊戲獲取Url
func loginAPI(uid string) (LoginDTO, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	req := LoginReq{
		Action:    string(LoginUri),
		AppID:     sc.AppID,
		AppSecret: sc.AppSecret,
		Uid:       uid,
		GameCode:  "QHSlotsRes",
		SignKey:   sc.Sign_key,
	}
	resp := LoginDTO{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 試玩帳號進入指定遊戲獲取Url
func demoLoginAPI() (DemoLoginDTO, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	req := DemoLoginReq{
		Action:    string(DemoLoginUri),
		AppID:     sc.AppID,
		AppSecret: sc.AppSecret,
		GameCode:  "QHSlotsRes",
		SignKey:   sc.Sign_key,
	}
	resp := DemoLoginDTO{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}

	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	if err := apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler); err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

func lobbyLoginAPI(uid string) (string, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	req := LobbyLoginReq{
		Action:    string(LobbyLoginUri),
		AppID:     sc.AppID,
		AppSecret: sc.AppSecret,
		Uid:       uid,
		SignKey:   sc.Sign_key,
	}
	resp := LobbyLoginDTO{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return "", err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return "", err
	}

	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return "", err
	}
	return resp.Path, nil
}

func CheckPoint(member service.MemberGameAccountInfo) (float64, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	resp, err := getMoneyAPI(member.UserName)
	if err != nil {
		return 0, err
	}
	amount, err := strconv.ParseFloat(resp.Amount, 64)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

// 獲取玩家當前餘額
func getMoneyAPI(uid string) (GetMoneyDTO, error) {

	req := GetMoneyReq{
		Action:    string(GetMoneyUri),
		AppID:     sc.AppID,
		AppSecret: sc.AppSecret,
		Uid:       uid,
		SignKey:   sc.Sign_key,
	}
	resp := GetMoneyDTO{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		fmt.Println(string(r.Body()))
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler)
	if err != nil {
		mlog.Error(err.Error())
		return resp, err
	}
	return resp, nil
}

func TransferOutPoint(member service.MemberGameAccountInfo, point float64) error {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	_, err := transferOutAPI(member.UserName, strconv.FormatFloat(point, 'f', -1, 64))
	if err != nil {
		return err
	}
	return nil
}

func TransferPoint(member service.MemberGameAccountInfo, point float64) error {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	_, err := transferAPI(member.UserName, strconv.FormatFloat(point, 'f', -1, 64))
	if err != nil {
		return err
	}
	return nil
}

// 攜出玩家錢包
func transferOutAPI(uid string, amount string) (TransferOutDTO, error) {

	orderNo := uuid.New().String()
	req := TransferReq{
		Action:    string(TransferUri),
		AppID:     sc.AppID,
		AppSecret: sc.AppSecret,
		Uid:       uid,
		OrderNo:   orderNo,
		Amount:    amount,
		SignKey:   sc.Sign_key,
	}
	resp := TransferOutDTO{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}

	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		// err := json.Unmarshal([]byte(data), &resp)
		// if err != nil {
		// 	mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
		// 	return err
		// }
		return nil
	}
	err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 攜入玩家錢包
func transferAPI(uid string, amount string) (TransferDTO, error) {

	orderNo := uuid.New().String()
	req := TransferReq{
		Action:    string(TransferUri),
		AppID:     sc.AppID,
		AppSecret: sc.AppSecret,
		Uid:       uid,
		OrderNo:   orderNo,
		Amount:    amount,
		SignKey:   sc.Sign_key,
	}
	resp := TransferDTO{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}

	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 儲值訂單查詢
func getTransferStateAPI(req GetTransferStateReq) (GetTransferStateDTO, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	resp := GetTransferStateDTO{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}

	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler)
	if err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 回傳資料屬性判斷
func ParseGetTransferStateData(data []byte) string {
	var resp map[string]interface{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		mlog.Error(fmt.Sprintf("回傳資料解析失敗 %s", err.Error()))
		return ""
	}
	if resp["returnCode"] != "0000" {
		return ""
	}

	// 判断是否存在 "data" 字段
	dataField, ok := resp["data"]
	if !ok {
		fmt.Println("data field not found")
		return ""
	}
	switch v := dataField.(type) {
	case map[string]interface{}, []interface{}:
		orderData, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(orderData)
	default:
		return ""
	}
}

func GetBetRecord(frequency time.Duration) error {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	// 確認遊戲商帳號是否存在
	resp, err := getDetailedAPI(
		GetDetailedReq{
			Action:    string(GetDetailedUri),
			AppID:     sc.AppID,
			AppSecret: sc.AppSecret,
			StartTime: time.Now().Add((-frequency) * 2).Format("2006-01-02 15:04:05"),
			EndTime:   time.Now().Add(+(1 * time.Minute)).Format("2006-01-02 15:04:05"),
			SignKey:   sc.Sign_key,
		},
	)
	if err != nil {
		return err
	}

	// 會員啟用狀態
	enableTpye := typeparam.TypeParam{
		MainType: "member_state",
		SubType:  "enable",
	}
	state, err := enableTpye.Get()
	if err != nil {
		return err
	}
	// 代理商遊戲類型
	agentType := typeparam.TypeParam{
		MainType: "game_type",
		SubType:  "elect",
	}
	gameType, err := agentType.Get()
	if err != nil {
		return err
	}

	tx, err := database.GAME.TX()
	if err != nil {
		return err
	}

	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}

	slecStr := `
		SELECT id
		FROM Member
		WHERE username = ? AND member_state = ?;
	`
	records := []service.BetRecord{}
	for _, v := range resp.List {
		var mid int
		if err = db.QueryRow(slecStr, v.UID, state).Scan(&mid); err != nil {
			return err
		}
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", v.GameDate, time.Local)
		if err != nil {
			continue
		}
		bet, err := strconv.ParseFloat(v.Bet, 64)
		if err != nil {
			continue
		}
		effectBet, err := strconv.ParseFloat(v.ValidBet, 64)
		if err != nil {
			continue
		}
		winlose, err := strconv.ParseFloat(v.NetWin, 64)
		if err != nil {
			continue
		}
		info, err := json.Marshal(v)
		if err != nil {
			mlog.Error(fmt.Sprintf("解析注單資訊失敗: %s", err.Error()))
			continue
		}
		record := service.BetRecord{
			MemberId:   mid,
			AgentId:    AGENT_ID,
			BetUnique:  "gb-" + v.No,
			GameTypeId: gameType,
			BetAt:      parsedTime,
			Bet:        bet,
			EffectBet:  effectBet,
			WinLose:    winlose,
			BetInfo:    string(info),
		}
		records = append(records, record)
	}
	return service.WriteBetRecord2(tx, records, func() error {
		return tx.Commit()
	})
}

// 查詢遊戲詳細信息
func getDetailedAPI(req GetDetailedReq) (GetDetailedDataDto, error) {
	resp := GetDetailedDataDto{}
	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	if err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler); err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 獲取遊戲列表
func getGameListAPI() (GetGameListDTO, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	req := GetGameListReq{
		Action:    string(GetGameListUri),
		AppID:     sc.AppID,
		AppSecret: sc.AppSecret,
		SignKey:   sc.Sign_key,
	}
	resp := GetGameListDTO{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	if err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler); err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 獲取遊戲歷史轉帳記錄
func getHistoryTransferAPI(req GetHistoryTransferReq) (GetHistoryTransferDTO, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	resp := GetHistoryTransferDTO{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	if err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler); err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 注單統計
func getOrderStatAPI(req GetOrderStatReq) (GetOrderStatDTO, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	resp := GetOrderStatDTO{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	if err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler); err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 玩家活動列表
func getActivityListsAPI(req ActivityListsReq) (ActivityListsDTO, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	resp := ActivityListsDTO{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	if err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler); err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

// 活動中獎名單
func getActivityWinnerListAPI(req ActivityWinnerListReq) (ActivityWinnerListData, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	resp := ActivityWinnerListData{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	if err = apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler); err != nil {
		return resp, err
	}
	return resp, nil
}

// 查詢玩家是否有未完成遊戲
func getUserGameStateAPI(uid string) (UserGameStateData, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	req := UserGameStateReq{
		Action:    "35",
		AppID:     sc.AppID,
		AppSecret: sc.AppSecret,
		Uid:       uid,
		SignKey:   sc.Sign_key,
	}
	resp := UserGameStateData{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}

	if err := apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler); err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

func getUserGameNoteAPI(uid string) (UserGameNoteRes, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	req := UserGameNoteReq{
		Action:    "36",
		AppID:     sc.AppID,
		AppSecret: sc.AppSecret,
		Uid:       uid,
		SignKey:   sc.Sign_key,
	}
	resp := UserGameNoteRes{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	if err := apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler); err != nil {
		return resp, err
	}
	return resp, nil
}

func KickOut(member service.MemberGameAccountInfo) (int, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	_, err := kickOutAPI(member.UserName)
	if err != nil {
		mlog.Error(fmt.Sprintf("踢出失敗: %s", err.Error()))
		return 0, err
	}
	return 1, nil
}

func kickOutAPI(uid string) (KickOutRes, error) {

	req := KickOutReq{
		Action:    "37",
		AppID:     sc.AppID,
		AppSecret: sc.AppSecret,
		Uid:       uid,
		SignKey:   sc.Sign_key,
	}
	resp := KickOutRes{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp.Data)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	if err := apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler); err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}

func getUserCompensationAmountRecordAPI(req GetUserCompensationAmountRecordReq) (GetUserCompensationAmountRecordRes, error) {
	if !GB_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !GB_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	resp := GetUserCompensationAmountRecordRes{}

	// 獲取MD5加密字符串
	queryString := encoder.GetMD5string(req)

	// 獲取進行MD5加密
	signature := encoder.SignatureData(queryString)

	// 設置簽名
	req.SignKey = signature
	body, err := json.Marshal(req)
	if err != nil {
		mlog.Error(fmt.Sprintf("資料解析錯誤 %s", err.Error()))
		return resp, err
	}

	// DES-ECB 加密
	encrypted, err := encoder.DESEncryptECBBase64(body, sc.EBC_key)
	if err != nil {
		mlog.Error(fmt.Sprintf("加密失敗 %s", err.Error()))
		return resp, err
	}
	var data string
	handler := func(r *fasthttp.Response) error {
		if r.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("HTTP status code: %d", r.StatusCode())
		}
		data = ParseGetTransferStateData(r.Body())
		if data == "" {
			mlog.Error(fmt.Sprintf("無回應資料 %s", data))
			return nil
		}
		err := json.Unmarshal([]byte(data), &resp)
		if err != nil {
			mlog.Error(fmt.Sprintf("Unmarshal錯誤 %s", err.Error()))
			return err
		}
		return nil
	}
	if err := apicaller.GBGameSendPostRequest(sc.Httpdomain, nil, encrypted, handler); err != nil {
		mlog.Error(fmt.Sprintf("請求錯誤 %s", err.Error()))
		return resp, err
	}
	return resp, nil
}
