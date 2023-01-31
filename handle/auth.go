package handle

import (
	"bj-pfd2/model"
	"bj-pfd2/pkg/log"
	"bj-pfd2/pkg/web"
	"context"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := model.CheckToken(r)
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Login GET /login
// Show the Login page
func Login(writer http.ResponseWriter, request *http.Request) {
	web.GenerateHTML(writer, nil, "layout", "public.navbar", "login")
}

// Authenticate the user given the email and password,  POST /authenticate
func Authenticate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		log.Error("Cannot parse form")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	token := request.PostFormValue("token")
	if token == "" {
		log.Error("Token is empty")
		web.ResponseWithUnauthorized(writer, "Token is empty")
	}
	if !TokenValid(token) {
		web.ResponseWithUnauthorized(writer, "无效的 Notion Token，请检查后重试。")
	}
	web.SetTokenToCookie(writer, token)
	writer.WriteHeader(http.StatusOK)
	//http.Redirect(writer, request, "/", 302)
	return
}

// Logout GET /logout
// Logs the user out
func Logout(writer http.ResponseWriter, request *http.Request) {
	_, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		log.Warn("Failed to get cookie")
	}
	http.SetCookie(writer, &http.Cookie{
		Name:     "_cookie",
		Value:    "",
		HttpOnly: true,
	})
	http.Redirect(writer, request, "/", 302)
}
