package handle

import (
    "bj-pfd2/model"
    "bj-pfd2/pkg/log"
    "bj-pfd2/pkg/web"
    "net/http"
)

func Home(writer http.ResponseWriter, request *http.Request) {
    query := request.URL.Query()
    refresh := false
    if query.Get("refresh") != "" {
        refresh = true
    }
    log.Infof("Report With Cache? %v", !refresh)
    token := model.GetToken(request)
    fullData := GetAllData(token, refresh)
    fullData.StatisticAll()
    if refresh {
        http.Redirect(writer, request, "/", http.StatusFound)
    }
    web.SetTokenToCookie(writer, token)
    web.GenerateHTML(writer, fullData, "layout", "private.navbar", "home")
}
