package main

import (
	"fmt"
	"server"
)

type Web struct {
}

func (this *Web) Print(w server.ResponseWriter, r server.Request) {
	w.ResponseWriter.Write([]byte("print"))
}

func (this *Web) Hello(w server.ResponseWriter, r server.Request) {
	w.ResponseWriter.Write([]byte("hello"))
}

func (this *Web) JsonHello(r server.Request) map[string]interface{} {
	m := make(map[string]interface{})
	m["baidu"] = "www.baidu.com"
	return m
}

func main() {

	s := server.NewServer()
	fmt.Println(s.Register(new(Web)))
	s.Start(":8080")

}
