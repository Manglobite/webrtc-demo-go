package main

import (
	"log"

	"github.com/webrtc-demo-go/bootstrap"
	"github.com/webrtc-demo-go/config"
	"github.com/webrtc-demo-go/http"
	openmqtt "github.com/webrtc-demo-go/openapi/mqtt"
	"github.com/webrtc-demo-go/websocket"

	"sync"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	var wg sync.WaitGroup

	wg.Add(1)

	startWebRTCSample()
	log.Print("start webRTC sample")

	wg.Wait()
}

func startWebRTCSample() {
	if err := config.LoadConfig(); err != nil {
		log.Printf("load webrtc.json to runtime fail: %s", err.Error())
		return
	}

	// Получите токен службы открытой платформы в соответствии
	// с кодом авторизации и регулярно обновляйте токен.
	if err := bootstrap.InitToken(); err != nil {
		log.Printf("init token fail: %s", err.Error())
		return
	}

	// Прежде чем mqtt получит доступ к открытой платформе,
	// вам необходимо получить соответствующую конфигурацию
	// через интерфейс Restful, чтобы запустить клиент mqtt.
	if config.App.OpenAPIMode == "mqtt" {
		if err := openmqtt.Start(); err != nil {
			log.Printf("start mqtt fail: %s", err.Error())

			return
		}
	}

	// Запустите веб-сервер
	go http.ListenAndServe()

	// Запустите веб-сервер
	go websocket.ListenAndServe()
}
