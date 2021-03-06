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

// -----------------------
// WebSocket 支持

func (this *Web) WebSocket(rw *freehttp.FreeHttp) {
	// HTTP -> WebSocket
	rw.NewWebSokcet(func(conn *freehttp.WSConn) {

		// WSConn = websocket.Conn
		conn.Write([]byte("Hello WebSokcet !!!"))
		r := bufio.NewReader(conn)
		for {
			v, err := r.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
				break
			}
			conn.Write(v)
			conn.Write([]byte("\n"))
		}
	})
}

// -----------------------
// 多结构类支持

type Support struct {
}

func (this *Support) Redirect() freehttp.Redirect {
	return "http://www.nljb.net"
}

// -----------------------

func main() {

	// New Service
	service := freehttp.NewService()

	// Web
	web := new(Web)
	service.Register(web)
	service.Router("/baidu", web.WriteJson)
	service.Router("/download", web.Download)
	service.Router("/readjson", web.ReadJson)
	service.Router("/writestream", web.WriteStream)
	service.Router("/redirect", web.Redirect)
	service.Router("/websocket", web.WebSocket)
	service.Router("/readwrite", web.ReadWrite)

	// 多结构类支持
	support := new(Support)
	service.Register(support)
	service.Router("/nljb", support.Redirect)

	// 启动服务器
	if err := service.Start(":8080"); err != nil {
		panic(err)
	}

}
