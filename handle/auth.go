package handle

import (
	"bj-pfd2/com/log"
	"bj-pfd2/com/web"
	"bj-pfd2/model"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Auth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		_, err := model.CheckToken(r)
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
		next(w, r, p)
	}
}

// Login GET /login
// Show the Login page
func Login(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	t := web.ParseTemplateFiles("login.layout", "public.navbar", "login")
	err := t.Execute(writer, nil)
	if err != nil {
		log.Error("Cannot execute template: " + err.Error())
		return
	}
}

// Authenticate the user given the email and password,  POST /authenticate
func Authenticate(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	err := request.ParseForm()
	if err != nil {
		log.Error("Cannot parse form")
	}
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    request.PostFormValue("token"),
		HttpOnly: true,
	}
	http.SetCookie(writer, &cookie)
	http.Redirect(writer, request, "/", 302)
}

// Logout GET /logout
// Logs the user out
func Logout(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
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
