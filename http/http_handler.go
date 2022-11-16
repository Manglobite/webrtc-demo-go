package http

import (
	"log"
	"net/http"
)

// ListenAndServe размещает веб-ресурсы для доступа через браузер, предпочтительно Chrome
func ListenAndServe() {
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", fs)

	log.Print("web server listen on :3333...")

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Printf("web serve fail: %s", err.Error())
	}
}
