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

// Slow handler
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

// Fast handler
func Fast(req *http.Request) wrap.R {
	return wrap.R{
		Err:        nil,
		StatusCode: http.StatusOK,
		Body:       strings.NewReader("Hello, World!"),
		Headers: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
	}
}

// Err handler
func Err(req *http.Request) wrap.R {
	return wrap.R{
		Err: errors.New("something went wrong"),
	}
}

// Empty handler
func Empty(req *http.Request) wrap.R {
	return wrap.R{}
}
