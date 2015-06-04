package freehttp

import (
	"bufio"
)

// Json
type Json interface{}

// Json Indent
type JsonIndent interface{}

// HTTP-Status
type HttpStatus int

// Content-Type
type ContentType string

// Bufio.Reader
type Stream *bufio.Reader

// File
type File string
