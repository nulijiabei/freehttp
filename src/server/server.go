package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type ResponseWriter struct {
	ResponseWriter http.ResponseWriter
}

type Request struct {
	Request *http.Request
}

type Server struct {
	name    string
	rcvr    reflect.Value
	typ     reflect.Type
	methods map[string]*Method
}

type Method struct {
	method reflect.Method
	json   bool
}

func NewServer() *Server {
	server := new(Server)
	server.methods = make(map[string]*Method)
	return server
}

func (this *Server) Start(port string) error {
	return http.ListenAndServe(port, this)
}

func (this *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for mname, mmethod := range this.methods {
		if strings.ToLower("/"+this.name+"."+mname) == r.URL.Path {
			if mmethod.json {
				returnValues := mmethod.method.Func.Call([]reflect.Value{this.rcvr, reflect.ValueOf(ResponseWriter{w}), reflect.ValueOf(Request{r})})
				content := returnValues[0].Interface()
				fmt.Println(content)
				if content != nil {
					data, err := json.MarshalIndent(content, "", "  ")
					if err != nil {
						w.WriteHeader(500)
						w.Write(data)
						return
					} else {
						w.Write(data)
					}
				}
			} else {
				mmethod.method.Func.Call([]reflect.Value{this.rcvr, reflect.ValueOf(ResponseWriter{w}), reflect.ValueOf(Request{r})})
			}
		}
	}
}

/*
	func (this *Hello) JsonHello(r server.Request) {}
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
		if strings.HasPrefix(mname, "Json") {
			if mtype.NumIn() != 2 {
				return fmt.Errorf("method %s has wrong number of ins: %d", mname, mtype.NumIn())
			}
			arg := mtype.In(1)
			if arg.String() != "server.Request" {
				return fmt.Errorf("%s argument type not exported: %s", mname, arg)
			}
			this.methods[mname] = &Method{method, true}
		} else {
			if mtype.NumIn() != 3 {
				return fmt.Errorf("method %s has wrong number of ins: %d", mname, mtype.NumIn())
			}
			reply := mtype.In(1)
			if reply.String() != "server.ResponseWriter" {
				return fmt.Errorf("%s argument type not exported: %s", mname, reply)
			}
			arg := mtype.In(2)
			if arg.String() != "server.Request" {
				return fmt.Errorf("%s argument type not exported: %s", mname, arg)
			}
			this.methods[mname] = &Method{method, false}
		}
	}
	return nil
}
