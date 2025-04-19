package gameclient

import (
	"fmt"
	"game_service/api/handler/private/v1/game"
	"game_service/internal/locals"
	redisDriver "game_service/internal/redis"
	gws "game_service/internal/ws/gameclient"

	"github.com/gomodule/redigo/redis"
	mlog "github.com/mike504110403/goutils/log"

	"sync"

	"github.com/gofiber/websocket/v2"
)

// 遊戲頁面
func PageHandler(c *websocket.Conn) {
	mlog.Info("Starting PageHandler function")
	localUser, err := locals.GetUserInfo(c)
	if err != nil {
		mlog.Error(fmt.Sprintf("Error getting user info: %v", err))
		c.Close()
		return
	}
	mid := localUser.MemberId
	mlog.Info(fmt.Sprintf("WebSocket connection started for user %d with details: %+v", mid, localUser))

	// 創建一個done channel用於控制goroutine退出
	done := make(chan struct{})
	errChan := make(chan error, 1) // 使用buffered channel

	var once sync.Once // 確保只關閉一次
	cleanup := func() {
		once.Do(func() {
			mlog.Info(fmt.Sprintf("Starting cleanup process for user %d with connection details: %+v", mid, c.RemoteAddr()))
			close(done)
			handleDisconnect(mid)
			c.Close()
			mlog.Info(fmt.Sprintf("Cleanup completed for user %d, all resources released", mid))
		})
	}
	defer cleanup()

	if err := gws.GameWebSocket.AddClient(mid, c, done); err != nil {
		mlog.Error(fmt.Sprintf("Failed to add client for user %d: %v, connection details: %+v", mid, err, c.RemoteAddr()))
		return
	}
	mlog.Info(fmt.Sprintf("Successfully added client for user %d to GameWebSocket manager", mid))

	// 啟動讀取消息的goroutine
	go func() {
		mlog.Info(fmt.Sprintf("Starting message reading goroutine for user %d, connection: %v", mid, c.RemoteAddr()))
		for {
			select {
			case <-done:
				mlog.Info(fmt.Sprintf("Received done signal for user %d, stopping message reading. Connection: %v", mid, c.RemoteAddr()))
				return
			default:
				messageType, message, err := c.ReadMessage()
				if err != nil {
					mlog.Error(fmt.Sprintf("Error reading message for user %d: %v, connection: %v", mid, err, c.RemoteAddr()))
					errChan <- err
					return
				}
				mlog.Info(fmt.Sprintf("Received message from user %d: %s, message type: %d", mid, string(message), messageType))

				switch gws.WSMessage(message) {
				case gws.WSMessageGameClose:
					mlog.Info(fmt.Sprintf("Processing game close request for user %d, connection: %v", mid, c.RemoteAddr()))
					handleDisconnect(mid)
					if err := c.WriteMessage(messageType, []byte("leave game success")); err != nil {
						mlog.Error(fmt.Sprintf("Error sending leave game success message to user %d: %v, connection: %v", mid, err, c.RemoteAddr()))
					}
					errChan <- nil
					mlog.Info(fmt.Sprintf("Game close request completed for user %d, connection closed successfully", mid))
					return
				case gws.WSMessageGameNew:
					mlog.Info(fmt.Sprintf("Processing new game request for user %d", mid))
					if err := gws.GameWebSocket.AddClient(mid, c, done); err != nil {
						mlog.Error(fmt.Sprintf("Failed to add client for user %d: %v", mid, err))
						errChan <- err
						return
					}
					if err := c.WriteMessage(messageType, []byte("join game success")); err != nil {
						mlog.Error(fmt.Sprintf("Error sending join game success message to user %d: %v", mid, err))
					}
					mlog.Info(fmt.Sprintf("New game request completed for user %d", mid))
				case gws.WSMessageGameHeart:
					mlog.Info(fmt.Sprintf("Processing heartbeat for user %d", mid))
					if err := c.WriteMessage(messageType, []byte("pong")); err != nil {
						mlog.Error(fmt.Sprintf("Error sending pong message to user %d: %v", mid, err))
						errChan <- err
						return
					}
					mlog.Info(fmt.Sprintf("Heartbeat processed for user %d", mid))
				}
			}
		}
	}()

	// 等待goroutine完成或出現錯誤
	mlog.Info(fmt.Sprintf("Waiting for message handling completion for user %d", mid))
	if err := <-errChan; err != nil {
		mlog.Error(fmt.Sprintf("WebSocket error occurred for user %d: %v", mid, err))
	}
	mlog.Info(fmt.Sprintf("WebSocket handler completed for user %d", mid))
}

func handleDisconnect(mid int) {
	mlog.Info(fmt.Sprintf("Starting handleDisconnect for user %d", mid))
	conn := redisDriver.GetRedisConn()
	if conn == nil {
		mlog.Error(fmt.Sprintf("Failed to get Redis connection for user %d disconnect process", mid))
		return
	}
	defer conn.Close()

	mlog.Info(fmt.Sprintf("Starting cleanup process for user %d in handleDisconnect", mid))

	// 移除客戶端
	defer func() {
		if err := gws.GameWebSocket.RemoveClient(mid); err != nil {
			mlog.Error(fmt.Sprintf("Failed to remove client for mid: %d, error: %v, stack trace: %+v", mid, err, err))
		} else {
			mlog.Info(fmt.Sprintf("Successfully removed client for mid: %d from GameWebSocket", mid))
		}
	}()

	// 踢出玩家
	if err := game.KickOutAll(mid); err != nil && err != redis.ErrNil {
		mlog.Error(fmt.Sprintf("Kick out all failed for mid: %d, error: %v, details: %+v", mid, err, err))
	} else {
		mlog.Info(fmt.Sprintf("Successfully kicked out user %d from all games", mid))
	}

	// 清理 Redis
	redisKey := fmt.Sprintf("member:%d:game", mid)
	mlog.Info(fmt.Sprintf("Attempting to delete Redis key: %s for user %d", redisKey, mid))
	_, err := conn.Do("DEL", redisKey)
	if err != nil {
		mlog.Error(fmt.Sprintf("Redis DEL failed for mid: %d, key: %s, error: %v", mid, redisKey, err))
	} else {
		mlog.Info(fmt.Sprintf("Successfully deleted Redis key: %s for user %d", redisKey, mid))
	}

	mlog.Info(fmt.Sprintf("Completed handleDisconnect process for user: %d", mid))
}
