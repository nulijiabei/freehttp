package freehttp

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// Handler
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// Server
type Server struct {
	m map[string]Handler
}

// New
func NewServer() *Server {
	server := new(Server)
	server.m = make(map[string]Handler)
	return server
}

// 添加服务
func (this *Server) Default(handler interface{}) *Service {
	service := NewService()
	if err := service.Register(handler); err != nil {
		panic(err)
	}
	this.m[service.name] = service
	return service
}

// 添加服务
func (this *Server) Append(handler interface{}) {
	this.m[reflect.Indirect(reflect.ValueOf(handler)).Type().Name()] = handler.(Handler)
}

// 启动服务
func (this *Server) Start(port string) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// ListenAndServe(addr string, handler Handler)
	return http.ListenAndServe(port, this)
}

// 包含 ServeHTTP(ResponseWriter, *Request) 函数符合 Handler 接口
func (this *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ok := false
	url := strings.Split(Trim(r.URL.Path), "/")
	for k, v := range this.m {
		if strings.ToLower(Trim(url[1])) == strings.ToLower(k) {
			ok = true
			v.ServeHTTP(w, r)
		}
	}
	if !ok {
		http.NotFound(w, r)
	}
}
