package config

import (
	"encoding/json"
	"io/ioutil"
)

// Easy В простом режиме вам нужно вручную заполнить uId IPC, к которому вы хотите получить доступ
type Easy struct {
	UID string `json:"uId"`
}

// Auth В режиме кода авторизации для авторизации необходимо ввести
// пароль учетной записи пользователя на странице авторизации Tuya Open Platform,
// а возвращенный код авторизации заполнить здесь
type Auth struct {
	Code string `json:"code"`
}

// APPConfig Sample Запуск конфигурации приложения
type APPConfig struct {
	OpenAPIMode string `json:"openAPIMode"` // Режим подключения к открытой платформе Tuya пока поддерживает только mqtt
	OpenAPIURL  string `json:"openAPIUrl"`  // URL-адрес открытой платформы Tuya

	ClientID string `json:"clientId"` // AppId приложения открытой платформы Tuya
	Secret   string `json:"secret"`   // Секретный ключ приложения открытой платформы Tuya

	AuthMode string `json:"authMode"` // Режим авторизации, для "easy" / "auth"

	Easy Easy `json:"easy"`
	Auth Auth `json:"auth"`

	DeviceID string `json:"deviceId"` // Пример идентификатора устройства для подключения

	UID          string `json:"-"` // Идентификатор пользователя, которому принадлежит идентификатор устройства, заполненный вручную в простом режиме, access_token, полученный через код авторизации в режиме авторизации, вернет идентификатор пользователя
	MQTTUID      string `json:"-"` // Идентификатор темы на веб-странице при общении с Tuya MQTT
	AccessToken  string `json:"-"` // Access_token, возвращаемый режимом кода авторизации Tuya Open Platform
	RefreshToken string `json:"-"` // Refresh_token, возвращаемый режимом кода авторизации Tuya Open Platform.
	ExpireTime   int64  `json:"-"` // Срок действия токена, возвращаемого режимом кода авторизации Tuya Open Platform.
}

var App = APPConfig{
	OpenAPIMode: "mqtt",
	OpenAPIURL:  "openapi.tuyacn.com",
}

// LoadConfig нагрузка webrtc.json Настроен на среду выполнения
func LoadConfig() error {
	return parseJSON("webrtc.json", &App)
}

func parseJSON(path string, v interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	return err
}
