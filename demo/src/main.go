package main

import (
	"bufio"
	"fmt"
	"freehttp"
	"strings"
)

type Web struct {
}

func (this *Web) ReadWrite(w *freehttp.ResponseWriter, r *freehttp.Request) {
	w.ResponseWriter.Write([]byte("print"))
}

func (this *Web) WriteJson() (freehttp.Json, freehttp.JsonIndent) {
	m := make(map[string]interface{})
	m["baidu"] = "www.baidu.com"
	return m, m
}

func (this *Web) ReadBody(t freehttp.ContentType) error {
	//	fmt.Println(string(body))
	fmt.Println(t)
	return fmt.Errorf("...")
}

func (this *Web) ReadBodyJson(bodyJson freehttp.BodyJson) {
	// bodyJson.(map[string]interface{})
}

func (this *Web) WriteReturn() freehttp.HttpStatus {
	return 404
}

func (this *Web) WriteStatus() freehttp.ContentType {
	return "image/jpeg"
}

func (this *Web) ReadBufioStream(stream freehttp.BufioStream) {
	// stream.(*bufio.Reader)
}

func (this *Web) WriteBufioStream() freehttp.BufioStream {
	return bufio.NewReader(strings.NewReader("..."))
}

func main() {

	s := freehttp.NewServer()
	if err := s.Register(new(Web)); err != nil {
		fmt.Println(err)
	}
	s.Start(":8080")

}
