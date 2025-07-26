package mqtt

import (
	"WebSocket-I/src/application"
	"WebSocket-I/src/domain/entities"
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	payloadStr := string(msg.Payload())
	var payload entities.Message

	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		log.Printf("Error al decodificar el mensaje JSON: %v\n", err)
		return
	}

	log.Printf("Mensaje recibido en [%s]: %s\n", msg.Topic(), payloadStr)

	userID := payload.IdUser
	code := payload.Code

	if code == 0 {
		log.Printf("No code recibido para userID %d, no se enviará mensaje WS\n", userID)
		return
	}

	if payload.Alcohol != 0 {
		data := map[string]interface{}{
			"user_id":               userID,
			"alcohol_concentration": payload.Alcohol,
		}
		sendToAPI("https://fermest-api.it2id.cc/api/sensor-alcohol/create", data)
		sendToUser(userID, code, data)
	}

	if payload.Temperatura != 0 {
		data := map[string]interface{}{
			"user_id":     userID,
			"temperature": payload.Temperatura,
		}
		sendToAPI("https://fermest-api.it2id.cc/api/sensor-temperatura/create", data)
		sendToUser(userID, code, data)
	}

	if payload.Conductividad != 0 {
		data := map[string]interface{}{
			"user_id":      userID,
			"conductivity": payload.Conductividad,
		}
		sendToAPI("https://fermest-api.it2id.cc/api/sensor-conductividad/create", data)
		sendToUser(userID, code, data)
	}

	if payload.Turbuidez != 0 {
		data := map[string]interface{}{
			"user_id":   userID,
			"turbidity": payload.Turbuidez,
		}
		sendToAPI("https://fermest-api.it2id.cc/api/sensor-turbuidez/create", data)
		sendToUser(userID, code, data)
	}

	if payload.PH != 0 {
		data := map[string]interface{}{
			"user_id":  userID,
			"ph_value": payload.PH,
		}
		sendToAPI("https://fermest-api.it2id.cc/api/sensor-ph/create", data)
		sendToUser(userID, code, data)
	}

	if payload.Densidad != 0 {
		data := map[string]interface{}{
			"user_id": userID,
			"density": payload.Densidad,
		}
		sendToAPI("https://fermest-api.it2id.cc/api/sensor-densidad/create", data)
		sendToUser(userID, code, data)
	}
	if payload.Rpm != 0 {
		data := map[string]interface{}{
			"user_id": userID,
			"rpm":     payload.Rpm,
		}
		sendToAPI("https://fermest-api.it2id.cc/api/motor/create", data)
		sendToUser(userID, code, data)
	}
}

func sendToAPI(apiURL string, data map[string]interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error al convertir los datos a JSON: %v\n", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error al hacer la solicitud HTTP a %s: %v\n", apiURL, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Respuesta de la API [%s]: %s\n", apiURL, resp.Status)
}

func sendToUser(userID int, code int, data map[string]interface{}) {
	msgJSON, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error serializando dato para enviar WS: %v\n", err)
		return
	}
	application.Manager.SendTo(userID, code, msgJSON)
}

func StartMQTTClient() {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://52.202.0.30:1883").
		SetClientID("mi-consumidor").
		SetUsername("milton").
		SetPassword("milton123").
		SetDefaultPublishHandler(messageHandler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Error al conectar al broker MQTT:", token.Error())
	}

	topics := []string{"sensores"}
	for _, topic := range topics {
		if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
			log.Fatalf("Error al suscribirse al tópico %s: %v", topic, token.Error())
		}
		log.Println("Suscrito al tópico:", topic)
	}
}
