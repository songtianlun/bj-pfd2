package handle

import (
	"bj-pfd2/pkg/web"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Err(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	val := request.URL.Query()
	msg := val.Get("msg")
	web.GenerateHTML(writer, msg, "layout", "public.navbar", "error")
}

func Index(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	query := request.URL.Query()
	refresh := false
	if query.Get("refresh") != "" {
		refresh = true
	}
	web.GenerateHTML(writer, refresh, "layout", "empty.navbar", "index")
}
