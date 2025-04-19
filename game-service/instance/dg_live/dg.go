package dg_live

import (
	"encoding/json"
	"errors"
	"fmt"
	"game_service/internal/cachedata"
	"game_service/internal/database"
	"game_service/internal/service"
	"strings"
	"time"

	redisDriver "game_service/internal/redis"

	"github.com/gomodule/redigo/redis"

	"github.com/google/uuid"
	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

type DGLiveGamer struct {
}

// 投注記錄
func GetBetRecord(frequency time.Duration) error {
	report, err := GetReport()
	if err != nil {
		return err
	}
	gameNameIdMap := cachedata.AgentNameIdMap()
	AGENT_ID := gameNameIdMap["dg_live"]
	// 會員啟用狀態
	enableTpye := typeparam.TypeParam{
		MainType: "member_state",
		SubType:  "enable",
	}
	state, err := enableTpye.Get()
	if err != nil {
		return err
	}
	agentType := typeparam.TypeParam{
		MainType: "game_type",
		SubType:  "live",
	}
	gameType, err := agentType.Get()
	if err != nil {
		return err
	}
	records, recordIds := []service.BetRecord{}, []int{}
	slecStr := `
		SELECT id
		FROM Member
		WHERE username = ? COLLATE utf8mb4_general_ci AND member_state = ?;
	`
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}
	for _, v := range report.ReportList {
		var mid int
		if err := db.QueryRow(slecStr, v.UserName, state).Scan(&mid); err != nil {
			continue
		}
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", v.BetTime, time.Local)
		if err != nil {
			continue
		}
		info, err := json.Marshal(v)
		if err != nil {
			continue
		}
		if v.IsRevocation != 1 {
			continue
		}
		record := service.BetRecord{
			MemberId:   mid,
			AgentId:    AGENT_ID,
			BetUnique:  "dg-" + v.Ext,
			GameTypeId: gameType,
			BetAt:      parsedTime,
			Bet:        v.BetPoints,
			EffectBet:  v.AvailableBet,
			WinLose:    v.WinLose - v.BetPoints,
			BetInfo:    string(info),
		}
		records = append(records, record)
		recordIds = append(recordIds, v.Id)
	}

	tx, err := database.GAME.TX()
	if err != nil {
		return err
	}
	return service.WriteBetRecord2(tx, records, func() error {
		if err := MarkReport(recordIds); err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

// 抓取注單紀錄
func GetReport() (ReportRes, error) {
	res := ReportRes{}
	if err := apiCaller(ReportUri, nil, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("踢出失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.CodeId != 0 {
				return errors.New(res.Msg)
			}
		}
		return nil
	}); err != nil {
		return res, err
	}
	return res, nil
}

// 標記注單
func MarkReport(list []int) error {
	res := MarkReportRes{}
	req := map[string]interface{}{
		"list": list,
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return apiCaller(MarkReportUri, &reqBody, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("踢出失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.CodeId != 0 {
				return errors.New(res.Msg)
			}
		}
		return nil
	})
}

// 加入遊戲
func (m *DGLiveGamer) JoinGame(member service.MemberGameAccountInfo) (string, error) {
	url := ""

	url, err := Login(member)
	if err != nil {
		mlog.Error(fmt.Sprintf("加入遊戲失敗: %s", err.Error()))
		return url, err
	}
	return url, nil
}

func (m *DGLiveGamer) CheckPoint(agentId int, member service.MemberGameAccountInfo) (int, error) {
	return 1, nil
}

func (m *DGLiveGamer) LeaveGame() error {
	return nil
}

func (m *DGLiveGamer) GetBetRecord(agid int, starttime time.Time, endtime time.Time) (int, error) {
	return 1, nil
}

func (m *DGLiveGamer) PointOut(mid string, agentId int, point int) (int, error) {
	return 1, nil
}

