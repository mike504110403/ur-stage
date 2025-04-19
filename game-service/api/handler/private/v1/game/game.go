package game

import (
	"database/sql"
	"errors"
	"fmt"
	"game_service/instance/dg_live"
	"game_service/instance/gb_elec"
	"game_service/instance/mt_live"
	"game_service/instance/mt_lottery"
	"game_service/instance/rsg_elec"
	"game_service/instance/sa_live"
	"game_service/instance/wg_sport"
	"game_service/internal/cachedata"
	"game_service/internal/database"
	redisDriver "game_service/internal/redis"
	"game_service/internal/service"
	"math"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	mlog "github.com/mike504110403/goutils/log"
)

// Game 遊戲接口
// 定義遊戲基本操作
type Gamer interface {
	JoinGame(service.MemberGameAccountInfo) (string, error)
	//PointIn(string, int, int) (int, error)
	//PointOut(string, int, int) (int, error)
	CheckPoint(int, service.MemberGameAccountInfo) (int, error)
	LeaveGame() error
	GetRecord(string, string, time.Time, time.Time) ([]any, error)
}

// JointGame 進入遊戲
func JoinGame(mid int, agentId int) (string, error) {
	// 準備加入遊戲
	member, err := prepareJoin(mid, agentId)
	if err != nil {
		mlog.Error(fmt.Sprintf("加入遊戲前準備失敗: %s", err.Error()))
		return "", err
	}
	conn := redisDriver.GetRedisConn()
	if conn == nil {
		return "", fmt.Errorf("无法获取 Redis 连接")
	}
	defer conn.Close()

	// 添加重试机制
	var retryCount int
	for retryCount < 3 {
		if conn != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
		conn = redisDriver.GetRedisConn()
		retryCount++
	}

	if conn == nil {
		return "", fmt.Errorf("无法获取 Redis 连接，已重试 %d 次", retryCount)
	}

	gameNameIdMap := cachedata.AgentNameIdMap()
	if url, err := func(agentId int) (string, error) {
		switch agentId {
		case gameNameIdMap["mt_live"]:
			mt_live_instance := &mt_live.MTLiveGamer{}
			return mt_live_instance.JoinGame(member)
		case gameNameIdMap["mt_lottery"]:
			mt_lottery_instance := &mt_lottery.MTLotteryGamer{}
			return mt_lottery_instance.JoinGame(member)
		case gameNameIdMap["gb_elec"]:
			gb_elec_instance := &gb_elec.GBELECGamer{}
			return gb_elec_instance.JoinGame(member)
		case gameNameIdMap["dg_live"]:
			dg_live_instance := &dg_live.DGLiveGamer{}
			return dg_live_instance.JoinGame(member)
		case gameNameIdMap["rsg_elec"]:
			rsg_elec_instance := &rsg_elec.RSGElecGamer{}
			return rsg_elec_instance.JoinGame(member)
		case gameNameIdMap["sa_live"]:
			sa_live_instance := &sa_live.SALiveGamer{}
			return sa_live_instance.JoinGame(member)
		case gameNameIdMap["wg_sport"]:
			wg_sport_instance := &wg_sport.WGSportGamer{}
			return wg_sport_instance.JoinGame(member)
		default:
			return "", errors.New("無此代理")
		}
	}(agentId); err != nil {
		mlog.Error(fmt.Sprintf("加入遊戲失敗: %s", err.Error()))
		return "", err
	} else {
		// 更新 Redis，記錄玩家當前的遊戲
		redisKey := fmt.Sprintf("member:%d:game", mid)
		_, err = conn.Do("SET", redisKey, agentId)
		if err != nil {
			mlog.Info(fmt.Sprintf("更新 Redis 失敗: %s", err.Error()))
			return "", err
		}
		return url, nil
	}
}

