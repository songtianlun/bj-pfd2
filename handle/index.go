package handle

import (
	"bj-pfd2/com/web"
	"bj-pfd2/model"
	"net/http"
)

func Err(writer http.ResponseWriter, request *http.Request) {
	val := request.URL.Query()
	msg := val.Get("msg")
	web.GenerateHTML(writer, msg, "layout", "public.navbar", "error")
}

func Index(writer http.ResponseWriter, request *http.Request) {
	token := model.GetToken(request)
	indexData := model.Index{
		Token: token,
	}
	web.GenerateHTML(writer, indexData, "layout", "private.navbar", "index")
}
