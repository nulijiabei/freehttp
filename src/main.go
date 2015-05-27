package main

import (
	"fmt"
	"server"
)

type Hello struct {
}

func (this *Hello) Print(w server.ResponseWriter, r server.Request) map[string]interface{} {
	w.ResponseWriter.Write([]byte("print"))
	return nil
}

func (this *Hello) Hello(w server.ResponseWriter, r server.Request) {
	w.ResponseWriter.Write([]byte("hellp"))
}

func (this *Hello) JsonHello(r server.Request) {

}

func main() {

	s := server.NewServer()
	fmt.Println(s.Register(new(Hello)))
	s.Start(":8080")

}
