package httputils

import (
	"fmt"
	"net/http"
)

type HttpResult[T any] struct {
	err        error
	statusCode int
	header     http.Header
	rawBody    []byte
	o          T
}

func (r HttpResult[T]) Err() error {
	return r.err
}

func (r HttpResult[T]) StatusCode() int {
	return r.statusCode
}

func (r HttpResult[T]) Header() http.Header {
	return r.header
}

func (r HttpResult[T]) ToObject() T {
	return r.o
}

func (r HttpResult[T]) RawBody() []byte {
	return r.rawBody
}

func (r HttpResult[T]) RawBodyString() string {
	return string(r.rawBody)
}

func (r HttpResult[T]) String() string {
	return fmt.Sprintf("err: %v, statusCode: %d, rawBody: %s, o: %v", r.err, r.statusCode, r.rawBody, r.o)
}
