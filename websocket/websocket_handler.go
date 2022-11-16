package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/webrtc-demo-go/bootstrap"
	"github.com/webrtc-demo-go/config"
	openmqtt "github.com/webrtc-demo-go/openapi/mqtt"
	"github.com/webrtc-demo-go/types"
)

// Передача WebSocket в этом образце — json, а тип сообщения — 1 (текст).
// В этом образце отключается проверка исходного адреса запроса, которая должна быть включена в производственной среде.
var upgrader = websocket.Upgrader{
	Subprotocols: []string{"json"},
	CheckOrigin:  checkOrigin,
}

// Отключить проверку исходного адреса запроса
func checkOrigin(r *http.Request) bool {
	return true
}

// ListenAndServe Предоставьте запись службы WebSocket/webrtc, обновите протокол HTTP до WebSocket
func ListenAndServe() {
	http.HandleFunc("/webrtc", webrtc)

	log.Print("websocket server listen on :5555...")

	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Printf("websocket serve fail: %s", err.Error())
	}
}

// Функция обработки соединения WebSocket,
// в Golang каждое соединение имеет
// свою собственную сопрограмму (аналогично потокам в C++/Java, более легковесная)
func webrtc(w http.ResponseWriter, r *http.Request) {
	// Обновите протокол подключения до WebSocket
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade to websocket fail: %s", err.Error())

		return
	}
	defer c.Close()

	log.Printf("new ws client, addr: %s", r.RemoteAddr)

	// Сохраните идентификатор прокси текущего клиента WebSocket.
	agentID := ""

	// Опрос сообщений из соединения WebSocket
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("ws read fail: %s", err.Error())

			break
		}

		log.Printf("ws recv: %s", string(message))

		msg := &types.WsMessage{
			Conn: c,
		}

		err = json.Unmarshal(message, msg)
		if err != nil {
			log.Printf("unmarshal ws message fail: %s", err.Error())

			break
		}

		// Установить идентификатор прокси
		agentID = msg.AgentID

		// Увеличьте отношение сопоставления
		// между идентификатором сеанса и соединением WebSocket.
		bootstrap.AddLink(msg.AgentID, msg.SessionID, msg)

		dispatch(msg)
	}

	if agentID != "" {
		bootstrap.RemoveLinkByConnLost(agentID)
	}
}

func sendIceServers(c *websocket.Conn) {
	iceServers := &types.WsMessage{
		Type:    "webrtcConfigs",
		Payload: openmqtt.IceServers(),
		Success: true,
	}

	sendBytes, err := json.Marshal(iceServers)
	if err != nil {
		log.Printf("marshal iceServers fail: %s", err.Error())

		return
	}

	// iceServers send back to Javascript
	err = c.WriteMessage(1, sendBytes)
	if err != nil {
		log.Printf("ws write fail: %s", err.Error())
	}
}

func dispatch(msg *types.WsMessage) {
	switch config.App.OpenAPIMode {
	case "mqtt":
		// Каждый раз, когда страница браузера нажимает «Позвонить», она вытягивает webrtc. configs
		if msg.Type == "webRTCConfigs" {
			if err := openmqtt.FetchWebRTCConfigs(); err != nil {
				log.Printf("%s fetch webrtc configs fail", msg.AgentID)
			} else {
				sendIceServers(msg.Conn)
			}
		} else {
			openmqtt.Post(msg)
		}
	default:
		log.Printf("OpenAPI webRTC only support [mqtt], mode: %s", config.App.OpenAPIMode)
	}
}
