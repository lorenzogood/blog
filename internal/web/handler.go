package web

import (
	"net/http"

	"go.uber.org/zap"
)

type Handler func(ctx *Ctx) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Ctx{
		w: w,
		r: r,
	}

	if err := h(c); err != nil {
		zap.L().Error("unhandled web error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Internal Server Error")); err != nil {
			zap.L().Error("failed to send web error response", zap.Error(err))
			return
		}
		return
	}
}
