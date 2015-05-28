package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// 包装 http.ResponseWriter
type ResponseWriter struct {
	ResponseWriter http.ResponseWriter
}

// 包装 *http.Request
type Request struct {
	Request *http.Request
}

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

// 启动服务
func (this *Server) Start(port string) error {
	return http.ListenAndServe(port, this)
}

// 内部方法 http 需要
func (this *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := false
	for name, method := range this.methods {
		if strings.ToLower(fmt.Sprintf("/%s.%s", this.name, name)) == strings.ToLower(r.URL.Path) {
			status = true
			value := make([]reflect.Value, method.Type.NumIn())
			value[0] = this.rcvr
			for n := 1; n < method.Type.NumIn(); n++ {
				switch method.Type.In(n).String() {
				case "server.Request":
					value[n] = reflect.ValueOf(Request{r})
				case "server.ResponseWriter":
					value[n] = reflect.ValueOf(ResponseWriter{w})
				}
			}
			returnValues := method.Func.Call(value)
			if method.Type.NumOut() == 1 {
				if content := returnValues[0].Interface(); content != nil {
					reType := method.Type.Out(0).String()
					switch reType {
					case "error":
						fmt.Println(content.(error).Error())
					default:
						fmt.Printf("unsupported return type: %s\n", reType)
					}
				}
			}
		}
	}
	if !status {
		http.NotFound(w, r)
	}
}

/*
	// 一个以Json开头的命令的函数, 返回的map[string]interface{}自动被处理成Json发送
	func (this *Hello) JsonHello(r server.Request) map[string]interface{} {}
	// 非以Json开头的命令的函数, 则默认为HTTP函数
	func (this *Hello) Hello(w server.ResponseWriter, r server.Request) {}
*/
func (this *Server) Register(rcvr interface{}) error {
	this.typ = reflect.TypeOf(rcvr)
	this.rcvr = reflect.ValueOf(rcvr)
	this.name = reflect.Indirect(this.rcvr).Type().Name()
	if this.name == "" {
		return fmt.Errorf("no service name for type ", this.typ.String())
	}
	for m := 0; m < this.typ.NumMethod(); m++ {
		method := this.typ.Method(m)
		mtype := method.Type
		mname := method.Name
		for n := 1; n < mtype.NumIn(); n++ {
			switch mtype.In(n).String() {
			case "server.Request":
			case "server.ResponseWriter":
			default:
				return fmt.Errorf("%s argument type not exported: Request or ResponseWriter", mname)
			}
		}
		this.methods[mname] = method
	}
	return nil
}

// 将 map[string]interface{} 转 Json 并排版回写
func (this *ResponseWriter) WriterJsonIndent(content map[string]interface{}) (int, error) {
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return 0, err
	}
	return this.ResponseWriter.Write(data)
}

// 将 map[string]interface{} 转 Json 回写 (一行一数据)
func (this *ResponseWriter) WriterJsonLine(content map[string]interface{}) (int, error) {
	data, err := json.Marshal(content)
	if err != nil {
		return 0, err
	}
	return this.ResponseWriter.Write(append(data, '\n'))
}
