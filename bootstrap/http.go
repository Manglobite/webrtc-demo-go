package bootstrap

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/webrtc-demo-go/config"
)

// Rest Отправьте HTTP-запрос на открытую платформу и верните
// полезную нагрузку ответа открытой платформы на верхний уровень.
func Rest(method string, url string, body io.Reader) (res []byte, err error) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("create http request fail: %s", err.Error())

		return
	}

	ts := time.Now().UnixNano() / 1000000
	sign := calBusinessSign(ts)

	request.Header.Set("Accept", "*")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Access-Control-Allow-Origin", "*")
	request.Header.Set("Access-Control-Allow-Methods", "*")
	request.Header.Set("Access-Control-Allow-Headers", "*")
	request.Header.Set("mode", "no-cors")
	request.Header.Set("client_id", config.App.ClientID)
	request.Header.Set("access_token", config.App.AccessToken)
	request.Header.Set("sign", sign)
	request.Header.Set("t", strconv.FormatInt(ts, 10))

	response, err := client.Do(request)
	if err != nil {
		log.Printf("http request fail: %s", err.Error())

		return
	}
	defer response.Body.Close()

	res, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read http response fail", err.Error())

		return
	}

	return
}
