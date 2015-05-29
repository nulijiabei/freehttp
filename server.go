package freehttp

import (
	"fmt"
	"net/http"
	"reflect"
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
		fmt.Println(err.Error())
	}
}

// 启动服务
func (this *Server) Start(port string) error {
	return http.ListenAndServe(port, this)
}

// 内部方法 http 需要
func (this *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := false
	request := &Request{r}
	responseWriter := &ResponseWriter{w}
	for name, method := range this.methods {
		if strings.ToLower(fmt.Sprintf("/%s.%s", this.name, name)) == strings.ToLower(r.URL.Path) {
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
				case "freehttp.Body":
					value[n] = reflect.ValueOf(request.ReadAllBody())
				case "freehttp.BodyJson":
					value[n] = reflect.ValueOf(request.ReadAllBodyJson())
				default:
					this.Error(fmt.Errorf("unsupported in type: %s\n", inType))
				}
			}
			returnValues := method.Func.Call(value)
			for t := 0; t < method.Type.NumOut(); t++ {
				reType := method.Type.Out(t).String()
				switch reType {
				case "freehttp.Status":
					responseWriter.WriteHeader(returnValues[t].Interface().(int))
				case "freehttp.Json":
					this.Error(responseWriter.WriterJson(returnValues[t].Interface()))
				case "freehttp.JsonIndent":
					this.Error(responseWriter.WriterJsonIndent(returnValues[t].Interface()))
				case "error":
					this.Error(returnValues[t].Interface().(error))
				default:
					this.Error(fmt.Errorf("unsupported out type: %s\n", reType))
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
		return fmt.Errorf("no service name for type ", this.typ.String())
	}
	for m := 0; m < this.typ.NumMethod(); m++ {
		method := this.typ.Method(m)
		// mtype := method.Type
		mname := method.Name
		this.methods[mname] = method
	}
	return nil
}
