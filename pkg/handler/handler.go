package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/wayneashleyberry/hwrap/pkg/wrap"
)

// Slow handler demonstrates streaming to the response writer, as well
// as timing out with an operation that doesn't get cancelled
func Slow(req *http.Request) wrap.R {
	time.Sleep(3 * time.Second)

	payload := map[string]string{
		"message": "Hello, World!",
	}

	r, w := io.Pipe()

	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		_ = w.Close()
	}()

	return wrap.R{
		Err:        nil,
		StatusCode: http.StatusOK,
		Body:       r,
		Headers: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
	}
}

// Fast handler shows a simpler way to return a response body
func Fast(req *http.Request) wrap.R {
	return wrap.R{
		StatusCode: http.StatusOK,
		Body:       strings.NewReader("Hello, World!"),
		Headers: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
	}
}

// Err handler only returns an error
func Err(req *http.Request) wrap.R {
	return wrap.R{
		Err: errors.New("something went wrong"),
	}
}

// Empty handler returns nothing
func Empty(req *http.Request) wrap.R {
	return wrap.R{}
}
