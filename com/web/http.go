package web

import (
	"bj-pfd2/com/cfg"
	"bj-pfd2/com/log"
	"embed"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"time"
)

type Middleware func(httprouter.Handle) httprouter.Handle

type Chain struct {
	middlewares []func(handler httprouter.Handle) httprouter.Handle
}

//var mux *http.ServeMux
var router *httprouter.Router
var gEfs *embed.FS

// Init initializes the web server
// 导入时自动实例化
func init() {
	//mux = http.NewServeMux()
	router = httprouter.New()
}

func RegisterHandle(method string, path string, handle httprouter.Handle, m ...func(handlerFunc httprouter.Handle) httprouter.Handle) {
	c := Chain{}
	c.middlewares = append(c.middlewares, m...)
	//mux.HandleFunc(path, c.Then(handle))
	//(*router).GET(path, handle)
	if method == "post" {
		(*router).POST(path, c.Then(handle))
	} else {
		(*router).GET(path, c.Then(handle))
	}
}

func RegisterDir(path string, file string, strip bool) {
	router.ServeFiles(path, http.Dir(file))
}

func RegisterEmbedFs(path string, efs *embed.FS, strip bool) {
	router.ServeFiles(path, http.FS(efs))
}

func RegisterTplEmbedFs(efs *embed.FS) {
	if gEfs == nil {
		gEfs = efs
	}
}

func (c Chain) Then(next httprouter.Handle) httprouter.Handle {
	for i := range c.middlewares {
		prev := c.middlewares[len(c.middlewares)-1-i]
		next = prev(next)
	}
	return next
}

func Run(address string) {
	server := &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.GetInt64("ReadTimeout") * int64(time.Second)),
		WriteTimeout:   time.Duration(cfg.GetInt64("WriteTimeout") * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.ErrorF("web server error: %s", err.Error())
		return
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
		log.ErrorF("Generate HTML error: %v", err.Error())
	}
}

func ParseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.tmpl", file))
	}
	t = template.Must(t.ParseFS(gEfs, files...))
	return
}
