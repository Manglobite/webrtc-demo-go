package openmqtt

import (
	"log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/webrtc-demo-go/bootstrap"
	"github.com/webrtc-demo-go/config"
)

var (
	client mqtt.Client

	motoID string
	auth   string

	iceServers string

	publishTopic   string
	subscribeTopic string
)

func Start() (err error) {
	motoID, auth, iceServers, err = bootstrap.GetMotoIDAndAuth()
	if err != nil {
		log.Printf("allocate motoID fail: %s", err.Error())

		return
	}

	log.Printf("motoID: %s", motoID)
	log.Printf("auth: %s", auth)
	log.Printf("iceServers: %s", iceServers)

	hubConfig, err := bootstrap.LoadHubConfig()
	if err != nil {
		log.Printf("loadConfig fail: %s", err.Error())

		return
	}

	log.Printf("hubConfig: %+v", *hubConfig)

	publishTopic = hubConfig.SinkTopic.IPC
	subscribeTopic = hubConfig.SourceSink.IPC

	publishTopic = strings.Replace(publishTopic, "moto_id", motoID, 1)
	publishTopic = strings.Replace(publishTopic, "{device_id}", config.App.DeviceID, 1)

	log.Printf("publish topic: %s", publishTopic)
	log.Printf("subscribe topic: %s", subscribeTopic)

	// !!!При отправке сообщений mqtt from не является идентификатором
	// пользователя в webrtc.json, его необходимо обновить до идентификатора
	// в теме подписки, возвращаемой открытой платформой!!!
	parts := strings.Split(subscribeTopic, "/")
	config.App.MQTTUID = parts[3]

	opts := mqtt.NewClientOptions().AddBroker(hubConfig.Url).
		SetClientID(hubConfig.ClientID).
		SetUsername(hubConfig.Username).
		SetPassword(hubConfig.Password).
		SetOnConnectHandler(onConnect).
		SetConnectTimeout(10 * time.Second)

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("create mqtt client fail: %s", token.Error().Error())

		err = token.Error()
		return
	}

	return
}

func FetchWebRTCConfigs() (err error) {
	_, _, iceServers, err = bootstrap.GetMotoIDAndAuth()
	if err != nil {
		log.Printf("get webrtc configs fail: %s", err.Error())

		return err
	}

	log.Printf("iceServers: %s", iceServers)

	return nil
}

// IceServers Верните лед, возвращенный открытой платформой Token
func IceServers() string {
	return iceServers
}

// Функция обратного вызова успешного соединения mqtt, подписка на тему, возвращаемую открытой платформой, и получение сообщений mqtt
func onConnect(client mqtt.Client) {
	options := client.OptionsReader()

	log.Printf("%s connect to mqtt success", options.ClientID())

	if token := client.Subscribe(subscribeTopic, 1, consume); token.Wait() && token.Error() != nil {
		log.Printf("subcribe fail: %s, topic: %s", token.Error().Error(), subscribeTopic)

		return
	}

	log.Print("subscribe mqtt topic success")
}

func Unsubscribe() {
	if token := client.Unsubscribe(subscribeTopic); token.Wait() && token.Error() != nil {
		log.Printf("unsubscribe fail: %s, topic: %s", token.Error().Error(), subscribeTopic)
	}
}

func Disconnect() {
	client.Disconnect(1000)
}
