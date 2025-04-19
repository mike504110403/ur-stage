package sa_live

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"game_service/internal/database"
	"game_service/internal/service"
	"strconv"
	"time"

	"github.com/google/uuid"
	mlog "github.com/mike504110403/goutils/log"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"

	"github.com/valyala/fasthttp"
)

type SALiveGamer struct {
}

// var loginURL = "https://www.sai.slgaming.net/app.aspx?"
// var lobbyCode = "A11024"
// var lang = "zh-hans"
// var secretKey = "25699DAA28A74D87AE3F074D6E65FBDC"
// var md5Key = "GgaIMaiNNtg"
// var encryptKey = "g9G16nTs"
// var currencyType = "USDT"

// 加入遊戲
func (m *SALiveGamer) JoinGame(member service.MemberGameAccountInfo) (string, error) {
	url := ""
	url, err := loginAPI(member.UserName)
	if err != nil {
		mlog.Error(fmt.Sprintf("加入遊戲失敗: %s", err.Error()))
		return url, err
	}
	return url, nil
}

// 檢查遊戲商帳號是否存在
func GameAccountExist(username string) (*VerifyUsernameRes, error) {
	// 確認遊戲商帳號是否存在
	resp, err := verifyUserNameAPI(username)
	if err != nil {
		return &resp, err
	}
	if resp.IsExist {
		return &resp, errors.New("遊戲商帳號已存在")
	}

	return &resp, err
}

func AccountRegister(member *service.MemberGameAccountInfo) error {
	tx, err := database.MEMBER.TX()
	if err != nil {
		return err
	}

	member.GamePassword = uuid.New().String()[:20]
	// 註冊會員
	if err = service.CreateGameUser(tx, member.MemberId, *member); err != nil {
		return err
	}

	resp, err := createUserAPI(member.UserName)
	if err != nil {
		tx.Rollback()
		return err
	}

	if resp.ErrorMsgId != 0 {
		mlog.Error(fmt.Sprintf("註冊失敗: %s", resp.ErrorMsg))
		return errors.New(resp.ErrorMsg)
	}

	return tx.Commit()
}

// 創建會員
func createUserAPI(username string) (RegUserInfoRes, error) {
	if !SA_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !SA_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	var res RegUserInfoRes
	time := time.Now().Format("20060102150405")
	req := RegUserInfoReq{
		Method:       string(RegUserInfoUri),
		Key:          se.SecretKey,
		Time:         time,
		Username:     username,
		CurrencyType: se.CurrencyType,
	}

	// 構建查詢語句 QS
	qs, err := ToQueryString(req)
	if err != nil {
		return res, err
	}

	// DES 加密 QS
	encryptedQS, err := DESEncryptToBase64([]byte(qs), []byte(se.EncryptKey))
	if err != nil {
		return res, err
	}

	// 構建簽名
	signature := BuildMD5(qs + se.Md5Key + req.Time + se.SecretKey)
	// 構建 POST 請求
	if err := apiCaller(encryptedQS, signature, se.HttpDomain, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("註冊失敗")
		}
		if err := xml.Unmarshal(r.Body(), &res); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return res, err
	}
	return res, nil
}

// 確認會員重複
func verifyUserNameAPI(username string) (VerifyUsernameRes, error) {
	if !SA_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !SA_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	var res VerifyUsernameRes
	time := time.Now().Format("20060102150405")
	req := VerifyUsernameReq{
		Method:   string(VerifyUsernameUri),
		Key:      se.SecretKey,
		Time:     time,
		Username: username,
	}

	// 構建查詢語句 QS
	qs, err := ToQueryString(req)
	if err != nil {
		return res, err
	}

	// DES 加密 QS
	encryptedQS, err := DESEncryptToBase64([]byte(qs), []byte(se.EncryptKey))
	if err != nil {
		return res, err
	}

	// 構建簽名
	signature := BuildMD5(qs + se.Md5Key + req.Time + se.SecretKey)
	// 構建 POST 請求
	if err := apiCaller(encryptedQS, signature, se.HttpDomain, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("帳號確認失敗")
		}
		if err := xml.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.ErrorMsgId != 0 {
				return errors.New(res.ErrorMsg)
			}
		}
		return nil
	}); err != nil {
		return res, err
	}
	return res, nil
}

