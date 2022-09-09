package handle

import (
	"bj-pfd2/com/web"
	"bj-pfd2/model"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Err(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	val := request.URL.Query()
	msg := val.Get("msg")
	web.GenerateHTML(writer, msg, "layout", "public.navbar", "error")
}

func Index(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	token := model.GetToken(request)
	fullData := GetAllData(token, false)
	fullData.StatisticAll()
	web.GenerateHTML(writer, fullData, "layout", "private.navbar", "index")
}
