package main

import (
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/justinas/alice"
	"github.com/wayneashleyberry/hwrap/pkg/handler"
	"github.com/wayneashleyberry/hwrap/pkg/middleware"
	"github.com/wayneashleyberry/hwrap/pkg/wrap"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	must(err)

	logger.Debug("starting up")

	r := httptreemux.NewContextMux()

	w := wrap.New(logger)

	r.GET("/slow", w.H(handler.Slow))
	r.GET("/fast", w.H(handler.Fast))
	r.GET("/err", w.H(handler.Err))
	r.GET("/empty", w.H(handler.Err))
	r.GET("/warn", w.H(handler.Warn))

	chain := alice.New(
		middleware.RequestID,
		middleware.Timeout(2*time.Second),
	).Then(r)

	err = http.ListenAndServe(":8080", chain)
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
