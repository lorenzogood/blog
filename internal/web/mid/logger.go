package mid

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func LogRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			if err := recover(); err != nil {
				zap.L().Error(
					"panic in http request handler",
					zap.Int("status", ww.Status()),
					zap.Duration("duration", time.Since(start)),
					zap.String("path", r.URL.Path),
					zap.Any("panic", err),
					zap.String("user_agent", r.Header.Get("User-Agent")),
				)
			} else {
				mills := time.Since(start) / time.Millisecond
				logFn := zap.L().Debug

				if mills >= 500 {
					logFn = zap.L().Info
				}
				if ww.Status() >= 500 {
					logFn = zap.L().Error
				}

				logFn(
					"http request",
					zap.Int("status", ww.Status()),
					zap.Duration("duration", time.Since(start)/time.Millisecond),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("user_agent", r.Header.Get("User-Agent")),
				)
			}
		}()

		next.ServeHTTP(ww, r)
	})
}
