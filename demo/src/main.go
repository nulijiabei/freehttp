package main

import (
	"bufio"
	"fmt"
	"freehttp"
	"io"
	"os"
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

func (this *Web) Download() freehttp.File {
	return "freehttp/README.md"
}

func (this *Web) WriteReturn() freehttp.HttpStatus {
	return 404
}

func (this *Web) WriteStatus() freehttp.ContentType {
	return "image/jpeg"
}

func (this *Web) ReadStream(w *freehttp.ResponseWriter, stream freehttp.Stream) error {
	// Content-Type ...
	f, err := os.Open("freehttp/README.md")
	if err != nil {
		return err
	}
	if _, err := io.Copy(w.Writer, bufio.NewReader(f)); err != nil {
		return err
	}
	return nil
}

func (this *Web) WriteBufioStream() freehttp.Stream {
	return bufio.NewReader(strings.NewReader("..."))
}

func main() {

	s := freehttp.NewServer()
	if err := s.Register(new(Web)); err != nil {
		fmt.Println(err)
	}

	s.Start(":8080")

}
