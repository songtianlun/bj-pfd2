package handle

import (
	"bj-pfd2/com/web"
	"bj-pfd2/model"
	"net/http"
)

func Err(writer http.ResponseWriter, request *http.Request) {
	val := request.URL.Query()
	msg := val.Get("msg")
	_, err := model.CheckSession(request)
	if err != nil {
		web.GenerateHTML(writer, msg, "layout", "public.navbar", "error")
	} else {
		web.GenerateHTML(writer, msg, "layout", "private.navbar", "error")
	}
}

func Index(writer http.ResponseWriter, request *http.Request) {
	web.GenerateHTML(writer, nil, "layout", "public.navbar", "index")
}
