package freehttp

import (
	"bufio"
)

type Json interface{}

type JsonIndent interface{}

type HttpStatus interface{}

func HttpStatusType(v interface{}) int {
	return v.(int)
}

type ContentType interface{}

func ContentTypeType(v interface{}) string {
	return v.(string)
}

type Stream interface{}

func StreamType(v interface{}) *bufio.Reader {
	return v.(*bufio.Reader)
}

type File interface{}

func FileType(v interface{}) string {
	return v.(string)
}

type Redirect interface{}

func RedirectType(v interface{}) string {
	return v.(string)
}
