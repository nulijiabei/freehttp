package freehttp

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// Server Json HTTP
// http://127.0.0.1:8080/MyStructName.MyFuncName
type Server struct {
	name    string
	rcvr    reflect.Value
	typ     reflect.Type
	methods map[string]reflect.Method
}

// 创建 Server
func NewServer() *Server {
	server := new(Server)
	server.methods = make(map[string]reflect.Method)
	return server
}

// 错误输出
func (this *Server) Error(err error) {
	if err != nil {
		fmt.Println("server exception:", err.Error())
	}
}

// 启动服务
func (this *Server) Start(port string) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// ListenAndServe(addr string, handler Handler)
	return http.ListenAndServe(port, this)
}

// 包含 ServeHTTP(ResponseWriter, *Request) 函数符合 Handler 接口
func (this *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := false
	request := NewRequest(r)
	responseWriter := NewResponseWriter(w)
	for name, method := range this.methods {
		path := strings.ToLower(strings.Split(r.URL.Path, "/")[len(strings.Split(r.URL.Path, "/"))-1])
		def := strings.ToLower(fmt.Sprintf("%s.%s", this.name, name))
		if def == path {
			status = true
			value := make([]reflect.Value, method.Type.NumIn())
			value[0] = this.rcvr
			for n := 1; n < method.Type.NumIn(); n++ {
				inType := method.Type.In(n).String()
				switch inType {
				case "*freehttp.Request":
					value[n] = reflect.ValueOf(request)
				case "*freehttp.ResponseWriter":
					value[n] = reflect.ValueOf(responseWriter)
				case "freehttp.Json":
					value[n] = reflect.ValueOf(request.ReadJson())
				case "freehttp.ContentType":
					value[n] = reflect.ValueOf(request.ReadContentType())
				case "freehttp.Stream":
					value[n] = reflect.ValueOf(request.ReadStream())
				default:
					this.Error(fmt.Errorf("unsupported in type: %s", inType))
				}
			}
			returnValues := method.Func.Call(value)
			for t := 0; t < method.Type.NumOut(); t++ {
				reType := method.Type.Out(t).String()
				content := returnValues[t].Interface()
				if content == nil {
					this.Error(fmt.Errorf("%s out value is null -> %s", name, reType))
					continue
				}
				switch reType {
				case "freehttp.HttpStatus":
					responseWriter.WriteHeader(content)
				case "freehttp.ContentType":
					responseWriter.SetContentType(content)
				case "freehttp.Json":
					this.Error(responseWriter.WriterJson(content))
				case "freehttp.JsonIndent":
					this.Error(responseWriter.WriterJsonIndent(content))
				case "freehttp.Stream":
					this.Error(responseWriter.WriterStream(content))
				case "error":
					this.Error(content.(error))
				default:
					this.Error(fmt.Errorf("unsupported out type: %s", reType))
				}
			}
		}
	}
	if !status {
		http.NotFound(w, r)
	}
}

// 将类及方法注册到FreeHttp
func (this *Server) Register(rcvr interface{}) error {
	this.typ = reflect.TypeOf(rcvr)
	this.rcvr = reflect.ValueOf(rcvr)
	this.name = reflect.Indirect(this.rcvr).Type().Name()
	if this.name == "" {
		return fmt.Errorf("no service name for type %s", this.typ.String())
	}
	for m := 0; m < this.typ.NumMethod(); m++ {
		method := this.typ.Method(m)
		// mtype := method.Type
		mname := method.Name
		this.methods[mname] = method
	}
	return nil
}
