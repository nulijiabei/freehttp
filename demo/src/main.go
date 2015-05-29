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

func (this *Web) ReadBody(body freehttp.Body, bodyJson freehttp.BodyJson) error {
	return fmt.Errorf("what you see is a error")
}

func (this *Web) WriteReturn() (freehttp.HttpStatus, freehttp.ContentType) {
	return 200, "image/jpeg"
}

func (this *Web) WriteBufioStream() freehttp.BufioStream {
	return bufio.NewReader(strings.NewReader("what you see is a stream"))
}

func (this *Web) ReadBufioStream(stream freehttp.BufioStream) {
}

func main() {

	s := freehttp.NewServer()
	if err := s.Register(new(Web)); err != nil {
		fmt.Println(err)
	}
	s.Start(":8080")

}
