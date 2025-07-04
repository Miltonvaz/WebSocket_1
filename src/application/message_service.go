package application

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientManager struct {
	mu      sync.Mutex
	clients map[int]map[int]*websocket.Conn
}

var Manager = ClientManager{
	clients: make(map[int]map[int]*websocket.Conn),
}

func (cm *ClientManager) AddClient(userID int, sessionID int, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.clients[userID] == nil {
		cm.clients[userID] = make(map[int]*websocket.Conn)
	}
	cm.clients[userID][sessionID] = conn
	log.Printf("Cliente agregado: userID=%d sessionID=%d\n", userID, sessionID)
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

func (cm *ClientManager) SendTo(userID int, sessionID int, message []byte) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if sessions, ok := cm.clients[userID]; ok {
		if conn, ok := sessions[sessionID]; ok {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error enviando a userID=%d sessionID=%d: %v\n", userID, sessionID, err)
				conn.Close()
				delete(sessions, sessionID)
			}
		}
	}
}
