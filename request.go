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

// 读取全部 Body 数据
func (this *Request) ReadBody() []byte {
	body, err := ioutil.ReadAll(this.Reader)
	if err != nil {
		return make([]byte, 0)
	}
	return body
}

// 读取全部 Body 数据转为 Json
func (this *Request) ReadBodyJson() interface{} {
	data, err := ioutil.ReadAll(this.Reader)
	if err != nil {
		return make(map[string]interface{})
	}
	var content interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return make(map[string]interface{})
	}
	return content
}

// 读取 Bufio Stream
func (this *Request) ReadBufioStream() *bufio.Reader {
	return bufio.NewReader(this.Reader)
}
