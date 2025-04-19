package atg_elect

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func Register(username string) error {
	req := RegisterReq{
		Username: username,
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}
	res := RegisterRes{}

	if err := apiCallerPost(string(RegisterUri), &reqBody, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("遊戲商註冊失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func getBalance(username string) (string, error) {
	res := BalanceRes{}
	req := BalanceReq{
		Username: username,
	}
	urlParamsStr := StructToQueryString(req)
	uri := GameProvidersUri + ProviserId + BalanceUri
	if err := apiCallerGet(string(uri), urlParamsStr, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("獲取錢包資訊失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}
	return res.Data.Balance, nil
}

func getGameLobby(username string) (string, error) {
	res := GameLobbyRes{}
	req := GameLobbyReq{
		Username: username,
	}
	urlParamsStr := StructToQueryString(req)
	uri := GameProvidersUri + ProviserId + GameLobbyUri
	if err := apiCallerGet(string(uri), urlParamsStr, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("獲取大廳連結失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}
	return res.Data.Url, nil
}

func transferIn(username string, point float64) (TransferRes, error) {
	req := TransferReq{
		Username:   username,
		Balance:    point,
		Action:     "IN",
		TransferId: uuid.New().String(),
	}
	res := TransferRes{}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return res, err
	}
	uri := GameProvidersUri + ProviserId + BalanceUri

	if err := apiCallerPost(string(uri), &reqBody, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("攜入失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return res, err
	}
	return res, nil
}

func transferOut(username string, point float64) (TransferRes, error) {
	req := TransferReq{
		Username:   username,
		Balance:    point,
		Action:     "OUT",
		TransferId: uuid.New().String(),
	}
	res := TransferRes{}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return res, err
	}
	uri := GameProvidersUri + ProviserId + BalanceUri

	if err := apiCallerPost(string(uri), &reqBody, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("攜出失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return res, err
	}
	return res, nil
}

func kickOut(username string) error {

	res := TransferRes{}
	uri := string(PlayUri) + username

	if err := apiCallerDelete(uri, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("踢出失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func betRecord(req BetRecordReq) (BetRecordRes, error) {
	res := BetRecordRes{}
	formData, err := structToMap(req)
	if err != nil {
		return res, err
	}

	if err := apiCallerFormDataPost(string(TransactionUri), formData, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("取得注單失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return res, err
	}
	return res, nil

}

// func GetBetRecord(frequency time.Duration) error {
// 	if !RSG_ELEC_SECRECT_INFO.RefreshTime.After(time.Now()) || !RSG_ELEC_SECRECT_INFO.Ischeck {
// 		Init()
// 	}

// 	resp, err := betRecord(frequency)
// 	if err != nil {
// 		return err
// 	}

// 	// 會員啟用狀態
// 	enableTpye := typeparam.TypeParam{
// 		MainType: "member_state",
// 		SubType:  "enable",
// 	}
// 	state, err := enableTpye.Get()
// 	if err != nil {
// 		return err
// 	}
// 	// 代理商遊戲類型
// 	agentType := typeparam.TypeParam{
// 		MainType: "game_type",
// 		SubType:  "elect",
// 	}
// 	gameType, err := agentType.Get()
// 	if err != nil {
// 		return err
// 	}

// 	tx, err := database.GAME.TX()
// 	if err != nil {
// 		return err
// 	}

// 	db, err := database.MEMBER.DB()
// 	if err != nil {
// 		return err
// 	}

// 	slecStr := `
// 		SELECT id
// 		FROM Member
// 		WHERE username = ? AND member_state = ?;
// 	`
// 	records := []service.BetRecord{}
// 	for _, v := range resp {
// 		var mid int
// 		if err = db.QueryRow(slecStr, v.MemberName, state).Scan(&mid); err != nil {
// 			return err
// 		}
// 		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", v.BettingDate, time.Local)
// 		if err != nil {
// 			continue
// 		}
// 		effectBet, err := strconv.ParseFloat(v.ValidBet, 64)
// 		if err != nil {
// 			fmt.Println("轉換失敗:", err)
// 			continue
// 		}
// 		info, err := json.Marshal(v)
// 		if err != nil {
// 			mlog.Error(fmt.Sprintf("解析注單資訊失敗: %s", err.Error()))
// 			continue
// 		}

// 		record := service.BetRecord{
// 			MemberId:   mid,
// 			AgentId:    AGENT_ID,
// 			BetUnique:  "atg-" + v.BettingID,
// 			GameTypeId: gameType,
// 			BetAt:      parsedTime,
// 			Bet:        v.BettingAmount,
// 			EffectBet:  effectBet,
// 			WinLose:    v.WinLoseAmount,
// 			BetInfo:    string(info),
// 		}
// 		records = append(records, record)
// 	}
// 	return service.WriteBetRecord2(tx, records, func() error {
// 		return tx.Commit()
// 	})
// }
