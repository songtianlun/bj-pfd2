package handle

import (
	"bj-pfd2/pkg/log"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Log(next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		log.Infof("Started %s %s for %s", request.Method, request.URL.Path, request.RemoteAddr)
		next(writer, request, params)
		//log.Infof("Completed %s %s for %s", request.Method, request.URL.Path, request.RemoteAddr)
	}
}
