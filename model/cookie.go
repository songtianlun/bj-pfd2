package model

import (
	"fmt"
	"net/http"
)

func CheckToken(request *http.Request) (token string, err error) {
	cookie, err := request.Cookie("_cookie")
	if err != nil {
		return "", err
	}
	token = cookie.Value
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}
	return
}
