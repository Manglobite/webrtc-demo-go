package openmqtt

import (
	"encoding/json"
	"log"
	"time"

	"github.com/webrtc-demo-go/config"
	"github.com/webrtc-demo-go/types"
)

func Post(msg *types.WsMessage) {
	switch msg.Type {
	case "offer":
		sendOffer(msg, msg.Payload)
	case "candidate":
		sendCandidate(msg, msg.Payload)
	case "disconnect":
		sendDisconnect(msg)
	default:
		log.Printf("unsupported ws message, type: %s", msg.Type)
	}
}

func sendOffer(msg *types.WsMessage, sdp string) {
	offerFrame := struct {
		Mode       string `json:"mode"`        // режим предложения, по умолчанию — webrtc
		Sdp        string `json:"sdp"`         // Предложение, сгенерированное браузером
		StreamType uint32 `json:"stream_type"` // тип потока кода, по умолчанию 1
		Auth       string `json:"auth"`        // Код авторизации, необходимый для аутентификации через webRTC, получен с открытой платформы
	}{
		Mode:       "webrtc",
		Sdp:        sdp,
		StreamType: 1,
		Auth:       auth,
	}

	offerMqtt := &MqttMessage{
		Protocol: 302,
		Pv:       "2.2",
		T:        time.Now().Unix(),
		Data: MqttFrame{
			Header: MqttFrameHeader{
				Type:      "offer",
				From:      config.App.MQTTUID,
				To:        config.App.DeviceID,
				SubDevID:  "",
				SessionID: msg.SessionID,
				MotoID:    motoID,
			},
			Message: offerFrame,
		},
	}

	sendBytes, err := json.Marshal(offerMqtt)
	if err != nil {
		log.Printf("marshal offer mqtt to bytes fail: %s", err.Error())

		return
	}

	publish(sendBytes)
}

func sendCandidate(msg *types.WsMessage, candidate string) {
	candidateFrame := struct {
		Mode      string `json:"mode"`      // Кандидатский режим, по умолчанию webrtc
		Candidate string `json:"candidate"` // адрес кандидата，a=candidate:1922393870 1 UDP 2130706431 192.168.1.171 51532 typ host
	}{
		Mode:      "webrtc",
		Candidate: candidate,
	}

	candidateMqtt := &MqttMessage{
		Protocol: 302,
		Pv:       "2.2",
		T:        time.Now().Unix(),
		Data: MqttFrame{
			Header: MqttFrameHeader{
				Type:      "candidate",
				From:      config.App.MQTTUID,
				To:        config.App.DeviceID,
				SubDevID:  "",
				SessionID: msg.SessionID,
				MotoID:    motoID,
			},
			Message: candidateFrame,
		},
	}

	sendBytes, err := json.Marshal(candidateMqtt)
	if err != nil {
		log.Printf("marshal candidate mqtt to bytes fail: %s", err.Error())

		return
	}

	publish(sendBytes)
}

func sendDisconnect(msg *types.WsMessage) {
	disconnectFrame := struct {
		Mode string `json:"mode"`
	}{
		Mode: "webrtc", // Режим отключения, по умолчанию — webrtc
	}

	disconnectMqtt := &MqttMessage{
		Protocol: 302,
		Pv:       "2.2",
		T:        time.Now().Unix(),
		Data: MqttFrame{
			Header: MqttFrameHeader{
				Type:      "disconnect",
				From:      config.App.MQTTUID,
				To:        config.App.DeviceID,
				SubDevID:  "",
				SessionID: msg.SessionID,
				MotoID:    motoID,
			},
			Message: disconnectFrame,
		},
	}

	sendBytes, err := json.Marshal(disconnectMqtt)
	if err != nil {
		log.Printf("marshal candidate mqtt to bytes fail: %s", err.Error())

		return
	}

	publish(sendBytes)
}

// Опубликовать сообщение mqtt
func publish(payload []byte) {
	token := client.Publish(publishTopic, 1, false, payload)
	if token.Error() != nil {
		log.Printf("mqtt publish fail: %s, topic: %s", token.Error().Error(),
			publishTopic)
	}
}
