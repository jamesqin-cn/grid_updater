package util

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HttpServer struct {
	host   string
	router *httprouter.Router
}

func NewHttpServer(host string) *HttpServer {
	router := httprouter.New()
	router.GET("/version", ShowVersion)

	return &HttpServer{
		host,
		router,
	}
}

func (s *HttpServer) GET(path string, handle httprouter.Handle) {
	s.router.GET(path, handle)
}

func (s *HttpServer) POST(path string, handle httprouter.Handle) {
	s.router.POST(path, handle)
}

func (s *HttpServer) RunHttpServer() {
	log.Fatal(http.ListenAndServe(s.host, s.router))
}

func ShowVersion(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, GetAppInfo())
}

func GetQuery(r *http.Request, key string, defaultVal string) string {
	values, ok := r.URL.Query()[key]
	if ok && len(values) > 0 && len(values[0]) > 0 {
		return values[0]
	}
	return defaultVal
}

func GetPost(r *http.Request, key string, defaultVal string) string {
	r.ParseMultipartForm(32 << 20)
	if r.MultipartForm != nil {
		values := r.MultipartForm.Value[key]
		if len(values) > 0 && len(values[0]) > 0 {
			return values[0]
		}
	}
	return defaultVal
}

func GetCookie(r *http.Request, key string, defaultVal string) string {
	cookie, err := r.Cookie(key)
	if err == nil && len(cookie.Value) > 0 {
		return cookie.Value
	}
	return defaultVal
}
