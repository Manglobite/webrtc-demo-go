package bootstrap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"github.com/webrtc-demo-go/config"
	"github.com/webrtc-demo-go/types"
)

// GetMotoIDAndAuth Получите идентификатор службы сигнализации moto и код авторизации,
// необходимый для аутентификации через webRTC, с открытой платформы в соответствии с идентификатором пользователя и идентификатором устройства.
func GetMotoIDAndAuth() (motoID, auth, iceServers string, err error) {
	url := fmt.Sprintf("https://%s/v1.0/users/%s/devices/%s/webrtc-configs", config.App.OpenAPIURL, config.App.UID, config.App.DeviceID)

	body, err := Rest("GET", url, nil)
	if err != nil {
		log.Printf("GET webrtc-configs fail: %s, body: %s", err.Error(), string((body)))

		return
	}

	motoIDValue := gjson.GetBytes(body, "result.moto_id")
	if !motoIDValue.Exists() {
		log.Printf("moto_id not exist in webrtc-configs, body: %s", string(body))

		return "", "", "", errors.New("moto_id not exist")
	}

	authValue := gjson.GetBytes(body, "result.auth")
	if !authValue.Exists() {
		log.Printf("auth not exist in webrtc-configs, body: %s", string(body))

		return "", "", "", errors.New("auth not exist")
	}

	iceServersValue := gjson.GetBytes(body, "result.p2p_config.ices")
	if !iceServersValue.Exists() {
		log.Printf("iceServers not exist in webrtc-configs, body: %s", string(body))

		return "", "", "", errors.New("p2p_config.ices not exist")
	}

	var tokens []types.Token
	err = json.Unmarshal([]byte(iceServersValue.String()), &tokens)
	if err != nil {
		log.Printf("unmarshal to tokens fail: %s", err.Error())
		return "", "", "", err
	}

	ices := make([]types.WebToken, 0)
	for _, token := range tokens {
		if strings.HasPrefix(token.Urls, "stun") {
			ices = append(ices, types.WebToken{
				Urls: token.Urls,
			})
		} else if strings.HasPrefix(token.Urls, "turn") {
			ices = append(ices, types.WebToken{
				Urls:       token.Urls,
				Username:   token.Username,
				Credential: token.Credential,
			})
		}
	}

	iceServersBytes, err := json.Marshal(&ices)
	if err != nil {
		log.Printf("marshal token to web tokens fail: %s", err.Error())
		return "", "", "", err
	}

	motoID = motoIDValue.String()
	auth = authValue.String()
	iceServers = string(iceServersBytes)

	return
}

// LoadHubConfig Получить информацию о соединении mqtt с открытой платформы
func LoadHubConfig() (config *types.OpenIoTHubConfig, err error) {
	body, err := getOpenIoTHubConfig()
	if err != nil {
		log.Printf("get OpenIoTHub config from http fail: %s", err.Error())

		return
	}

	if gjson.GetBytes(body, "success").Bool() != true {
		log.Printf("request OpenIoTHub Config fail, body: %s", string(body))

		return nil, errors.New("request hub config fail")
	}

	config = &types.OpenIoTHubConfig{}

	err = json.Unmarshal([]byte(gjson.GetBytes(body, "result").String()), config)
	if err != nil {
		log.Printf("unmarshal OpenIoTHub config to object fail: %s, body: %s", err.Error(), string(body))

		return
	}

	return
}

func getOpenIoTHubConfig() ([]byte, error) {
	url := fmt.Sprintf("https://%s/v2.0/open-iot-hub/access/config", config.App.OpenAPIURL)

	request := &types.OpenIoTHubConfigRequest{
		UID:      config.App.UID,
		UniqueID: uuid.New().String(),
		LinkType: "mqtt",
		Topics:   "ipc",
	}

	payload, err := json.Marshal(request)
	if err != nil {
		log.Printf("marshal OpenIoTHubConfig Request fail: %s", err.Error())

		return nil, err
	}

	return Rest("POST", url, bytes.NewReader(payload))
}