// 加入遊戲前準備
func prepareJoin(mid int, agentId int) (service.MemberGameAccountInfo, error) {
	member := service.MemberGameAccountInfo{}

	// 記錄開始準備加入遊戲
	mlog.Info(fmt.Sprintf("[PrepareJoin] 開始準備加入遊戲流程 - 會員ID: %d, 遊戲商ID: %d", mid, agentId))

	db, err := database.MEMBER.DB()
	if err != nil {
		mlog.Error(fmt.Sprintf("[PrepareJoin] 資料庫連線失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
		return member, err
	}

	member, err = service.GetMember(db, agentId, mid)
	if err != nil {
		if err == sql.ErrNoRows {
			mlog.Error(fmt.Sprintf("[PrepareJoin] 會員帳號不存在 - 會員ID: %d", mid))
			return member, errors.New("會員帳號不存在")
		} else {
			mlog.Error(fmt.Sprintf("[PrepareJoin] 獲取會員資料失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
			return member, err
		}
	}
	mlog.Info(fmt.Sprintf("[PrepareJoin] 成功獲取會員資料 - 會員ID: %d, 遊戲帳號: %s", mid, member.UserName))

	thisMember := MyMemberGameAccountInfo{member}
	thisMember.MemberGameAccountInfo.GameAgentId = agentId

	// 檢查遊戲密碼
	if thisMember.MemberGameAccountInfo.GamePassword == "" {
		mlog.Info(fmt.Sprintf("[PrepareJoin] 遊戲密碼為空，進行遊戲帳號檢查 - 會員ID: %d", mid))
		if err := thisMember.CheckGameUser(); err != nil {
			mlog.Error(fmt.Sprintf("[PrepareJoin] 遊戲帳號檢查失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
			return member, err
		}
		mlog.Info(fmt.Sprintf("[PrepareJoin] 遊戲帳號檢查成功 - 會員ID: %d", mid))
	}

	// 踢出其他遊戲
	mlog.Info(fmt.Sprintf("[PrepareJoin] 開始踢出玩家其他遊戲 - 會員ID: %d", mid))
	if err = thisMember.KickOutUser(db, agentId); err != nil && err != redis.ErrNil {
		mlog.Error(fmt.Sprintf("[PrepareJoin] 踢出玩家失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
		return member, err
	}
	mlog.Info(fmt.Sprintf("[PrepareJoin] 成功踢出玩家其他遊戲 - 會員ID: %d", mid))

	// 確認點數
	mlog.Info(fmt.Sprintf("[PrepareJoin] 開始確認點數 - 會員ID: %d", mid))
	if err := CheckPoint(agentId, thisMember.MemberGameAccountInfo.MemberId); err != nil {
		mlog.Error(fmt.Sprintf("[PrepareJoin] 確認點數失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
		return member, err
	}
	mlog.Info(fmt.Sprintf("[PrepareJoin] 點數確���成功 - 會員ID: %d", mid))

	mlog.Info(fmt.Sprintf("[PrepareJoin] 完成所有準備工作 - 會員ID: %d", mid))
	return thisMember.MemberGameAccountInfo, nil
}

type Config struct {
	Frequency time.Duration
}

// GetRecord Cron
func BetRecordCron() {
	// 獲取 agentMap
	agentMap := cachedata.AgentNameIdMap()

	// quit channel 用來結束 cron
	quit := make(chan struct{})

	var tickerMtLive, tickerMtLottery, tickerGb, tickerDg, tickerRsg, tickerSa, tickerWGSP *time.Ticker

	if _, ok := agentMap["mt_live"]; ok {
		tickerMtLive = time.NewTicker(3 * time.Minute)
		defer tickerMtLive.Stop()
	}

	if _, ok := agentMap["mt_lottery"]; ok {
		tickerMtLottery = time.NewTicker(3 * time.Minute)
		defer tickerMtLottery.Stop()
	}

	if _, ok := agentMap["gb_elec"]; ok {
		tickerGb = time.NewTicker(5 * time.Minute)
		defer tickerGb.Stop()
	}

	if _, ok := agentMap["dg_live"]; ok {
		tickerDg = time.NewTicker(3 * time.Minute)
		defer tickerDg.Stop()
	}

	if _, ok := agentMap["rsg_elec"]; ok {
		tickerRsg = time.NewTicker(3 * time.Minute)
		defer tickerRsg.Stop()
	}
	if _, ok := agentMap["sa_live"]; ok {
		tickerSa = time.NewTicker(3 * time.Minute)
		defer tickerSa.Stop()
	}
	if _, ok := agentMap["wg_sport"]; ok {
		tickerWGSP = time.NewTicker(3 * time.Minute)
		defer tickerWGSP.Stop()
	}
	// 開始執行
	for {
		select {
		case <-quit:
			return

		case <-getTickerChannel(tickerMtLive):
			go mt_live.GetLiveBetRecord(3 * time.Minute)
			go mt_live.GetDonateRecord(3 * time.Minute)

		case <-getTickerChannel(tickerMtLottery):
			go mt_lottery.GetBetOrder(3 * time.Minute)

		case <-getTickerChannel(tickerGb):
			go gb_elec.GetBetRecord(5 * time.Minute)

		case <-getTickerChannel(tickerDg):
			go dg_live.GetBetRecord(3 * time.Minute)

		case <-getTickerChannel(tickerRsg):
			go rsg_elec.GetBetRecord(3 * time.Minute)

		case <-getTickerChannel(tickerSa):
			go sa_live.GetLiveBetRecord(3 * time.Minute)

		case <-getTickerChannel(tickerWGSP):
			go wg_sport.GetBetRecord(3 * time.Minute)
		default:
			// 添加一個休眠時間以避免忙等待
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// getTickerChannel 是一個幫助函數，確保在 Ticker 不為 nil 的情況下返回對應的 Channel
func getTickerChannel(ticker *time.Ticker) <-chan time.Time {
	if ticker != nil {
		return ticker.C
	}
	return nil
}

// PointIn 積分轉入
func PointIn(mid string, agentId int, point float64) (int, error) {
	switch agentId {
	case MT_LIVE_Agent.Id:
		mt_live_instance := &mt_live.MTLiveGamer{}
		return mt_live_instance.PointIn(mid, agentId, point)
	default:
		return 0, nil
	}
}

// CheckPoint 點數確認
func CheckPoint(agentId int, mid int) error {
	mlog.Info(fmt.Sprintf("[CheckPoint] ���始點數確認流程 - 會員ID: %d, 遊戲商ID: %d", mid, agentId))

	var balanceAPI func(int) (float64, float64, error)
	var inAPI func(int) (float64, error)

	// 連接資料庫
	mlog.Info(fmt.Sprintf("[CheckPoint] 嘗試連接資料庫 - 會員ID: %d", mid))
	db, err := database.MEMBER.DB()
	if err != nil {
		mlog.Error(fmt.Sprintf("[CheckPoint] 資料庫連線失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
		return err
	}
	mlog.Info(fmt.Sprintf("[CheckPoint] 資料庫連接成功 - 會員ID: %d", mid))

	// 準備 SQL statement
	mlog.Info(fmt.Sprintf("[CheckPoint] 準備 SQL Statement - 會員ID: %d", mid))
	getMember, err := service.PrepareMember(db)
	if err != nil {
		mlog.Error(fmt.Sprintf("[CheckPoint] SQL Statement 準備失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
		return err
	}
	mlog.Info(fmt.Sprintf("[CheckPoint] SQL Statement 準備成功 - 會員ID: %d", mid))

	// 準備 API 函數
	mlog.Info(fmt.Sprintf("[CheckPoint] 準備遊戲商 API - 會員ID: %d", mid))
	balanceAPI = SwitchBalanceAPI(mid, getMember)
	inAPI = SwitchInApi(mid, getMember)
	mlog.Info(fmt.Sprintf("[CheckPoint] 遊戲商 API 準備完成 - 會員ID: %d", mid))

	// 執行點數準備
	mlog.Info(fmt.Sprintf("[CheckPoint] 開始執行點數準備 - 會員ID: %d", mid))
	if err := PointPrepare(mid, agentId, balanceAPI, inAPI); err != nil {
		mlog.Error(fmt.Sprintf("[CheckPoint] 點數準備失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
		return err
	}
	mlog.Info(fmt.Sprintf("[CheckPoint] 點數準備完成 - 會員ID: %d", mid))

	mlog.Info(fmt.Sprintf("[CheckPoint] 點數確認流程完成 - 會員ID: %d", mid))
	return nil
}

// Recycle 回收點數
func Recycle(mid int) error {
	mlog.Info(fmt.Sprintf("[Recycle] 開始回收點數流程 - 會員ID: %d", mid))

	conn := redisDriver.GetRedisConn()
	defer conn.Close()

	if err := KickOutPageForRecyle(mid, conn); err != nil && err != redis.ErrNil {
		mlog.Error(fmt.Sprintf("[Recycle] 踢出玩家失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
		return err
	}
	mlog.Info(fmt.Sprintf("[Recycle] 踢出玩家成功 - 會員ID: %d", mid))

	if _, err := CheckMemberInGame(conn, mid); err != nil && err != redis.ErrNil {
		mlog.Error(fmt.Sprintf("[Recycle] 檢查玩家遊戲狀態失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
		return err
	}
	mlog.Info(fmt.Sprintf("[Recycle] 檢查玩家遊戲狀態成功 - 會員ID: %d", mid))

	var balanceAPI func(int) (float64, float64, error)

	db, err := database.MEMBER.DB()
	if err != nil {
		mlog.Error(fmt.Sprintf("[Recycle] 資料庫連線失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
		return err
	}
	mlog.Info(fmt.Sprintf("[Recycle] 資料庫連線成功 - 會員ID: %d", mid))

	getMember, err := service.PrepareMember(db)
	if err != nil {
		mlog.Error(fmt.Sprintf("[Recycle] 準備會員資料失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
		return err
	}
	mlog.Info(fmt.Sprintf("[Recycle] 準備會員資料成功 - 會員ID: %d", mid))

	// Prepare 遊戲商 API
	balanceAPI = SwitchBalanceAPI(mid, getMember)
	mlog.Info(fmt.Sprintf("[Recycle] 準備遊戲商 API 成功 - 會員ID: %d", mid))

	if err := RecyclePoint(mid, balanceAPI); err != nil {
		mlog.Error(fmt.Sprintf("[Recycle] 回收點數失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
		return err
	}
	mlog.Info(fmt.Sprintf("[Recycle] 回收點數成功 - 會員ID: %d", mid))

	return nil
}

// SwitchBalanceAPI：切換取得餘額API
func SwitchBalanceAPI(mid int, stmt *sql.Stmt) func(int) (float64, float64, error) {
	mlog.Info(fmt.Sprintf("[SwitchBalanceAPI] 開始準備餘額檢查 API - 會員ID: %d", mid))
	agentNameMap := cachedata.AgentNameIdMap()

	return func(agentId int) (float64, float64, error) {
		mlog.Info(fmt.Sprintf("[SwitchBalanceAPI] 開始檢查餘額 - 會員ID: %d, 遊戲商ID: %d", mid, agentId))

		switch agentId {
		case 0:
			mlog.Error(fmt.Sprintf("[SwitchBalanceAPI] 無效的遊戲商ID - 會員ID: %d", mid))
			return 0, 0, errors.New("無此代理")

		case agentNameMap["mt_live"]:
			mlog.Info(fmt.Sprintf("[SwitchBalanceAPI] 開始處理 MT Live 餘額 - 會員ID: %d", mid))
			userData := service.MemberGameAccountInfo{}
			err := stmt.QueryRow(agentId, mid).Scan(&userData.UserName, &userData.NickName, &userData.GamePassword, &userData.GameAgentId)
			if err != nil {
				mlog.Error(fmt.Sprintf("[SwitchBalanceAPI] 查詢會員資料失敗 [MT Live] - 會員ID: %d, 錯誤: %s", mid, err.Error()))
				return 0, 0, err
			}
			if userData.GamePassword == "" {
				mlog.Error(fmt.Sprintf("[SwitchBalanceAPI] 遊戲商帳號不存在 [MT Live] - 會員ID: %d", mid))
				return 0, 0, errors.New("遊戲商帳號不存在")
			}

			res, err := mt_live.CheckBalance(userData.UserName)
			if err != nil {
				mlog.Error(fmt.Sprintf("[SwitchBalanceAPI] 確認點數失敗 [MT Live] - 會員ID: %d, 錯誤: %s", mid, err.Error()))
				return 0, 0, err
			}

			balance, err := strconv.ParseFloat(res.Data.Balance, 64)
			if err != nil {
				mlog.Error(fmt.Sprintf("[SwitchBalanceAPI] ���額轉換失敗 [MT Live] - 會員ID: %d, 錯誤: %s", mid, err.Error()))
				return 0, 0, err
			}

			if math.Trunc(balance) >= 1 {
				mlog.Info(fmt.Sprintf("[SwitchBalanceAPI] 開始轉出點數 [MT Live] - 會員ID: %d, 金額: %f", mid, math.Trunc(balance)))
				if _, err := mt_live.TransferOut(userData.UserName, math.Trunc(balance)); err != nil {
					mlog.Error(fmt.Sprintf("[SwitchBalanceAPI] 轉出失敗 [MT Live] - 會員ID: %d, 錯誤: %s", mid, err.Error()))
					return 0, 0, err
				}
			}

			mlog.Info(fmt.Sprintf("[SwitchBalanceAPI] 完成餘額檢查 [MT Live] - 會員ID: %d, 餘額: %f", mid, balance))
			return balance, math.Trunc(balance), nil

		case agentNameMap["mt_lottery"]:
			userData := service.MemberGameAccountInfo{}
			err := stmt.QueryRow(agentId, mid).Scan(&userData.UserName, &userData.NickName, &userData.GamePassword, &userData.GameAgentId)
			if err != nil {
				return 0, 0, err
			}
			if userData.GamePassword == "" {
				return 0, 0, errors.New("遊戲帳號不存在")
			}
			balance, err := mt_lottery.CheckPoint(userData)
			if err != nil {
				mlog.Info(fmt.Sprintf("確認點數失敗 [%d mt_lottery]: %s", mid, err.Error()))
				return 0, 0, err
			} else {
				if math.Trunc(balance) >= 1 {
					if _, err := mt_lottery.TransferPoint(userData, -math.Trunc(balance)); err != nil {
						mlog.Info(fmt.Sprintf("攜出失敗 [%d mt_lottery]: %s", mid, err.Error()))
						return 0, 0, err
					}
				}

				return balance, math.Trunc(balance), nil
			}
		case agentNameMap["gb_elec"]:
			userData := service.MemberGameAccountInfo{}
			err := stmt.QueryRow(agentId, mid).Scan(&userData.UserName, &userData.NickName, &userData.GamePassword, &userData.GameAgentId)
			if err != nil {
				return 0, 0, err
			}
			if userData.GamePassword == "" {
				return 0, 0, errors.New("遊戲商帳號不存在")
			}
			balance, err := gb_elec.CheckPoint(userData)
			if err != nil {
				mlog.Info(fmt.Sprintf("確認點數失敗 [%d gb_elec]: %s", mid, err.Error()))
				return 0, 0, err
			} else {
				if math.Trunc(balance) >= 1 {
					err := gb_elec.TransferOutPoint(userData, -math.Trunc(balance))
					if err != nil {
						mlog.Info(fmt.Sprintf("攜出失敗 [%d gb_elec]: %s", mid, err.Error()))
						return 0, 0, err
					}
				}
				mlog.Info(fmt.Sprintf("確認點數成功 [%d gb_elec]: %f", mid, balance))
				return balance, math.Trunc(balance), nil
			}
		case agentNameMap["dg_live"]:
			userData := service.MemberGameAccountInfo{}
			err := stmt.QueryRow(agentId, mid).Scan(
				&userData.UserName, &userData.NickName,
				&userData.GamePassword,
				&userData.GameAgentId,
			)
			if err != nil {
				return 0, 0, err
			}
			if userData.GamePassword == "" {
				return 0, 0, errors.New("遊戲商帳號不存在")
			}

			balance, err := dg_live.CheckBlance(userData)
			if err != nil {
				mlog.Info(fmt.Sprintf("確認點數失敗 [%d dg_live]: %s", mid, err.Error()))
				return 0, 0, err
			} else {
				if math.Trunc(balance) >= 1 {
					err := dg_live.TransferPoint(userData, -math.Trunc(balance))
					if err != nil {
						mlog.Error(fmt.Sprintf("攜出失敗 [%d dg_live]: %s", mid, err.Error()))
						return 0, 0, err
					}
				}

				return balance, math.Trunc(balance), nil
			}
		case agentNameMap["rsg_elec"]:
			userData := service.MemberGameAccountInfo{}
			err := stmt.QueryRow(agentId, mid).Scan(&userData.UserName, &userData.NickName, &userData.GamePassword, &userData.GameAgentId)
			if err != nil {
				return 0, 0, err
			}
			if userData.GamePassword == "" {
				return 0, 0, errors.New("遊戲商帳號不存在")
			}
			balance, err := rsg_elec.GetBalance(userData.UserName)
			if err != nil {
				mlog.Info(fmt.Sprintf("確認點數失敗 [%d rsg_elec]: %s", mid, err.Error()))
				return 0, 0, err
			} else {
				if math.Trunc(balance) >= 1 {
					if err := rsg_elec.KickOut(userData); err != nil {
						mlog.Info(fmt.Sprintf("踢出失敗 [%d rsg_elec]: %s", mid, err.Error()))
						return 0, 0, err
					}
					if _, err := rsg_elec.PointOut(userData.UserName, math.Trunc(balance)); err != nil {
						mlog.Info(fmt.Sprintf("攜出失敗 [%d rsg_elec]: %s", mid, err.Error()))
						return 0, 0, err
					}
				}
				return balance, math.Trunc(balance), nil
			}
		case agentNameMap["sa_live"]:
			userData := service.MemberGameAccountInfo{}
			err := stmt.QueryRow(agentId, mid).Scan(&userData.UserName, &userData.NickName, &userData.GamePassword, &userData.GameAgentId)
			if err != nil {
				return 0, 0, err
			}
			if userData.GamePassword == "" {
				return 0, 0, errors.New("遊戲商帳號不存在")
			}
			balance, err := sa_live.GetBalance(userData)
			if err != nil {
				mlog.Info(fmt.Sprintf("確認點數失敗 [%d sa_live]: %s", mid, err.Error()))
				return 0, 0, err
			} else {
				if math.Trunc(balance) >= 1 {
					_, err := sa_live.DebitBalanceDVAPI(userData.UserName, math.Trunc(balance))
					if err != nil {
						mlog.Info(fmt.Sprintf("攜出失敗 [%d sa_live]: %s", mid, err.Error()))
						return 0, 0, err
					}
				}

				return balance, math.Trunc(balance), nil
			}
		case agentNameMap["wg_sport"]:
			userData := service.MemberGameAccountInfo{}
			err := stmt.QueryRow(agentId, mid).Scan(&userData.UserName, &userData.NickName, &userData.GamePassword, &userData.GameAgentId)
			if err != nil {
				return 0, 0, err
			}
			if userData.GamePassword == "" {
				return 0, 0, errors.New("遊戲商帳號不存在")
			}
			balance, err := wg_sport.PointUserAPI(userData.UserName)
			if err != nil {
				mlog.Info(fmt.Sprintf("確認點數失敗 [%d wg_sport]: %s", mid, err.Error()))
				return 0, 0, err
			} else {
				if math.Trunc(balance) >= 1 {
					_, err := wg_sport.TransferOut(userData, math.Trunc(balance))
					if err != nil {
						mlog.Info(fmt.Sprintf("攜出失敗 [%d wg_sport]: %s", mid, err.Error()))
						return 0, 0, err
					}
				}
				return balance, math.Trunc(balance), nil
			}
		default:
			mlog.Info(fmt.Sprintf("[SwitchBalanceAPI] 未知的遊戲商ID - 會員ID: %d, 遊戲商ID: %d", mid, agentId))
			return 0, 0, nil
		}
	}
}

// SwitchInApi：切換存提API
func SwitchInApi(mid int, stmt *sql.Stmt) func(int) (float64, error) {
	agentNameMap := cachedata.AgentNameIdMap()
	return func(agentId int) (float64, error) {
		mlog.Info(fmt.Sprintf("[SwitchInApi] 開始轉入點數流程 - 會員ID: %d, 遊戲商ID: %d", mid, agentId))

		// 做遊戲轉帳，從Ｗallet裡面獲取
		amount, err := service.GetBalance(mid)
		if err != nil {
			mlog.Error(fmt.Sprintf("[SwitchInApi] 獲取錢包餘額失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
			return 0, err
		}
		mlog.Info(fmt.Sprintf("[SwitchInApi] 獲取錢包餘額成功 - 會員ID: %d, 餘額: %f", mid, amount))

		transAmount := math.Trunc(amount)
		if transAmount <= 1 {
			mlog.Info(fmt.Sprintf("[SwitchInApi] 轉入金額過小，不執行轉入 - 會員ID: %d, 金額: %f", mid, transAmount))
			return 0, nil
		}

		switch agentId {
		case 0:
			mlog.Info(fmt.Sprintf("[SwitchInApi] 遊戲商ID為0，不執行轉入 - 會員ID: %d", mid))
			return 0, nil
		case agentNameMap["mt_live"]:
			mlog.Info(fmt.Sprintf("[SwitchInApi] 開始處理 MT Live 轉入 - 會員ID: %d", mid))
			userData := service.MemberGameAccountInfo{}
			if err := stmt.QueryRow(agentId, mid).Scan(
				&userData.UserName,
				&userData.NickName,
				&userData.GamePassword,
				&userData.GameAgentId,
			); err != nil {
				mlog.Error(fmt.Sprintf("[SwitchInApi] MT Live 查詢用戶資料失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
				return 0, err
			}
			mlog.Info(fmt.Sprintf("[SwitchInApi] MT Live 開始轉入點數 - 會員ID: %d, 用戶名: %s, 金額: %f", mid, userData.UserName, transAmount))
			if _, err = mt_live.TransferIn(userData.UserName, transAmount); err != nil {
				mlog.Error(fmt.Sprintf("[SwitchInApi] MT Live 轉入失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
				return 0, err
			}
		case agentNameMap["mt_lottery"]:
			userData := service.MemberGameAccountInfo{}
			if err := stmt.QueryRow(agentId, mid).Scan(
				&userData.UserName,
				&userData.NickName,
				&userData.GamePassword,
				&userData.GameAgentId,
			); err != nil {
				return 0, err
			}
			if _, err = mt_lottery.TransferPoint(userData, transAmount); err != nil {
				return 0, err
			}
		case agentNameMap["gb_elec"]:
			userData := service.MemberGameAccountInfo{}
			if err := stmt.QueryRow(agentId, mid).Scan(
				&userData.UserName,
				&userData.NickName,
				&userData.GamePassword,
				&userData.GameAgentId,
			); err != nil {
				return 0, err
			}
			if err = gb_elec.TransferPoint(userData, transAmount); err != nil {
				return 0, err
			}
		case agentNameMap["dg_live"]:
			userData := service.MemberGameAccountInfo{}
			err := stmt.QueryRow(agentId, mid).Scan(&userData.UserName, &userData.NickName, &userData.GamePassword, &userData.GameAgentId)
			if err != nil {
				return 0, err
			}
			err = dg_live.TransferPoint(userData, transAmount)
			if err != nil {
				mlog.Error(err.Error())
				return 0, err
			}
		case agentNameMap["rsg_elec"]:
			userData := service.MemberGameAccountInfo{}
			err := stmt.QueryRow(agentId, mid).Scan(&userData.UserName, &userData.NickName, &userData.GamePassword, &userData.GameAgentId)
			if err != nil {
				return 0, err
			}
			if _, err = rsg_elec.PointIn(userData.UserName, transAmount); err != nil {
				return 0, err
			}
		case agentNameMap["sa_live"]:
			mlog.Info(fmt.Sprintf("[SwitchInApi] 開始處理 SA Live 轉入 - 會員ID: %d", mid))
			userData := service.MemberGameAccountInfo{}
			if err := stmt.QueryRow(agentId, mid).Scan(&userData.UserName, &userData.NickName, &userData.GamePassword, &userData.GameAgentId); err != nil {
				mlog.Error(fmt.Sprintf("[SwitchInApi] SA Live 查詢用戶資料失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
				return 0, err
			}
			mlog.Info(fmt.Sprintf("[SwitchInApi] SA Live 開始轉入點數 - 會員ID: %d, 用戶名: %s, 金額: %f", mid, userData.UserName, transAmount))
			if _, err = sa_live.CreditBalanceDVAPI(userData.UserName, transAmount); err != nil {
				mlog.Error(fmt.Sprintf("[SwitchInApi] SA Live 轉入失敗 - 會員ID: %d, 錯誤: %s", mid, err.Error()))
				return 0, err
			}
		case agentNameMap["wg_sport"]:
			userData := service.MemberGameAccountInfo{}
			err := stmt.QueryRow(agentId, mid).Scan(&userData.UserName, &userData.NickName, &userData.GamePassword, &userData.GameAgentId)
			if err != nil {
				return 0, err
			}
			balance, err := wg_sport.PointUserAPI(userData.UserName)
			if err != nil {
				mlog.Info(fmt.Sprintf("確認點數失敗 [%d wg_sport]: %s", mid, err.Error()))
				return 0, err
			}
			if _, err = wg_sport.TransferIn(userData, balance, transAmount); err != nil {
				return 0, err
			}
		default:
			mlog.Error(fmt.Sprintf("[SwitchInApi] 未知的遊戲商ID - 會員ID: %d, 遊戲商ID: %d", mid, agentId))
			return 0, nil
		}

		mlog.Info(fmt.Sprintf("[SwitchInApi] 轉入點數成功 - 會員ID: %d, 遊戲商ID: %d, 金額: %f", mid, agentId, transAmount))
		return transAmount, nil
	}
}
