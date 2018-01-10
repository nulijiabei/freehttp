package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"../../../freehttp"
)

type Web struct {
}

//func (this *Web) RradConf(conf *freehttp.INI) {
//	conf.Show()
//	conf.Set("default", "freehttp", "initalize")
//	conf.GetString("default.freehttp", "default value")
//	conf.Del("default", "freehttp")
//	conf.Save()
//}

func (this *Web) Redirect() freehttp.Redirect {
	return "http://www.baidu.com"
}

func (this *Web) ReadWrite(rw *freehttp.FreeHttp) {
	rw.SuperRequest.Request.ParseForm()
	fmt.Println("->", rw.SuperRequest.Request.FormValue("val"))
	rw.SuperResponseWriter.ResponseWriter.Write([]byte("print"))
}

func (this *Web) WriteJson() (freehttp.Json, freehttp.JsonIndent) {
	m := make(map[string]interface{})
	m["baidu"] = "www.baidu.com"
	return m, m
}

func (this *Web) Download(rw *freehttp.FreeHttp) freehttp.File {
	return "E:\\MyCore\\git\\github\\freehttp\\README.md"
}

func (this *Web) WriteReturn(w *freehttp.ResponseWriter) freehttp.HttpStatus {
	return 404
}

func (this *Web) WriteStatus() freehttp.ContentType {
	return "image/jpeg"
}

func (this *Web) ReadJson(json *freehttp.Json) {
	if json == nil {
		fmt.Println("hahaha")
	}
}

func (this *Web) ReadConfig(conf *freehttp.INI) {
	if conf == nil {
		fmt.Println("hahaha")
	}
}

func (this *Web) ReadStream(rw *freehttp.FreeHttp, stream freehttp.Stream) error {
	f, err := os.OpenFile("E:\\MyCore\\git\\github\\freehttp\\a", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if nil != err {
		panic(err)
	}
	defer f.Close()
	if _, err := io.Copy(bufio.NewWriter(f), freehttp.StreamType(stream)); err != nil {
		return err
	}
	return nil
}

func (this *Web) WriteStream() freehttp.Stream {
	return bufio.NewReader(strings.NewReader("..."))
}

// http://127.0.0.1:8080/MyStructName/MyFuncName
func main() {

	service := freehttp.NewService(new(Web))
	//	service.Config("/profile")

	// 启动服务器
	if err := service.Start(":8080"); err != nil {
		panic(err)
	}

}
