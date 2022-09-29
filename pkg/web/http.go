package web

import (
	"bj-pfd2/pkg/log"
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

var gEfs *embed.FS

func RegisterTplEmbedFs(efs *embed.FS) {
	if gEfs == nil {
		gEfs = efs
	}
}

func GenerateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.tmpl", file))
	}

	templates := template.Must(template.ParseFS(gEfs, files...))
	err := templates.ExecuteTemplate(writer, "layout", data)
	if err != nil {
		log.Errorf("Generate HTML error: %v", err.Error())
	}
}

func ResponseWithUnauthorized(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte(msg))
	if err != nil {
		return
	}
	return
}
