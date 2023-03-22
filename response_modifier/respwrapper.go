package main

import (
	"io"
)

type ResponseWrapper interface {
	Data() map[string]interface{}
	Io() io.Reader
	IsComplete() bool
	Headers() map[string][]string
	StatusCode() int
}

type responseWrapper struct {
	data       map[string]interface{}
	isComplete bool
	io         io.Reader
	headers    map[string][]string
	statusCode int
}

func (r responseWrapper) Data() map[string]interface{} { return r.data }
func (r responseWrapper) IsComplete() bool             { return r.isComplete }
func (r responseWrapper) Io() io.Reader                { return r.io }
func (r responseWrapper) Headers() map[string][]string { return r.headers }
func (r responseWrapper) StatusCode() int              { return r.statusCode }
