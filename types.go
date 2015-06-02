package freehttp

import (
	"bufio"
)

// Json 普通格式
type Json map[string]interface{}

// Json 排版格式
type JsonIndent map[string]interface{}

// HTTP-Status
type HttpStatus int

// Content-Type
type ContentType string

// Cookie
type Cookie string

// Bufio.Reader
type Stream *bufio.Reader

type Stream struct {
	Reader *bufio.Reader
}
