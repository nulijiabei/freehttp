package freehttp

/*
	说明：
		同时使用到ResponseWriter和Request的在FreeHttp下进行封装
		只使用到Request的封装到SuperRequest中
		只使用到ResponseWriter的封装到SuperResponseWriter中
*/

import (
	"net/http"
)

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
	http.ServeFile(this.SuperResponseWriter.ResponseWriter, this.SuperRequest.Request, string(content.(File)))
}

// Redirect
func (this *FreeHttp) Redirect(content interface{}) {
	http.Redirect(this.SuperResponseWriter.ResponseWriter, this.SuperRequest.Request, string(content.(Redirect)), http.StatusFound)
}
