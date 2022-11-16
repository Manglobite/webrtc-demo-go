package bootstrap

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tidwall/gjson"
	"github.com/webrtc-demo-go/config"
)

// InitToken Получить токен на основе кода авторизации
func InitToken() (err error) {
	var url string

	switch config.App.AuthMode {
	case "easy":
		url = fmt.Sprintf("https://%s/v1.0/token?grant_type=1", config.App.OpenAPIURL)
	case "auth":
		url = fmt.Sprintf("https://%s/v1.0/token?grant_type=2&code=%s", config.App.OpenAPIURL, config.App.Auth.Code)
	default:
		return fmt.Errorf("unsupported auth mode %s", config.App.AuthMode)
	}

	body, err := Rest("GET", url, nil)
	if err != nil {
		log.Printf("GET token fail: %s, body: %s", err.Error(), string((body)))

		return
	}

	err = syncToConfig(body)
	if err != nil {
		log.Printf("sync OpenAPI ressponse to config fail: %s", err.Error())

		return
	}

	// Запустите сопрограмму обновления обслуживания токенов
	go maintainToken()

	return
}

// Последующее использование refresh_token для обновления токена и получения нового refresh_token
func refreshToken() (err error) {
	url := fmt.Sprintf("https://%s/v1.0/token/%s", config.App.OpenAPIURL, config.App.RefreshToken)

	body, err := Rest("GET", url, nil)
	if err != nil {
		log.Printf("GET token fail: %s, body: %s", err.Error(), string((body)))

		return
	}

	err = syncToConfig(body)
	if err != nil {
		log.Printf("sync OpenAPI ressponse to config fail: %s", err.Error())

		return
	}

	return
}

// Синхронизируйте ответ интерфейса службы токенов OpenAPI с образцом глобальной конфигурации приложения.
func syncToConfig(body []byte) error {
	uIdValue := gjson.GetBytes(body, "result.uid")
	if !uIdValue.Exists() {
		log.Printf("uid not exits in body: %s", string(body))

		return errors.New("uid not exist")
	}

	accessTokenValue := gjson.GetBytes(body, "result.access_token")
	if !accessTokenValue.Exists() {
		log.Printf("access_token not exits in body: %s", string(body))

		return errors.New("access_token not exist")
	}

	refreshTokenValue := gjson.GetBytes(body, "result.refresh_token")
	if !refreshTokenValue.Exists() {
		log.Printf("refresh_token not exist")

		return errors.New("refresh_token not exist")
	}

	expireTimeValue := gjson.GetBytes(body, "result.expire_time")
	if !expireTimeValue.Exists() {
		log.Printf("expire_time not exist")

		return errors.New("expire_time not exist")
	}

	switch config.App.AuthMode {
	case "easy":
		config.App.UID = config.App.Easy.UID
	case "auth":
		config.App.UID = uIdValue.String()
	default:
		return fmt.Errorf("unsupported auth mode %s", config.App.AuthMode)
	}

	config.App.AccessToken = accessTokenValue.String()
	config.App.RefreshToken = refreshTokenValue.String()
	config.App.ExpireTime = expireTimeValue.Int()

	log.Printf("UID: %s", config.App.UID)
	log.Printf("AccessToken: %s", config.App.AccessToken)
	log.Printf("RefreshToken: %s", config.App.RefreshToken)
	log.Printf("ExpireTime: %d", config.App.ExpireTime)

	return nil
}

// После успешного получения токена в первый раз требуется регулярное обслуживание и обновление токена.
// Если обновление не удалось, обновляйте снова каждые 60 сек.
// Если обновление прошло успешно, оно будет обновлено за 300 секунд до истечения срока действия токена.
func maintainToken() {
	interval := config.App.ExpireTime - 300

	for {
		timer := time.NewTimer(time.Duration(interval) * time.Second)

		select {
		case <-timer.C:
			if err := refreshToken(); err != nil {
				log.Printf("refresh token fail: %s", err.Error())

				interval = 60
			} else {
				interval = config.App.ExpireTime - 300
			}
		}
	}
}
