package handler

import "io"

// CustomResponse definition
type CustomResponse struct {
	statusCode int
	err        error
	body       io.Reader
	headers    map[string]string
}

// StatusCode implementation
func (r CustomResponse) StatusCode() int {
	return r.statusCode
}

// Err implementation
func (r CustomResponse) Err() error {
	return r.err
}

// Body implementation
func (r CustomResponse) Body() io.Reader {
	return r.body
}

// Headers implementation
func (r CustomResponse) Headers() map[string]string {
	return r.headers
}
