package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

func Serve(ctx context.Context, addr string, handler http.Handler) {
	logger := zap.L().With(zap.String("addr", addr))

	s := &http.Server{
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second, //nolint:mnd // fine
		ReadHeaderTimeout: 2 * time.Second,  //nolint:mnd // fine
		Handler:           handler,
		Addr:              addr,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start http server", zap.Error(err))
			os.Exit(1)
		}
	}()

	go func() {
		<-ctx.Done()

		sCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.Shutdown(sCtx); err != nil {
			logger.Warn("failed to gracefully stop http server", zap.Error(err))
		}

		logger.Debug("stopped http server")
	}()

	logger.Debug("started http server")
}
