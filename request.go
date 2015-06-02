package freehttp

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// *http.Request
type Request struct {
	Request *http.Request
	Reader  *bufio.Reader
}

// New Request
func NewRequest(r *http.Request) *Request {
	request := new(Request)
	request.Request = r
	request.Reader = bufio.NewReader(r.Body)
	return request
}

// Close Request
func (this *Request) Close() {
	this.Request.Body.Close()
}

// 读取全部 Body 数据转为 Json
func (this *Request) ReadJson() Json {
	data, err := ioutil.ReadAll(this.Reader)
	if err != nil {
		return map[string]interface{}{}
	}
	var content interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return map[string]interface{}{}
	}
	return content.(map[string]interface{})
}

// 读取 Bufio Stream
func (this *Request) ReadStream() Stream {
	return bufio.NewReader(this.Reader)
}

// 读取 Content-Type
func (this *Request) ReadContentType() ContentType {
	return ContentType(this.Request.Header.Get("Content-Type"))
}
