package wrap

import (
	"fmt"
	"io"
	"net/http"

	"github.com/wayneashleyberry/hwrap/pkg/middleware"
	"go.uber.org/zap"
)

// R implementation
type R struct {
	Body       io.Reader
	Err        error
	Headers    map[string]string
	StatusCode int
}

// H implementation
func H(logger *zap.Logger, h func(r *http.Request) R) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger = logger.With(
			zap.String("request.id", middleware.GetReqID(ctx)),
		)

		res := make(chan R, 1)

		go func() {
			res <- h(r)
		}()

		select {
		case <-ctx.Done():
			logger.Error("timeout")
			return
		case r := <-res:
			write(logger, w, r)
		}
	}
}

func write(z *zap.Logger, w http.ResponseWriter, r R) {
	// Log errors if present
	if r.Err != nil {
		z.Error(r.Err.Error())
	}

	// Set a default status code
	if r.StatusCode < http.StatusContinue {
		r.StatusCode = http.StatusInternalServerError
	}

	// Write the headers
	if r.Headers != nil {
		for key, value := range r.Headers {
			w.Header().Set(key, value)
		}
	}

	// Write the status code
	w.WriteHeader(r.StatusCode)

	// Write the response body, if not present in the response we'll
	// generate a generic response like "500 Internal Server Error"
	switch {
	case r.Body != nil:
		io.Copy(w, r.Body)
	default:
		fmt.Fprintf(w, "%d %s", r.StatusCode, http.StatusText(r.StatusCode))
	}
}
