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

// Redirect
func (this *FreeHttp) Redirect(name string, content interface{}, def func(string, string) string) {
	url := string(content.(Redirect))
	if url[0] == ':' {
		url = def(name, url[1:])
	}
	http.Redirect(this.SuperResponseWriter.ResponseWriter, this.SuperRequest.Request, url, http.StatusFound)
}
