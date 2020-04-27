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

	r.GET("/slow", wrap.H(logger, handler.Slow))
	r.GET("/fast", wrap.H(logger, handler.Fast))
	r.GET("/err", wrap.H(logger, handler.Err))
	r.GET("/empty", wrap.H(logger, handler.Err))
	r.GET("/warn", wrap.H(logger, handler.Warn))

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
