// websocket/ws_handler.go
package websocket

import (
	"WebSocket-I/src/application"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar a WebSocket:", err)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	sessionIDStr := r.URL.Query().Get("session_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("user_id inválido")
		conn.Close()
		return
	}

	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		log.Println("session_id inválido")
		conn.Close()
		return
	}

	application.Manager.AddClient(userID, sessionID, conn)
	defer application.Manager.RemoveClient(userID, sessionID)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
