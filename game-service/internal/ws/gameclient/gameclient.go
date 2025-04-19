package gameclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

type WSMessage string

var GameWebSocket GameWebSocketManager

const (
	WSMessageGameClose WSMessage = "LeaveGame"
	WSMessageGameNew   WSMessage = "JoinGame"
	WSMessageGameHeart WSMessage = "ping"
)

type GameWebSocketManager struct {
	CloseChans   map[int]chan struct{}
	Clients      map[int]*websocket.Conn
	ClientsMutex *sync.RWMutex
	Ctx          context.Context
	cancel       context.CancelFunc
}

// 初始化 WebSocketManager
func NewWebSocketManager() {
	ctx, cancel := context.WithCancel(context.Background())
	GameWebSocket = GameWebSocketManager{
		CloseChans:   make(map[int]chan struct{}),
		Clients:      make(map[int]*websocket.Conn),
		ClientsMutex: &sync.RWMutex{},
		Ctx:          ctx,
		cancel:       cancel,
	}
}

// 添加客戶端連接
func (wsm *GameWebSocketManager) AddClient(userId int, conn *websocket.Conn, closeChan chan struct{}) error {
	if conn == nil {
		return fmt.Errorf("websocket connection cannot be nil")
	}
	if userId <= 0 {
		return fmt.Errorf("invalid user ID: %d", userId)
	}
	if closeChan == nil {
		closeChan = make(chan struct{}, 1)
	}

	wsm.ClientsMutex.Lock()
	defer wsm.ClientsMutex.Unlock()

	if oldConn, exists := wsm.Clients[userId]; exists {
		oldConn.Close()
		if oldChan := wsm.CloseChans[userId]; oldChan != nil {
			close(oldChan)
		}
	}

	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	wsm.Clients[userId] = conn
	wsm.CloseChans[userId] = closeChan
	return nil
}

// 刪除客戶端連接
func (wsm *GameWebSocketManager) RemoveClient(userId int) error {
	if userId <= 0 {
		return fmt.Errorf("invalid user ID: %d", userId)
	}

	// 首先获取连接并移除映射
	wsm.ClientsMutex.Lock()
	conn := wsm.Clients[userId]
	closeChan := wsm.CloseChans[userId]
	delete(wsm.Clients, userId)
	delete(wsm.CloseChans, userId)
	wsm.ClientsMutex.Unlock()

	if conn == nil {
		return nil
	}

	// 在释放锁后执行清理操作
	closeTimeout := time.Second * 5
	deadline := time.Now().Add(closeTimeout)
	if err := conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		deadline); err != nil {
		fmt.Printf("error sending close message: %v\n", err)
	}

	conn.Close()
	if closeChan != nil {
		close(closeChan)
	}

	return nil
}

// 發送消息給指定客戶端
func (wsm *GameWebSocketManager) SendMessage(userId int, message string, conn *websocket.Conn) error {
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}

	wsm.ClientsMutex.RLock()
	if storedConn, exists := wsm.Clients[userId]; !exists || storedConn != conn {
		wsm.ClientsMutex.RUnlock()
		return fmt.Errorf("invalid connection for user %d", userId)
	}
	wsm.ClientsMutex.RUnlock()

	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}
