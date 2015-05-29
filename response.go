package freehttp

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
)

// http.ResponseWriter
type ResponseWriter struct {
	ResponseWriter http.ResponseWriter
	Writer         *bufio.Writer
}

// New ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	responseWriter := new(ResponseWriter)
	responseWriter.ResponseWriter = w
	responseWriter.Writer = bufio.NewWriter(w)
	return responseWriter
}

// 回写 HTTP Status
func (this *ResponseWriter) WriteHeader(content interface{}) {
	this.ResponseWriter.WriteHeader(content.(int))
}

// 转 Json 并回写
func (this *ResponseWriter) WriterJson(content interface{}) error {
	data, err := json.Marshal(content.(Json))
	if err != nil {
		return err
	}
	_, err = this.ResponseWriter.Write(data)
	return err
}

// 转 Json 并回写（排版）
func (this *ResponseWriter) WriterJsonIndent(content interface{}) error {
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	_, err = this.ResponseWriter.Write(data)
	return err
}

// 回写 ContentType
func (this *ResponseWriter) SetContentType(content interface{}) {
	this.ResponseWriter.Header().Set("Content-Type", content.(string))
}

// 回写 Bufio Stream
func (this *ResponseWriter) WriterBufioStream(content interface{}) error {
	if _, err := io.Copy(this.Writer, content.(*bufio.Reader)); err != nil {
		return err
	}
	return this.Writer.Flush()
}
