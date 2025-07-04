package routes

import (
	"WebSocket-I/src/infraestructure/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/ws", func(c *gin.Context) {
		userID := c.Query("user_id")
		sessionID := c.Query("session_id")

		if userID == "" || sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id y session_id son requeridos"})
			return
		}
		websocket.WSHandler(c.Writer, c.Request)
	})
}
