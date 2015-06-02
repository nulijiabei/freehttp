package freehttp

import (
	"bufio"
)

// Json
type Json map[string]interface{}

// Json Indent
type JsonIndent map[string]interface{}

// HTTP-Status
type HttpStatus int

// Content-Type
type ContentType string

// Bufio.Reader
type Stream *bufio.Reader

// File
type File string
