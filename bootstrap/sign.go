package bootstrap

import (
	"crypto/md5"
	"fmt"

	"github.com/webrtc-demo-go/config"
)

// В соответствии с текущим временем (миллисекунды) сгенерировать подпись для запроса токена Restful открытой платформы.
func calTokenSign(ts int64) string {
	data := fmt.Sprintf("%s%s%d", config.App.ClientID, config.App.Secret, ts)

	val := md5.Sum([]byte(data))

	// md5值转换为大写
	res := fmt.Sprintf("%X", val)
	return res
}

// По текущему времени (миллисекунды) сгенерировать подпись для бизнес-запросов Restful на открытой платформе.
func calBusinessSign(ts int64) string {
	data := fmt.Sprintf("%s%s%s%d", config.App.ClientID, config.App.AccessToken, config.App.Secret, ts)

	val := md5.Sum([]byte(data))

	// значение md5 преобразовано в верхний регистр
	res := fmt.Sprintf("%X", val)
	return res
}
