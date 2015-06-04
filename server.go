package freehttp

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// Server Json HTTP
// http://127.0.0.1:8080/MyStructName/MyFuncName
type Server struct {
	def     func(string, string) string
	name    string
	rcvr    reflect.Value
	typ     reflect.Type
	methods map[string]reflect.Method
}

// 创建 Server 其中 def 为路径处理函数 nil 则使用默认
func NewServer(def func(string, string) string) *Server {
	server := new(Server)
	server.methods = make(map[string]reflect.Method)
	if def != nil {
		server.def = def
	} else {
		server.def = _def
	}
	return server
}

// 路径定义
func _def(mname, name string) string {
	return strings.ToLower(fmt.Sprintf("%s/%s", mname, name))
}

// 错误输出
func (this *Server) error(err interface{}) {
	if err != nil {
		fmt.Println("server exception:", err.(error).Error())
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
	freehttp := NewFreeHttp(w, r)
	r.URL.Path = strings.ToLower(r.URL.Path)
	for name, method := range this.methods {
		if this.def(this.name, name) == r.URL.Path {
			status = true
			value := make([]reflect.Value, method.Type.NumIn())
			value[0] = this.rcvr
			for n := 1; n < method.Type.NumIn(); n++ {
				inType := method.Type.In(n).String()
				switch inType {
				case "*freehttp.FreeHttp":
					value[n] = reflect.ValueOf(freehttp)
				case "*freehttp.Request":
					value[n] = reflect.ValueOf(freehttp.SuperRequest)
				case "*freehttp.ResponseWriter":
					value[n] = reflect.ValueOf(freehttp.SuperResponseWriter)
				case "freehttp.Json":
					value[n] = reflect.ValueOf(freehttp.SuperRequest.ReadJson())
				case "freehttp.ContentType":
					value[n] = reflect.ValueOf(freehttp.SuperRequest.ReadContentType())
				case "freehttp.Stream":
					value[n] = reflect.ValueOf(freehttp.SuperRequest.ReadStream())
				default:
					this.error(fmt.Errorf("unsupported in type: %s", inType))
				}
			}
			returnValues := method.Func.Call(value)
			for t := 0; t < method.Type.NumOut(); t++ {
				reType := method.Type.Out(t).String()
				content := returnValues[t].Interface()
				if content == nil && reType != "error" {
					this.error(fmt.Errorf("%s out value is null -> %s", name, reType))
					continue
				}
				switch reType {
				case "freehttp.HttpStatus":
					freehttp.SuperResponseWriter.WriteHeader(content)
				case "freehttp.ContentType":
					freehttp.SuperResponseWriter.SetContentType(content)
				case "freehttp.Json":
					this.error(freehttp.SuperResponseWriter.WriterJson(content))
				case "freehttp.JsonIndent":
					this.error(freehttp.SuperResponseWriter.WriterJsonIndent(content))
				case "freehttp.Stream":
					this.error(freehttp.SuperResponseWriter.WriterStream(content))
				case "freehttp.File":
					freehttp.ServeFiles(content)
				case "error":
					this.error(content)
				default:
					this.error(fmt.Errorf("unsupported out type: %s", reType))
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