// 進入遊戲
func loginAPI(username string) (string, error) {
	if !SA_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !SA_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	var res LoginRes
	time := time.Now().Format("20060102150405")
	req := LoginReq{
		Method:       string(LoginRequestUri),
		Key:          se.SecretKey,
		Time:         time,
		Username:     username,
		CurrencyType: se.CurrencyType,
	}

	// 構建查詢語句 QS
	qs, err := ToQueryString(req)
	if err != nil {
		return "", err
	}

	// DES 加密 QS
	encryptedQS, err := DESEncryptToBase64([]byte(qs), []byte(se.EncryptKey))
	if err != nil {
		return "", err
	}

	// 構建簽名
	signature := BuildMD5(qs + se.Md5Key + req.Time + se.SecretKey)
	// 構建 POST 請求
	if err := apiCaller(encryptedQS, signature, se.HttpDomain, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("進入遊戲失敗")
		}
		if err := xml.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.ErrorMsgId != 0 {
				return errors.New(res.ErrorMsg)
			}
		}
		return nil
	}); err != nil {
		return "", err
	}

	// 組合登入 URL
	urlParams := LoginUrl{
		Username: username,
		Token:    res.Token,
		Lobby:    se.LobbyCode,
		Lang:     se.Lang,
	}
	loginUrl, err := ToQueryString(urlParams)
	if err != nil {
		return "", err
	}
	url := se.LoginURL + loginUrl

	return url, nil
}

// 攜入點數
func CreditBalanceDVAPI(username string, amount float64) (CreditBalanceDVRes, error) {
	if !SA_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !SA_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	var res CreditBalanceDVRes
	time := time.Now().Format("20060102150405")
	orderId := "IN" + time + username
	req := CreditBalanceDVReq{
		Method:       string(CreditBalanceDVUri),
		Key:          se.SecretKey,
		Time:         time,
		Username:     username,
		OrderId:      orderId,
		CreditAmount: amount,
		CurrencyType: se.CurrencyType,
	}

	// 構建查詢語句 QS
	qs, err := ToQueryString(req)
	if err != nil {
		return res, err
	}

	// DES 加密 QS
	encryptedQS, err := DESEncryptToBase64([]byte(qs), []byte(se.EncryptKey))
	if err != nil {
		return res, err
	}

	// 構建簽名
	signature := BuildMD5(qs + se.Md5Key + req.Time + se.SecretKey)
	// 構建 POST 請求
	if err := apiCaller(encryptedQS, signature, se.HttpDomain, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("攜入點數失敗")
		}
		if err := xml.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.ErrorMsgId != 0 {
				return errors.New(res.ErrorMsg)
			}
		}
		return nil
	}); err != nil {
		return res, err
	}
	return res, nil
}

func KickOut(member service.MemberGameAccountInfo) error {
	if _, err := kickUserAPI(member.UserName); err != nil {
		return err
	}
	return nil
}

// 踢出用戶
func kickUserAPI(username string) (string, error) {
	if !SA_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !SA_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	var res KickUserRes
	time := time.Now().Format("20060102150405")
	req := KickUserReq{
		Method:   string(KickUserUri),
		Key:      se.SecretKey,
		Time:     time,
		Username: username,
	}

	// 構建查詢語句 QS
	qs, err := ToQueryString(req)
	if err != nil {
		return "", err
	}

	// DES 加密 QS
	encryptedQS, err := DESEncryptToBase64([]byte(qs), []byte(se.EncryptKey))
	if err != nil {
		return "", err
	}

	// 構建簽名
	signature := BuildMD5(qs + se.Md5Key + req.Time + se.SecretKey)
	// 構建 POST 請求
	if err := apiCaller(encryptedQS, signature, se.HttpDomain, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("踢出失敗")
		}
		if err := xml.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.ErrorMsgId != 0 {
				return errors.New(res.ErrorMsg)
			}
		}
		return nil
	}); err != nil {
		return "", err
	}
	return res.ErrorMsg, nil
}

func GetBalance(member service.MemberGameAccountInfo) (float64, error) {
	memberStatus, err := getUserStatusDVAPI(member.UserName)
	if err != nil {
		return 0, err
	}
	return memberStatus.Balance, nil
}

// 獲取會員狀態
func getUserStatusDVAPI(username string) (GetUserStatusDVRes, error) {
	if !SA_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !SA_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	var res GetUserStatusDVRes
	time := time.Now().Format("20060102150405")
	req := GetUserStatusDVReq{
		Method:   string(GetUserStatusDVUri),
		Key:      se.SecretKey,
		Time:     time,
		Username: username,
	}

	// 構建查詢語句 QS
	qs, err := ToQueryString(req)
	if err != nil {
		return res, err
	}

	// DES 加密 QS
	encryptedQS, err := DESEncryptToBase64([]byte(qs), []byte(se.EncryptKey))
	if err != nil {
		return res, err
	}

	// 構建簽名
	signature := BuildMD5(qs + se.Md5Key + req.Time + se.SecretKey)
	// 構建 POST 請求
	if err := apiCaller(encryptedQS, signature, se.HttpDomain, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("獲取會員狀態失敗")
		}
		if err := xml.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.ErrorMsgId != 0 {
				return errors.New(res.ErrorMsg)
			}
		}
		return nil
	}); err != nil {
		return res, err
	}
	return res, nil
}

