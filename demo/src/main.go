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

func (this *Web) RradConf(conf *freehttp.INI) {
	conf.Show()
	conf.Set("default", "freehttp", "initalize")
	conf.GetString("default.freehttp", "default value")
	conf.Del("default", "freehttp")
	conf.Save()
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

	ss := freehttp.NewServer()

	// mname = StructName
	// name  = FuncName
	def := func(mname, name string) string {
		// return string == r.URL.Path
		return strings.ToLower(fmt.Sprintf("/%s/%s", mname, name))
	}

	ss.InitURLPath(def)

	ss.InitConfig("/Users/nljb/profile")

	if err := ss.Register(new(Web)); err != nil {
		panic(err)
	}

	if err := ss.Start(":8080"); err != nil {
		panic(err)
	}

}
