package main

import (
	"fmt"
	"freehttp"
)

type Web struct {
}

func (this *Web) ReadWrite(w *freehttp.ResponseWriter, r *freehttp.Request) {
	// r.Request.PostForm
	w.ResponseWriter.Write([]byte("print"))
}

func (this *Web) WriteJson() freehttp.Json {
	m := make(map[string]interface{})
	m["baidu"] = "www.baidu.com"
	return m
}

func (this *Web) Hello(w *freehttp.ResponseWriter, r *freehttp.Request, body freehttp.Body, bodyJson freehttp.BodyJson) (freehttp.Status, error) {
	fmt.Println(body)
	fmt.Println(bodyJson)
	return 404, fmt.Errorf("...")
}

func main() {

	s := freehttp.NewServer()
	if err := s.Register(new(Web)); err != nil {
		fmt.Println(err)
	}
	s.Start(":8080")

}
