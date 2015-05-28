package main

import (
	"fmt"
	"server"
)

type Web struct {
}

func (this *Web) ReadWrite(w server.ResponseWriter, r server.Request) {
	// r.Request.PostForm
	w.ResponseWriter.Write([]byte("print"))
}

func (this *Web) WriteJson() server.Json {
	m := make(map[string]interface{})
	m["baidu"] = "www.baidu.com"
	return m
}

func (this *Web) Hello(body server.Body, bodyJson server.BodyJson) error {
	fmt.Println(body)
	fmt.Println(bodyJson)
	return fmt.Errorf("...")
}

func main() {

	s := server.NewServer()
	if err := s.Register(new(Web)); err != nil {
		fmt.Println(err)
	}
	s.Start(":8080")

}
