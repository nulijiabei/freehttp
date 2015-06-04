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

func (this *Web) ReadWrite(rw *freehttp.FreeHttp) {
	rw.SuperResponseWriter.ResponseWriter.Write([]byte("print"))
}

func (this *Web) WriteJson() (freehttp.Json, freehttp.JsonIndent) {
	m := make(map[string]interface{})
	m["baidu"] = "www.baidu.com"
	return m, m
}

func (this *Web) Download(rw *freehttp.FreeHttp) freehttp.File {
	return "/Users/nljb/MyCore/git/github/freehttp/README.md"
}

func (this *Web) WriteReturn(w *freehttp.ResponseWriter) freehttp.HttpStatus {
	return 404
}

func (this *Web) WriteStatus() freehttp.ContentType {
	return "image/jpeg"
}

func (this *Web) ReadStream(rw *freehttp.FreeHttp, stream freehttp.Stream) error {
	// Content-Type ...
	f, err := os.Open("/Users/nljb/MyCore/git/github/freehttp/README.md")
	if err != nil {
		return err
	}
	if _, err := io.Copy(rw.SuperResponseWriter.Writer, bufio.NewReader(f)); err != nil {
		return err
	}
	return nil
}

func (this *Web) WriteStream() freehttp.Stream {
	return bufio.NewReader(strings.NewReader("..."))
}

func main() {

	s := freehttp.NewServer(func(mname, name string) string {
		return strings.ToLower(fmt.Sprintf("%s.%s", mname, name))
	})
	if err := s.Register(new(Web)); err != nil {
		fmt.Println(err)
	}

	s.Start(":8080")

}
