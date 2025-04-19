package rsg_elec

import (
	"encoding/json"
	"errors"
	"fmt"
	"game_service/internal/apicaller"
	"game_service/internal/database"
	"game_service/internal/service"
	encoder "game_service/pkg/encoder"
	"strconv"
	"time"

	"github.com/google/uuid"
	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

type RSGElecGamer struct {
}

// 加入遊戲
func (m *RSGElecGamer) JoinGame(member service.MemberGameAccountInfo) (string, error) {
	url := ""

	url, err := LoginLobby(member.UserName)
	if err != nil {
		mlog.Error(fmt.Sprintf("加入遊戲失敗: %s", err.Error()))
		return url, err
	}
	return url, nil
}

func (m *RSGElecGamer) CheckPoint(agentId int, member service.MemberGameAccountInfo) (int, error) {
	return 1, nil
}

func (m *RSGElecGamer) LeaveGame() error {
	return nil
}

func (m *RSGElecGamer) GetBetRecord(agid int, starttime time.Time, endtime time.Time) (int, error) {
	return 1, nil
}

func (m *RSGElecGamer) PointOut(mid string, agentId int, point int) (int, error) {
	return 1, nil
}

func KickOut(member service.MemberGameAccountInfo) error {
	if _, err := kickOutAPI(member.UserName); err != nil {
		return err
	}
	return nil
}

// 存款
func PointIn(username string, amount float64) (DepositRes, error) {
	if !RSG_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !RSG_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	req := DepositReq{
		SystemCode:    ce.SystemCode,
		WebId:         ce.WebID,
		UserId:        username,
		TransactionID: randStringBytes(20),
		Currency:      "USDT",
		Balance:       amount,
	}
	url := ce.Httpdomain + string(PointInUri)
	resp := DepositRes{}

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

// 提款
func PointOut(username string, amount float64) (WithdrawRes, error) {
	if !RSG_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !RSG_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	req := WithdrawReq{
		SystemCode:    ce.SystemCode,
		WebId:         ce.WebID,
		UserId:        username,
		TransactionID: randStringBytes(20),
		Currency:      "USDT",
		Balance:       amount,
	}
	url := ce.Httpdomain + string(PointOutUri)
	resp := WithdrawRes{}

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

// 獲取餘額
func GetBalance(username string) (float64, error) {
	if !RSG_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !RSG_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}
	req := GetBalanceReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     username,
		Currency:   "USDT",
	}
	url := ce.Httpdomain + string(GetBalanceUri)
	resp := GetBalanceRes{}

	reqData, err := reqDataProcessing(req)
	if err != nil {
		return 0, err
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
		return 0, err
	}
	return resp.Data.CurrentPlayerBalance, nil
}

// 加入遊戲大廳
func LoginLobby(username string) (string, error) {
	if !RSG_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !RSG_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	req := EnterLobbyGameReq{
		SystemCode: ce.SystemCode,
		WebId:      ce.WebID,
		UserId:     username,
		UserName:   username,
		Currency:   "USDT",
		Language:   "zh-TW",
	}
	url := ce.Httpdomain + string(LoginLobbyUri)
	resp := EnterLobbyGameRes{}

	reqData, err := reqDataProcessing(req)
	if err != nil {
		return "", err
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
		return "", err
	}
	return resp.Data.URL, nil
}

// RegisterValid 註冊檢查 應檢查註冊資訊格式是否正確 並非檢查帳號是否存在
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
	if resp.Response.ErrorCode != 0 {
		tx.Rollback()
		return errors.New("註冊失敗")
	}

	return tx.Commit()
}

func GetBetRecord(frequency time.Duration) error {
	if !RSG_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !RSG_ELEC_SECRECT_INFO.Ischeck {
		Init()
	}

	resp, err := getGameDetailAPI(frequency)
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
	for _, v := range resp.Data.GameDetail {
		var mid int
		if err = db.QueryRow(slecStr, v.UserId, state).Scan(&mid); err != nil {
			return err
		}
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", v.PlayTime, time.Local)
		if err != nil {
			continue
		}
		betUnique := strconv.Itoa(int(v.SequenNumber))
		info, err := json.Marshal(v)
		if err != nil {
			mlog.Error(fmt.Sprintf("解析注單資訊失敗: %s", err.Error()))
			continue
		}
		winLose := v.WinAmt - v.BetAmt
		record := service.BetRecord{
			MemberId:   mid,
			AgentId:    AGENT_ID,
			BetUnique:  "rsg-" + betUnique,
			GameTypeId: gameType,
			BetAt:      parsedTime,
			Bet:        v.BetAmt,
			EffectBet:  v.BetAmt,
			WinLose:    winLose,
			BetInfo:    string(info),
		}
		records = append(records, record)
	}
	return service.WriteBetRecord2(tx, records, func() error {
		return tx.Commit()
	})
}
