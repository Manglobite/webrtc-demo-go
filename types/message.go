package types

import "github.com/gorilla/websocket"

// WsMessage Структура верхнего уровня канала WebSocket для отправки и получения сообщений
type WsMessage struct {
	Conn *websocket.Conn `json:"-"` // Соединение WebSocket, соединение с веб-клиентом

	AgentID string `json:"agentId"` // Идентификатор прокси-сервера страницы браузера

	Type      string `json:"type"`              // тип сообщения, для offer、candidate、answer
	SessionID string `json:"sessionId"`         // идентификатор сессии
	Payload   string `json:"payload,omitempty"` // Содержимое сообщения, передаваемое WebSocket

	Success bool `json:"success"` // Отмечает успешность запроса WebSocket, действителен только при ответе веб-клиенту.
}

type OpenIoTHubConfig struct {
	Url      string `json:"url"`       // адрес подключения mqtt (включая протокол, ip, порт)
	ClientID string `json:"client_id"` // mqtt подключает client_id (уникальное и постоянное сопоставление, сгенерированное учетной записью пользователя и unique_id), clientId можно использовать для публикации или подписки.
	Username string `json:"username"`  // имя пользователя подключения mqtt (уникальное и постоянное сопоставление, созданное учетной записью пользователя)
	Password string `json:"password"`  // mqtt пароль для подключения, это поле остается неизменным в течение срока действия

	// Опубликуйте тему, управление устройством можно пройти через эту тему
	SinkTopic struct {
		IPC string `json:"ipc"`
	} `json:"sink_topic"`

	// Подпишитесь на тему, событие устройства, синхронизация состояния устройства, вы можете подписаться на эту тему
	SourceSink struct {
		IPC string `json:"ipc"`
	} `json:"source_topic"`

	ExpireTime int `json:"expire_time"` // Действительная продолжительность текущей конфигурации, все соединения будут отключены после того, как текущая конфигурация станет недействительной.
}

// OpenIoTHubConfigRequest Подать заявку на получение тела http-запроса mqtt-соединения с открытой платформой.
type OpenIoTHubConfigRequest struct {
	UID      string `json:"uid"`       // Идентификатор пользователя Туя
	UniqueID string `json:"unique_id"` // Конец соединения изолирован уникальным_идентификатором. Когда одному и тому же пользователю необходимо войти в систему на нескольких концах, вызывающая сторона должна убедиться, что уникальный_идентификатор отличается.
	LinkType string `json:"link_type"` // Тип подключения, на данный момент поддерживает только mqtt
	Topics   string `json:"topics"`    // mqtt беспокойства тема, этот пример посвящен только IPC topic
}

// Token ICE Token from OpenAPI
type Token struct {
	Urls       string `json:"urls"`
	Username   string `json:"username"`
	Credential string `json:"credential"`
	TTL        int    `json:"ttl"`
}

// WebToken ICE Token to Chrome
type WebToken struct {
	Urls       string `json:"urls,omitempty"`
	Username   string `json:"username,omitempty"`
	Credential string `json:"credential,omitempty"`
}
