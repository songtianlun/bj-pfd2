package handle

import (
	"bj-pfd2/com/web"
	"net/http"
)

func Err(writer http.ResponseWriter, request *http.Request) {
	val := request.URL.Query()
	msg := val.Get("msg")
	web.GenerateHTML(writer, msg, "layout", "public.navbar", "error")
}

func Index(writer http.ResponseWriter, request *http.Request) {
	web.GenerateHTML(writer, nil, "layout", "private.navbar", "index")
}
