package game

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"game_service/instance/dg_live"
	"game_service/instance/gb_elec"
	"game_service/instance/mt_live"
	"game_service/instance/mt_lottery"
	"game_service/instance/rsg_elec"
	"game_service/instance/sa_live"
	"game_service/instance/wg_sport"
	"sync"

	"game_service/internal/cachedata"
	"game_service/internal/database"
	redisDriver "game_service/internal/redis"
	"game_service/internal/service"
	"os"
	"strconv"

	"game_service/internal/ws/gameclient"

	"time"

	"github.com/gomodule/redigo/redis"
	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

// 踢出遊戲
func (m *MyMemberGameAccountInfo) KickOutUser(db *sql.DB, agentId int) error {
	mlog.Info(fmt.Sprintf("[KickOutUser] 開始踢出遊戲 會員ID:%d 代理ID:%d", m.MemberId, agentId))
	conn := redisDriver.GetRedisConn()
	defer conn.Close()
	// 踢出目前遊戲頁面
	if err := KickOutAll(m.MemberId); err != nil && err != redis.ErrNil {
		return err
	}
	mlog.Info(fmt.Sprintf("[KickOutUser] 成功踢出遊戲 會員ID:%d 代理ID:%d", m.MemberId, agentId))
	// gameclient.GameWebSocket.RemoveClient(m.MemberId) // KickOutAll 已經移除了 不用再移除
	return nil
}

// 檢查玩家是否在遊戲中
func CheckMemberInGame(conn redis.Conn, mid int) (int, error) {
	redisKey := fmt.Sprintf("member:%d:game", mid)
	agentId, err := redis.Int(conn.Do("GET", redisKey))
	if err != nil {
		return 0, err
	}
	return agentId, nil
}

// 遊戲頁面踢除
func KickOutPageForRecyle(mid int, conn redis.Conn) error {
	_, exists := gameclient.GameWebSocket.Clients[mid]
	agentId, err := CheckMemberInGame(conn, mid)
	if err != nil {
		return err
	}
	if !exists && agentId == 0 {
		return nil
	}
	// 如果已有連接，發送關閉命令給舊連接
	db, err := database.MEMBER.DB()
	if err != nil {
		return err
	}
	member, err := service.GetMember(db, agentId, mid)
	if err != nil {
		return err
	}
	thisMember := MyMemberGameAccountInfo{member}
	thisMember.KickOutUser(db, agentId)

	return nil
}

func MemberAgentKickOut(member service.MemberGameAccountInfo, agentId int) error {
	member.GameAgentId = agentId
	agentNameIdMap := cachedata.AgentNameIdMap()

	var err error
	switch agentId {
	case 0:
		return nil
	case agentNameIdMap["mt_live"]:
		_, err = mt_live.Kickout(member)
	case agentNameIdMap["mt_lottery"]:
		err = mt_lottery.KickOut(member)
	case agentNameIdMap["gb_elec"]:
		_, err = gb_elec.KickOut(member)
	case agentNameIdMap["dg_live"]:
		err = dg_live.KickOut(member)
	case agentNameIdMap["rsg_elec"]:
		err = rsg_elec.KickOut(member)
	case agentNameIdMap["sa_live"]:
		err = sa_live.KickOut(member)
	case agentNameIdMap["wg_sport"]:
		err = wg_sport.KickOut(member)
	default:
		mlog.Info("無此代理:")
	}
	if err != nil && err != redis.ErrNil {
		mlog.Info(fmt.Sprintf("會員踢出失敗: %s", err.Error()))
		return errors.New("踢出失敗")
	}

	return nil
}

// 確認遊戲帳號是否存在
func (m *MyMemberGameAccountInfo) CheckGameUser() error {
	agentNameIdMap := cachedata.AgentNameIdMap()
	switch m.MemberGameAccountInfo.GameAgentId {
	case 0:
		return errors.New("遊戲商不存在")
	case agentNameIdMap["mt_live"]:
		// 確認遊戲商帳號是否存在
		if gameExistRes, err := mt_live.GameAccountExist(m.MemberGameAccountInfo.UserName); err != nil || gameExistRes.Data.Result == 0 {
			// 其他類型錯誤
			if gameExistRes == nil {
				return err
			}
			// 帳號不存在
			if err := mt_live.AccountRegister(m.MemberGameAccountInfo.MemberId, &m.MemberGameAccountInfo); err != nil {
				return err
			}
		}
	case agentNameIdMap["mt_lottery"]:
		if err := mt_lottery.AccountRegister(&m.MemberGameAccountInfo); err != nil {
			return err
		}
	case agentNameIdMap["gb_elec"]:
		if err := gb_elec.AccountRegister(&m.MemberGameAccountInfo, m.MemberGameAccountInfo.GameAgentId); err != nil {
			return err
		}
	case agentNameIdMap["dg_live"]:
		if err := dg_live.AccountRegister(&m.MemberGameAccountInfo); err != nil {
			return err
		}
	case agentNameIdMap["rsg_elec"]:
		if err := rsg_elec.AccountRegister(&m.MemberGameAccountInfo); err != nil {
			return err
		}
	case agentNameIdMap["sa_live"]:
		// 確認遊戲商帳號是否存在
		if gameExistRes, err := sa_live.GameAccountExist(m.MemberGameAccountInfo.UserName); err != nil || gameExistRes.ErrorMsgId == 0 {
			// 其他類型錯誤
			if gameExistRes == nil {
				return err
			}
			if err := sa_live.AccountRegister(&m.MemberGameAccountInfo); err != nil {
				return err
			}
		}
	case agentNameIdMap["wg_sport"]:
		if err := wg_sport.AccountRegister(&m.MemberGameAccountInfo); err != nil {
			return err
		}
	default:
		return errors.New("遊戲商不存在")
	}
	return nil
}