// 攜出點數
func DebitBalanceDVAPI(username string, amount float64) (DebitBalanceDVRes, error) {
	if !SA_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !SA_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	var res DebitBalanceDVRes
	time := time.Now().Format("20060102150405")
	orderId := "OUT" + time + username
	req := DebitBalanceDVReq{
		Method:      string(DebitBalanceDVUri),
		Key:         se.SecretKey,
		Time:        time,
		Username:    username,
		OrderId:     orderId,
		DebitAmount: amount,
	}

	// 構建查詢語句 QS
	qs, err := ToQueryString(req)
	if err != nil {
		return res, err
	}

	// DES 加密 QS
	encryptedQS, err := DESEncryptToBase64([]byte(qs), []byte(se.EncryptKey))
	if err != nil {
		return res, err
	}

	// 構建簽名
	signature := BuildMD5(qs + se.Md5Key + req.Time + se.SecretKey)
	// 構建 POST 請求
	if err := apiCaller(encryptedQS, signature, se.HttpDomain, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("攜出點數失敗")
		}
		if err := xml.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.ErrorMsgId != 0 {
				return errors.New(res.ErrorMsg)
			}
		}
		return nil
	}); err != nil {
		return res, err
	}
	return res, nil
}

// 獲取投注紀錄
func getBetRecordAPI(frequency time.Duration) (GetAllBetDetailsForTimeIntervalDVRes, error) {
	if !SA_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !SA_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}
	var res GetAllBetDetailsForTimeIntervalDVRes
	timeStr := time.Now().Format("20060102150405")
	req := GetAllBetDetailsForTimeIntervalDVReq{
		Method:   string(GetAllBetDetailsForTimeIntervalDVUri),
		Key:      se.SecretKey,
		Time:     timeStr,
		FromTime: time.Now().Add((-frequency) * 2).Format("2006-01-02 15:04:05"),
		ToTime:   time.Now().Add(+(1 * time.Minute)).Format("2006-01-02 15:04:05"),
	}

	// 構建查詢語句 QS
	qs, err := ToQueryString(req)
	if err != nil {
		return res, err
	}

	// DES 加密 QS
	encryptedQS, err := DESEncryptToBase64([]byte(qs), []byte(se.EncryptKey))
	if err != nil {
		return res, err
	}

	// 構建簽名
	signature := BuildMD5(qs + se.Md5Key + req.Time + se.SecretKey)
	// 構建 POST 請求
	if err := apiCaller(encryptedQS, signature, se.HttpDomain, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("http status code is not 200")
		}
		if err := xml.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.ErrorMsgId != 0 {
				return errors.New(res.ErrorMsg)
			}
		}
		return nil
	}); err != nil {
		return res, err
	}
	return res, nil
}

// 真人遊戲商取得注單
func GetLiveBetRecord(frequency time.Duration) error {
	// 遊戲類型
	agentType := typeparam.TypeParam{
		MainType: "game_type",
		SubType:  "live",
	}
	gameType, err := agentType.Get()
	if err != nil {
		return err
	}
	res, err := getBetRecordAPI(frequency)
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
	tx, err := database.GAME.TX()
	if err != nil {
		return err
	}
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}
	records := []service.BetRecord{}
	slecStr := `
		SELECT id
		FROM Member
		WHERE username = ? AND member_state = ?;
	`
	for _, v := range res.BetDetailList {
		var mid int
		if err := db.QueryRow(slecStr, v.Username, state).Scan(&mid); err != nil {
			continue
		}
		// 解析字串為 time.Time 類型
		parsedTime, err := time.ParseInLocation("2006-01-02T15:04:05.999", v.BetTime, time.Local)
		if err != nil {
			continue
		}

		uniqueId := strconv.FormatInt(v.BetID, 10)
		info, err := json.Marshal(v)
		if err != nil {
			continue
		}
		record := service.BetRecord{
			MemberId:   mid,
			AgentId:    AGENT_ID,
			BetUnique:  "sa-live-" + uniqueId,
			GameTypeId: gameType,
			BetAt:      parsedTime,
			Bet:        v.BetAmount,
			EffectBet:  v.Rolling,
			WinLose:    v.ResultAmount,
			BetInfo:    string(info),
		}
		records = append(records, record)
	}
	return service.WriteBetRecord2(tx, records, func() error {
		return tx.Commit()
	})
}
