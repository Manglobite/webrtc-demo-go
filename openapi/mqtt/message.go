package openmqtt

// MqttFrameHeader заголовок фрейма сообщения mqtt
type MqttFrameHeader struct {
	// тип сообщения mqtt, предлагает отключение ответа-кандидата
	Type string `json:"type"`

	// отправитель сообщения mqtt
	From string `json:"from"`

	// получатель сообщений mqtt
	To string `json:"to"`

	// Если отправитель или получатель является устройством и дополнительным устройством, здесь указывается идентификатор дополнительного устройства.
	SubDevID string `json:"sub_dev_id"`

	// Идентификатор сеанса, которому принадлежит сообщение mqtt
	SessionID string `json:"sessionid"`

	// Идентификатор службы сигнализации moto, связанной с сообщением mqtt.
	MotoID string `json:"moto_id"`

	// Идентификатор транзакции, передаваемый при прозрачной передаче управляющих сигналов MQTT
	TransactionID string `json:"tid"`
}

// MqttFrame кадр сообщения mqtt
type MqttFrame struct {
	Header  MqttFrameHeader `json:"header"`
	Message interface{}     `json:"msg"` //В теле сообщения mqtt может быть предложено разъединение кандидата на ответ, так что это interface{}
}

// MqttMessage сообщение mqtt (включая заголовок протокола верхнего уровня)
type MqttMessage struct {
	Protocol int       `json:"protocol"` // Номер протокола сообщения mqtt, webRTC принадлежит службе потоковой передачи в реальном времени, равен 302.
	Pv       string    `json:"pv"`       // Номер версии протокола связи
	T        int64     `json:"t"`        // Временная метка Unix, единица измерения — секунда
	Data     MqttFrame `json:"data"`
}