// 進入遊戲前準備 - 點數轉入轉出
func PointPrepare(mid int, agentId int, balanceAPI func(int) (float64, float64, error), inAPI func(agentId int) (float64, error)) error {
	mlog.Info(fmt.Sprintf("[PointPrepare] 開始執行點數準備 會員ID:%d 代理ID:%d", mid, agentId))

	// 點數回收
	if err := RecyclePoint(mid, balanceAPI); err != nil {
		mlog.Error(fmt.Sprintf("[PointPrepare] 點數回收失敗: %s", err.Error()))
		return err
	}

	// 點數吸入
	amount, err := inAPI(agentId)
	if err != nil {
		mlog.Error(fmt.Sprintf("[PointPrepare] 點數吸入失敗 會員ID:%d 代理ID:%d 錯誤:%s", mid, agentId, err.Error()))
		return err
	}

	mlog.Info(fmt.Sprintf("[PointPrepare] 準備轉入點數 會員ID:%d 代理ID:%d 金額:%f", mid, agentId, amount))
	if err := SendTransfer(TransferReq{
		MemberId: mid,
		Type:     "main-to-sub",
		AgentId:  agentId,
		Amount:   amount,
	}); err != nil {
		mlog.Error(fmt.Sprintf("[PointPrepare] 點數轉入失敗: %s", err.Error()))
		return err
	}

	mlog.Info(fmt.Sprintf("[PointPrepare] 點數準備完成 會員ID:%d 代理ID:%d", mid, agentId))
	return nil
}

// 進入遊戲前準備 - 點數轉入轉出
func RecyclePoint(mid int, balanceAPI func(int) (float64, float64, error)) error {
	mlog.Info(fmt.Sprintf("[RecyclePoint] 開始執行點數回收 會員ID:%d", mid))

	var wg sync.WaitGroup
	agentIdMap := cachedata.AgentIdMap()
	bringOutMap := make(map[int]GameWalletMap)

	errChan := make(chan error, len(agentIdMap))

	for _, agentId := range agentIdMap {
		wg.Add(1)
		go func(agentIdStr string) {
			defer wg.Done()
			id, err := strconv.Atoi(agentIdStr)
			if err != nil {
				errChan <- fmt.Errorf("[RecyclePoint] 代理ID轉換失敗: %s", err.Error())
				return
			}

			mlog.Info(fmt.Sprintf("[RecyclePoint] 檢查代理錢包餘額 會員ID:%d 代理ID:%d", mid, id))
			balance, transBalance, err := balanceAPI(id)
			if err != nil {
				errChan <- err
				mlog.Error(fmt.Sprintf("[RecyclePoint] 獲取遊戲錢包餘額失敗 會員ID:%d 代理ID:%d 錯誤:%s", mid, id, err.Error()))
				return
			}

			mlog.Info(fmt.Sprintf("[RecyclePoint] 準備更新錢包餘額 會員ID:%d 代理ID:%d 餘額:%f 轉帳餘額:%f", mid, id, balance, transBalance))
			if err := SendTransfer(TransferReq{
				MemberId: mid,
				Amount:   balance - transBalance,
				Type:     "game-to-sub",
				AgentId:  id,
			}); err != nil {
				errChan <- err
				mlog.Error(fmt.Sprintf("[RecyclePoint] 更新遊戲錢包失敗 會員ID:%d 代理ID:%d 錯誤:%s", mid, id, err.Error()))
				return
			}

			if balance >= 1 {
				bringOutMap[id] = GameWalletMap{
					Balance:      balance,
					TransBalance: transBalance,
				}
				mlog.Info(fmt.Sprintf("[RecyclePoint] 新增提領記錄 會員ID:%d 代理ID:%d 餘額:%f", mid, id, balance))
			}
		}(agentId)
	}

	// 等待所有goroutine完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误发生
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	mlog.Info(fmt.Sprintf("[RecyclePoint] 點數提領記錄: %v", bringOutMap))

	// 點數吸出
	if err := SendTransfer(TransferReq{
		MemberId:      mid,
		Type:          "bring-out",
		GameWalletMap: bringOutMap,
	}); err != nil {
		mlog.Error(fmt.Sprintf("[RecyclePoint] 點數吸出失敗 會員ID:%d 錯誤:%s", mid, err.Error()))
		return err
	}

	mlog.Info(fmt.Sprintf("[RecyclePoint] 點數回收完成 會員ID:%d", mid))
	return nil
}

