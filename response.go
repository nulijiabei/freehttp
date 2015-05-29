package freehttp

import (
	"encoding/json"
	"net/http"
)

// 包装 http.ResponseWriter
type ResponseWriter struct {
	ResponseWriter http.ResponseWriter
}

// 回写 HTTP Status
func (this *ResponseWriter) WriteHeader(status int) {
	this.ResponseWriter.WriteHeader(status)
}

// 将 map[string]interface{} 转 Json 并回写
func (this *ResponseWriter) WriterJson(content interface{}) error {
	data, err := json.Marshal(content.(Json))
	if err != nil {
		return err
	}
	_, err = this.ResponseWriter.Write(data)
	return err
}

// 将 map[string]interface{} 转 Json 并回写
func (this *ResponseWriter) WriterJsonIndent(content interface{}) error {
	data, err := json.MarshalIndent(content.(JsonIndent), "", "  ")
	if err != nil {
		return err
	}
	_, err = this.ResponseWriter.Write(data)
	return err
}
