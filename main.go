package main

import (
	"WebSocket-I/src/infraestructure/mqtt"
	"WebSocket-I/src/infraestructure/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(engine)

	go mqtt.StartMQTTClient()
	err := engine.Run(":8082")
	if err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
