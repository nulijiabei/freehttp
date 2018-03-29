package freehttp

import (
	"bufio"
)

type Json interface{}

type JsonIndent interface{}

// HttpStatusType -> int
type HttpStatus interface{}

func HttpStatusType(v interface{}) int {
	return v.(int)
}

// ContentTypeType -> string
type ContentType interface{}

func ContentTypeType(v interface{}) string {
	return v.(string)
}

// StreamType -> *bufio.Reader
type Stream interface{}

func StreamType(v interface{}) *bufio.Reader {
	return v.(*bufio.Reader)
}

type File interface{}

// FileType -> string
func FileType(v interface{}) string {
	return v.(string)
}

// RedirectType -> string
type Redirect interface{}

func RedirectType(v interface{}) string {
	return v.(string)
}
