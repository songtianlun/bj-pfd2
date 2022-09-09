package model

import (
	"bj-pfd2/com/log"
	"fmt"
	"net/http"
)

func CheckToken(request *http.Request) (token string, err error) {
	cookie, err := request.Cookie("_cookie")
	if err != nil {
		log.ErrorF("Check Token fails - %v", err.Error())
		return
	}
	token = cookie.Value
	if token == "" {
		err = fmt.Errorf("token is empty")
	}
	log.InfoF("Check Token Pass")
	return
}

func GetToken(request *http.Request) (token string) {
	// 忽略错误，该方法理想情况下在 CheckToken 之后的方法中调用
	// 若出现错误则返回空字符串
	token, _ = CheckToken(request)
	return
}
