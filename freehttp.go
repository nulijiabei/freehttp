package freehttp

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
