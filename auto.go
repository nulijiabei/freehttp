package freehttp

import (
	"bufio"
)

// Json 普通格式
type Json map[string]interface{}

// Json 排版格式
type JsonIndent map[string]interface{}

// Body
type Body []byte

// Json Body
type BodyJson map[string]interface{}

// HTTP-Status
type HttpStatus int

// Content-Type
type ContentType string

// Bufio.Reader
type BufioStream *bufio.Reader
