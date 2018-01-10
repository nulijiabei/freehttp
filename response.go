package freehttp

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"
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
	this.ResponseWriter.WriteHeader(HttpStatusType(content))
}

// 转 Json 并回写
func (this *ResponseWriter) WriterJson(content interface{}) error {
	data, err := json.Marshal(content)
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
	this.ResponseWriter.Header().Set("Content-Type", ContentTypeType(content))
}

// 回写 Bufio Stream
func (this *ResponseWriter) WriterStream(content interface{}) error {
	if _, err := io.Copy(this.Writer, StreamType(content)); err != nil {
		return err
	}
	return this.Writer.Flush()
}

// 回写 文件
func (this *ResponseWriter) WriterFile(content interface{}) error {
	// Content-Type ...
	f, err := os.Open(FileType(content))
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(this.Writer, bufio.NewReader(f)); err != nil {
		return err
	}
	return this.Writer.Flush()
}
