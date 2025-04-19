package wg_sport

import (
	"encoding/json"
	"fmt"
	"game_service/internal/database"
	"game_service/internal/service"
	"strconv"
	"strings"
	"time"

	mlog "github.com/mike504110403/goutils/log"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

// WG_sport遊戲商
// WG_Game遊戲商
type WGSportGamer struct {
}

// 進入遊戲
func (g *WGSportGamer) JoinGame(member service.MemberGameAccountInfo) (string, error) {
	res := ""

	url, err := forwardGameAPI(member.UserName)
	if err != nil {
		mlog.Error(fmt.Sprintf("加入遊戲失敗: %s", err.Error()))
		return "加入遊戲失敗", err
	} else {
		res = url
	}
	return res, nil
}

// 遊戲商取得注單
func GetBetRecord(frequency time.Duration) error {
	if !WG_SPORT_SECRECT_INFO.RefreshTime.After(time.Now()) || !WG_SPORT_SECRECT_INFO.Ischeck {
		Init()
	}

	// 遊戲類型
	agentType := typeparam.TypeParam{
		MainType: "game_type",
		SubType:  "sport",
	}
	gameType, err := agentType.Get()
	if err != nil {
		return err
	}
	res, err := buyListGetAPI(BuyListGetReq{
		Prefix:   WSSI.Prefix,
		Type:     "3",
		Sdate:    time.Now().Add((-frequency) * 2).Format("2006-01-02 15:04:05"),
		Edate:    time.Now().Add(+(1 * time.Minute)).Format("2006-01-02 15:04:05"),
		Pri_type: "1",
	})
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
	for _, v := range res.Data {
		parts := strings.Split(v.Username, "_")
		username := ""
		if len(parts) > 1 {
			username = parts[1]
		} else {
			username = v.Username
		}
		var mid int
		if err := db.QueryRow(slecStr, username, state).Scan(&mid); err != nil {
			continue
		}
		// 解析字串為 time.Time 類型
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", v.UpdatedAt, time.Local)
		if err != nil {
			continue
		}

		uniqueId := strconv.FormatInt(int64(v.Id), 10)
		info, err := json.Marshal(v)
		if err != nil {
			continue
		}
		record := service.BetRecord{
			MemberId:   mid,
			AgentId:    AGENT_ID,
			BetUnique:  "wg-sport-" + uniqueId,
			GameTypeId: gameType,
			BetAt:      parsedTime,
			Bet:        v.Gold,
			EffectBet:  v.GoldOk,
			WinLose:    v.Result,
			BetInfo:    string(info),
		}
		records = append(records, record)
	}
	return service.WriteBetRecord2(tx, records, func() error {
		return tx.Commit()
	})
}
