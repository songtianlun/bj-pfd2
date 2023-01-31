package handle

import (
	"bj-pfd2/pkg/web"
	"net/http"
)

func Err(writer http.ResponseWriter, request *http.Request) {
	val := request.URL.Query()
	msg := val.Get("msg")
	web.GenerateHTML(writer, msg, "layout", "public.navbar", "error")
}

func Index(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	refresh := false
	if query.Get("refresh") != "" {
		refresh = true
	}
	web.GenerateHTML(writer, refresh, "layout", "empty.navbar", "index")
}
