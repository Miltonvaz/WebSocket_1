package application

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientInfo struct {
	Conn *websocket.Conn
	Code int
}

type ClientManager struct {
	mu      sync.Mutex
	clients map[int]map[int]ClientInfo
}

var Manager = ClientManager{
	clients: make(map[int]map[int]ClientInfo),
}

func (cm *ClientManager) AddClient(userID int, sessionID int, conn *websocket.Conn, code int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.clients[userID] == nil {
		cm.clients[userID] = make(map[int]ClientInfo)
	}
	cm.clients[userID][sessionID] = ClientInfo{
		Conn: conn,
		Code: code,
	}
	log.Printf("Cliente agregado: userID=%d sessionID=%d code=%d\n", userID, sessionID, code)
}

func (cm *ClientManager) RemoveClient(userID int, sessionID int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if sessions, ok := cm.clients[userID]; ok {
		delete(sessions, sessionID)
		if len(sessions) == 0 {
			delete(cm.clients, userID)
		}
	}
	log.Printf("Cliente eliminado: userID=%d sessionID=%d\n", userID, sessionID)
}

func (cm *ClientManager) SendTo(userID int, code int, message []byte) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if sessions, ok := cm.clients[userID]; ok {
		for sessionID, client := range sessions {
			if client.Code == code {
				err := client.Conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error enviando a userID=%d sessionID=%d: %v\n", userID, sessionID, err)
					client.Conn.Close()
					delete(sessions, sessionID)
				}
			}
		}
	}
}
