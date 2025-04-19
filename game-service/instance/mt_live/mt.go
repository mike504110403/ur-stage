package mt_live

import (
	"encoding/json"
	"fmt"
	"game_service/internal/database"
	"game_service/internal/service"
	"strconv"
	"time"

	mlog "github.com/mike504110403/goutils/log"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

// TODO: 這邊才是對“遊戲”介面的實作 層次要抽開

// Apollo apollo遊戲商
type MTLiveGamer struct {
}

// 進入遊戲
func (m *MTLiveGamer) JoinGame(member service.MemberGameAccountInfo) (string, error) {
	res := ""

	url, err := getURLTokenAPI(member.UserName)
	if err != nil {
		mlog.Error(fmt.Sprintf("加入遊戲失敗: %s", err.Error()))
		return "加入遊戲失敗", err
	} else {
		res = url
	}
	return res, nil
}

// 離開遊戲
func (m *MTLiveGamer) LeaveGame() error {
	return nil
}

// 積分轉入
func (m *MTLiveGamer) PointIn(mid string, agentId int, point float64) (int, error) {
	// 確認錢包金額 打api

	db, err := database.MEMBER.DB()
	if err != nil {
		return 0, err
	}

	midInt, err := strconv.Atoi(mid)
	if err != nil {
		return 0, err
	}

	member, err := service.GetMember(db, agentId, midInt)
	if err != nil {
		return 0, err
	}
	ckbalance, err := CheckBalance(member.UserName)
	if err != nil {
		mlog.Error(err.Error())
		return 0, err
	}
	fmt.Println(ckbalance.Data.Balance)
	// 確認玩家是否有遊玩過，為遊玩導入註冊
	gameExistRes, err := GameAccountExist(member.UserName)
	if err != nil || gameExistRes.Data.Result == 0 {
		// 其他類型錯誤
		if gameExistRes == nil {
			return 0, err
		}
	}

	res, err := TransferIn(member.UserName, point)
	if err != nil {
		mlog.Error(err.Error())
		return 0, err
	}
	cknewbalance, err := CheckBalance(member.UserName)
	if err != nil {
		mlog.Error(err.Error())
		return 0, err
	}
	fmt.Println(cknewbalance.Data.Balance)
	// // 檢查儲值單號
	// cksn, err := FindTransferRecord(res.Data.TransferId)
	// if err != nil {
	// 	mlog.Error(err.Error())
	// 	return 0, err
	// }
	// if res.Data.TransferId != cksn.Data.List[0].SN {
	// 	mlog.Error("儲值單號不符")
	// 	return 0, err
	// }
	// if res.Data.Balance != cksn.Data.List[0].Amount {
	// 	mlog.Error("儲值金額不符")
	// 	return 0, err
	// }
	// ckbalance, err := CheckBalance(member.UserName)
	// if err != nil {
	// 	mlog.Error(err.Error())
	// 	return 0, err
	// }
	// if res.Data.Balance != ckbalance.Data.Balance {
	// 	mlog.Error("提款失敗")
	// 	return 0, err
	// }
	balance, err := strconv.ParseFloat(res.Data.Balance, 64)
	if err != nil {
		return 0, err
	}
	return int(balance), nil
}

// 積分轉出
func (m *MTLiveGamer) PointOut(mid string, agentId int, point float64) (int, error) {

	db, err := database.MEMBER.DB()
	if err != nil {
		return 0, err
	}

	midInt, err := strconv.Atoi(mid)
	if err != nil {
		return 0, err
	}

	member, err := service.GetMember(db, agentId, midInt)
	if err != nil {
		return 0, err
	}

	res, err := TransferOut(member.UserName, point)
	if err != nil {
		mlog.Error(err.Error())
		return 0, err
	}

	// 檢查儲值單號
	cksn, err := FindTransferRecord(res.Data.TransferId)
	if err != nil {
		mlog.Error(err.Error())
		return 0, err
	}
	if res.Data.TransferId != cksn.Data.List[0].SN {
		mlog.Error("儲值單號不符")
		return 0, err
	}
	if res.Data.Balance != cksn.Data.List[0].Amount {
		mlog.Error("儲值金額不符")
		return 0, err
	}
	ckbalance, err := CheckBalance(member.UserName)
	if err != nil {
		return 0, err
	}
	if ckbalance.Data.Balance != "0" {
		mlog.Error("錢包未攜出")
		return 0, err
	}
	return 1, nil
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
	return GetBetRecord(frequency, gameType)
}

// 真人donate取得注單
func GetDonateRecord(frequency time.Duration) error {
	// 遊戲類型
	agentType := typeparam.TypeParam{
		MainType: "game_type",
		SubType:  "live_donate",
	}
	gameType, err := agentType.Get()
	if err != nil {
		return err
	}
	return GetDonate(frequency, gameType)
}

// 取得注單
func GetBetRecord(frequency time.Duration, gameType int) error {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}

	resp, err := getBetRecordAPI(
		GetBetRecordReq{
			SystemCode: ce.SystemCode,
			WebId:      ce.WebID,
			StartTime:  time.Now().Add((-frequency) * 2).Format("2006-01-02 15:04:05"),
			EndTime:    time.Now().Add(+(1 * time.Minute)).Format("2006-01-02 15:04:05"),
		},
		ce.Httpdomain+string(GetBetRecordUri),
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
	for _, v := range resp.Data.List {
		var mid int
		if err := db.QueryRow(slecStr, v.UserId, state).Scan(&mid); err != nil {
			continue
		}
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", v.OrderTime, time.Local)
		if err != nil {
			continue
		}
		orderMoney, err := strconv.ParseFloat(v.OrderMoney, 64)
		if err != nil {
			continue
		}
		validMoney, err := strconv.ParseFloat(v.ValidMoney, 64)
		if err != nil {
			continue
		}
		winLose, err := strconv.ParseFloat(v.Profit, 64)
		if err != nil {
			continue
		}
		info, err := json.Marshal(v)
		if err != nil {
			continue
		}
		record := service.BetRecord{
			MemberId:   mid,
			AgentId:    AGENT_ID,
			BetUnique:  "mt-live-" + v.Sn,
			GameTypeId: gameType,
			BetAt:      parsedTime,
			Bet:        orderMoney,
			EffectBet:  validMoney,
			WinLose:    winLose,
			BetInfo:    string(info),
		}
		records = append(records, record)
	}
	return service.WriteBetRecord2(tx, records, func() error {
		return tx.Commit()
	})
}

// 取得注單
func GetDonate(frequency time.Duration, gameType int) error {
	if !MT_LIVE_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LIVE_SECRECT_INFO.Ischeck {
		Init()
	}

	resp, err := getDonateRecordAPI(
		DonateRecordReq{
			SystemCode: ce.SystemCode,
			WebID:      ce.WebID,
			StartTime:  time.Now().Add((-frequency) * 2).Format("2006-01-02 15:04:05"),
			EndTime:    time.Now().Add(+(1 * time.Minute)).Format("2006-01-02 15:04:05"),
		},
		ce.Httpdomain+string(GetDonateRecordUri),
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
	for _, v := range resp.Data.List {
		var mid int
		if err := db.QueryRow(slecStr, v.UserID, state).Scan(&mid); err != nil {
			continue
		}
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", v.Time, time.Local)
		if err != nil {
			continue
		}
		info, err := json.Marshal(v)
		if err != nil {
			continue
		}
		record := service.BetRecord{
			MemberId:   mid,
			AgentId:    AGENT_ID,
			BetUnique:  "mt-donate-" + v.SN,
			GameTypeId: gameType,
			BetAt:      parsedTime,
			Bet:        float64(v.Money),
			EffectBet:  float64(v.Money),
			WinLose:    0,
			BetInfo:    string(info),
		}
		records = append(records, record)
	}
	return service.WriteBetRecord2(tx, records, func() error {
		return tx.Commit()
	})
}
