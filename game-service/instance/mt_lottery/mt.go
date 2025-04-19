package mt_lottery

import (
	"encoding/json"
	"fmt"
	"game_service/internal/database"
	"game_service/internal/service"
	"strings"
	"time"

	mlog "github.com/mike504110403/goutils/log"
	"gitlab.com/gogogo2712128/common_moduals/typeparam"
)

// Apollo apollo遊戲商
type MTLotteryGamer struct {
}

// 進入遊戲
func (m *MTLotteryGamer) JoinGame(member service.MemberGameAccountInfo) (string, error) {
	url, err := getLoginUrlAPI(member)
	if err != nil {
		return "", err
	}
	return url, nil
}

// 傳入需更改
func GetBetOrder(frequency time.Duration) error {
	if !MT_LOTTERY_SECRECT_INFO.RefreshTime.After(time.Now()) || !MT_LOTTERY_SECRECT_INFO.Ischeck {
		Init()
	}
	// 會員啟用狀態參數
	enableType := typeparam.TypeParam{
		MainType: "member_state",
		SubType:  "enable",
	}
	enableTypeInt, err := enableType.Get()
	if err != nil {
		return err
	}
	// 遊戲類別參數
	gameType := typeparam.TypeParam{
		MainType: "game_type",
		SubType:  "lottery",
	}
	gameTypeInt, err := gameType.Get()
	if err != nil {
		return err
	}
	// 確認遊戲商帳號是否存在
	resp, err := getBetOrderV2API(
		BetOrderV2Req{
			StartTime:  time.Now().Add((-frequency) * 2).Unix(),
			EndTime:    time.Now().Add(+(1 * time.Minute)).Unix(),
			GameTypeId: 4,
			GameId:     1,
			DateType:   2,
		},
		hc.Httpdomain+string(BetOrderV2Uri),
	)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}

	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}
	tx, err := database.GAME.TX()
	if err != nil {
		return err
	}

	records := []service.BetRecord{}
	slecStr := `
		SELECT id
		FROM Member
		WHERE username = ? AND member_state = ?;
	`
	for _, v := range resp.Rows {
		parts := strings.Split(v.Account, "-")
		username := ""
		if len(parts) > 1 {
			username = parts[1]
		} else {
			username = v.Account
		}
		var mid int
		if err = db.QueryRow(slecStr, username, enableTypeInt).Scan(&mid); err != nil {
			mlog.Error(fmt.Sprintf("查詢會員失敗: %s", err.Error()))
			continue
		}
		parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", v.DateCreate, time.Local)
		if err != nil {
			mlog.Error(fmt.Sprintf("解析時間失敗: %s", err.Error()))
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
			BetUnique:  "mt-lottery-" + v.No,
			GameTypeId: gameTypeInt,
			BetAt:      parsedTime,
			Bet:        float64(v.BetTotal),
			EffectBet:  float64(v.BetValid),
			WinLose:    float64(v.Winnings),
			BetInfo:    info,
		}
		records = append(records, record)
	}

	return service.WriteBetRecord2(tx, records, func() error {
		return tx.Commit()
	})
}

func (m *MTLotteryGamer) CheckPoint(agentId int, member service.MemberGameAccountInfo) (int, error) {
	return 1, nil
}

// 積分轉出
func (m *MTLotteryGamer) PointOut(mid string, agentId int, point int) (int, error) {
	return 1, nil
}