// 註冊遊戲帳號
func AccountRegister(member *service.MemberGameAccountInfo) error {
	tx, err := database.MEMBER.TX()
	if err != nil {
		mlog.Error(fmt.Sprintf("資料庫連線錯誤: %s", err.Error()))
		return err
	}
	// 註冊會員
	member.GamePassword = uuid.New().String()
	err = service.CreateUser(tx, member.MemberId, *member)
	if err != nil {
		mlog.Info(fmt.Sprintf("會員註冊寫入失敗: %s", err.Error()))
		return err
	}

	req := map[string]interface{}{
		"username":     member.UserName,
		"password":     GenerateMD5(member.GamePassword),
		"currencyName": "USDT",
		"winLimit":     0,
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}
	res := SignUpRes{}
	if err := apiCaller(SignUpUri, &reqBody, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("註冊失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.CodeId != 0 {
				return errors.New(res.Msg)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return tx.Commit()
}

// 踢出遊戲
func KickOut(member service.MemberGameAccountInfo) error {
	conn := redisDriver.GetRedisConn()
	defer conn.Close()
	var memberIdReq int
	// 取得會員ID
	if memberId, err := GetMemberId(conn, member.UserName); err != nil {
		// 取讀不到重打一次 online
		if err == redis.ErrNil {
			Online(conn)
			if memberId, err = GetMemberId(conn, member.UserName); err != nil {
				return err
			} else {
				memberIdReq = memberId
			}
		}
	} else {
		memberIdReq = memberId
	}

	req := OfflineReq{}
	req.OffLineP.List = append(req.OffLineP.List, memberIdReq)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}
	res := OfflineRes{}

	return apiCaller(OfflineUri, &reqBody, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("踢出失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.CodeId != 0 {
				return errors.New(res.Msg)
			}
		}
		return nil
	})
}

// 取得會員ID
func GetMemberId(conn redis.Conn, username string) (int, error) {
	if memberId, err := redis.Int(conn.Do("HGET", "dg-online-member", strings.ToLower(username))); err != nil {
		return memberId, err
	} else {
		return memberId, nil
	}
}

// 確認餘額
func CheckBlance(member service.MemberGameAccountInfo) (float64, error) {
	var blance float64

	req := BalanceReq{}
	req.UserName = strings.ToUpper(member.UserName)
	reqBody, err := json.Marshal(req)
	if err != nil {
		return blance, err
	}
	res := BalanceRes{}

	if err := apiCaller(BalanceUri, &reqBody, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("踢出失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.CodeId != 0 {
				return errors.New(res.Msg)
			}
		}
		return nil
	}); err != nil {
		return blance, err
	} else {
		blance = res.Balance
	}

	return blance, nil
}

// 存提
func TransferPoint(member service.MemberGameAccountInfo, point float64) error {
	req := map[string]interface{}{
		"username": member.UserName,
		"amount":   point,
		"serial":   uuid.New().String(),
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}
	res := TransferRes{}

	return apiCaller(TransferUri, &reqBody, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("踢出失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.CodeId != 0 {
				return errors.New(res.Msg)
			}
		}
		return nil
	})
}

// 登入
func Login(member service.MemberGameAccountInfo) (string, error) {
	url := ""
	req := map[string]interface{}{
		"username":     member.UserName,
		"password":     GenerateMD5(member.GamePassword),
		"currencyName": "USDT",
		"winLimit":     0,
		"language":     "cn",
		"limitGroup":   "A",
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return url, err
	}
	res := LoginRes{}

	if err := apiCaller(LoginUri, &reqBody, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("踢出失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		} else {
			if res.CodeId != 0 {
				return errors.New(res.Msg)
			}
		}
		return nil
	}); err != nil {
		return url, err
	} else {
		url = res.List[0]
	}

	return url, nil
}

// 在線會員
func Online(conn redis.Conn) error {
	res := OnlineRes{}

	return apiCaller(OnlineUri, nil, func(r *fasthttp.Response) error {
		if r.StatusCode() != 200 {
			return errors.New("踢出失敗")
		}
		if err := json.Unmarshal(r.Body(), &res); err != nil {
			return err
		}

		for _, v := range res.List {
			if _, err := conn.Do(
				"HSET",
				"dg-online-member",
				strings.ToLower(v.UserName),
				v.MemberId,
			); err != nil {
				mlog.Error(fmt.Sprintf("寫入redis失敗: %s", err.Error()))
				continue
			}
		}
		return nil
	})
}
