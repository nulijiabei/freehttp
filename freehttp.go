package freehttp

/*
	说明：
		同时使用到ResponseWriter和Request的在FreeHttp下进行封装
		只使用到Request的封装到SuperRequest中
		只使用到ResponseWriter的封装到SuperResponseWriter中
*/

import (
	"net/http"

	"./websocket"
)

// golang v1.9
// WebSocket
type WSConn = websocket.Conn

// FreeHttp
type FreeHttp struct {
	SuperResponseWriter *ResponseWriter
	SuperRequest        *Request
}

// New FreeHttp
func NewFreeHttp(w http.ResponseWriter, r *http.Request) *FreeHttp {
	freehttp := new(FreeHttp)
	freehttp.SuperResponseWriter = NewResponseWriter(w)
	freehttp.SuperRequest = NewRequest(r)
	return freehttp
}

// ServeFiles
func (this *FreeHttp) ServeFiles(content interface{}) {
	http.ServeFile(this.SuperResponseWriter.ResponseWriter, this.SuperRequest.Request, FileType(content))
}

// Redirect
func (this *FreeHttp) Redirect(content interface{}) {
	http.Redirect(this.SuperResponseWriter.ResponseWriter, this.SuperRequest.Request, RedirectType(content), http.StatusFound)
}

// WebSocket
func (this *FreeHttp) NewWebSokcet(handler func(*WSConn)) {
	s := websocket.Server{Handshake: websocket.CheckOrigin}
	s.ServeWebSocket(this.SuperResponseWriter.ResponseWriter, this.SuperRequest.Request, func(conn *WSConn) {
		handler(conn)
	})
}
