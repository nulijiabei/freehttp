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

func (this *Web) Ifconfig(r server.Request) {

}

func (this *Web) Json(w server.ResponseWriter) {
	m := make(map[string]interface{})
	m["baidu"] = "www.baidu.com"
	w.WriterJsonLine(m)
}

func main() {

	s := server.NewServer()
	if err := s.Register(new(Web)); err != nil {
		fmt.Println(err)
	}
	s.Start(":8080")

}