// 打 wallet service 轉帳
func SendTransfer(body TransferReq) error {
	mlog.Info(fmt.Sprintf("[SendTransfer] 開始執行轉帳 會員ID:%d 類型:%s 金額:%v", body.MemberId, body.Type, body.Amount))

	prefix := "/v1/wallet/transfer"
	domain := os.Getenv("WALLET_SERVER")
	url := domain + prefix
	mlog.Info(fmt.Sprintf("[SendTransfer] 請求URL: %s", url))

	reqBody, err := json.Marshal(body)
	if err != nil {
		mlog.Error(fmt.Sprintf("[SendTransfer] 請求資料序列化失敗: %s", err.Error()))
		return err
	}
	mlog.Info(fmt.Sprintf("[SendTransfer] 請求內容: %s", string(reqBody)))

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(reqBody)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := client.Do(req, resp); err != nil {
		mlog.Error(fmt.Sprintf("[SendTransfer] HTTP請求失敗 會員ID:%d 錯誤:%s", body.MemberId, err.Error()))
		return err
	}

	// 記錄響應內容
	respBody := resp.Body()
	mlog.Info(fmt.Sprintf("[SendTransfer] 響應狀態碼:%d 響應內容:%s", resp.StatusCode(), string(respBody)))

	if resp.StatusCode() != fasthttp.StatusOK {
		mlog.Error(fmt.Sprintf("[SendTransfer] 轉帳失敗 會員ID:%d 狀態碼:%d 響應內容:%s",
			body.MemberId, resp.StatusCode(), string(respBody)))
		return errors.New("轉帳失敗")
	}

	mlog.Info(fmt.Sprintf("[SendTransfer] 轉帳完成 會員ID:%d 類型:%s 金額:%v",
		body.MemberId, body.Type, body.Amount))
	return nil
}

// KickOutAll 踢出所有遊戲
func KickOutAll(mid int) error {
	mlog.Info(fmt.Sprintf("[KickOutAll] Start kicking out member: %d", mid))

	conn := redisDriver.GetRedisConn()
	defer conn.Close()

	_, exists := gameclient.GameWebSocket.Clients[mid]
	agentId, err := CheckMemberInGame(conn, mid)
	mlog.Info(fmt.Sprintf("[KickOutAll] WebSocket exists: %v, AgentId: %d", exists, agentId))

	if err != nil && err != redis.ErrNil {
		mlog.Error(fmt.Sprintf("[KickOutAll] Error checking member in game: %s", err.Error()))
		return err
	}

	if agentId == 0 {
		mlog.Info(fmt.Sprintf("[KickOutAll] Member %d not in game, skipping", mid))
		return nil
	}

	db, err := database.MEMBER.DB()
	if err != nil {
		mlog.Error(fmt.Sprintf("[KickOutAll] Failed to get DB connection: %s", err.Error()))
		return err
	}

	mlog.Info(fmt.Sprintf("[KickOutAll] Getting member info for mid: %d, agentId: %d", mid, agentId))
	member, err := service.GetMember(db, agentId, mid)
	if err != nil {
		if err == sql.ErrNoRows {
			mlog.Error(fmt.Sprintf("[KickOutAll] Member %d not found", mid))
			return errors.New("會員帳號不存在")
		} else {
			mlog.Error(fmt.Sprintf("[KickOutAll] Failed to get member data: %s", err.Error()))
			return err
		}
	}

	thisMember := MyMemberGameAccountInfo{member}
	thisMember.MemberGameAccountInfo.GameAgentId = agentId

	if thisMember.MemberGameAccountInfo.GamePassword == "" {
		mlog.Info(fmt.Sprintf("[KickOutAll] Game password empty for member %d, checking game user", mid))
		if err := thisMember.CheckGameUser(); err != nil {
			mlog.Error(fmt.Sprintf("[KickOutAll] Failed to check game user: %s", err.Error()))
			return err
		}
	}

	mlog.Info(fmt.Sprintf("[KickOutAll] Kicking out member %d from agent %d", mid, agentId))
	if err := MemberAgentKickOut(thisMember.MemberGameAccountInfo, agentId); err != nil && err != redis.ErrNil {
		mlog.Error(fmt.Sprintf("[KickOutAll] Failed to kick out member from agent: %s", err.Error()))
		return err
	}

	// 刪除redis中的玩家遊戲資訊
	mlog.Info(fmt.Sprintf("[KickOutAll] Cleaning up Redis and WebSocket for member %d", mid))
	conn.Do("DEL", fmt.Sprintf("member:%d:game", mid))
	gameclient.GameWebSocket.RemoveClient(mid)

	mlog.Info(fmt.Sprintf("[KickOutAll] Successfully kicked out member %d", mid))
	return nil
}
