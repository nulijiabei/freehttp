package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// 包装 *http.Request
type Request struct {
	Request *http.Request
}

// 读取全部 Body 数据
func (this *Request) ReadAllBody() []byte {
	data, err := ioutil.ReadAll(this.Request.Body)
	if err != nil {
		return nil
	}
	return data
}

// 读取全部 Body 数据转为 Json
func (this *Request) ReadAllBodyJson() map[string]interface{} {
	data, err := ioutil.ReadAll(this.Request.Body)
	if err != nil {
		return nil
	}
	var content interface{}
	if err := json.Unmarshal(data, &content); err != nil {
		return nil
	}
	return content.(map[string]interface{})
}
