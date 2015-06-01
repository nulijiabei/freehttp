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
func (this *Request) ReadBody() Body {
	body, err := ioutil.ReadAll(this.Reader)
	if err != nil {
		return []byte{}
	}
	return body
}

// 读取全部 Body 数据转为 Json
func (this *Request) ReadBodyJson() BodyJson {
	data, err := ioutil.ReadAll(this.Reader)
	if err != nil {
		return map[string]interface{}{}
	}
	var content map[string]interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return map[string]interface{}{}
	}
	return content
}

// 读取 Bufio Stream
func (this *Request) ReadBufioStream() BufioStream {
	return bufio.NewReader(this.Reader)
}

// 读取 Content-Type
func (this *Request) ReadContentType() ContentType {
	return ContentType(this.Request.Header.Get("Content-Type"))
}
