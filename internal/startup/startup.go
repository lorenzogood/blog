package startup

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go.uber.org/zap"
)

// Check if debug logger should be enabled.
func isDebug() bool {
	val := os.Getenv("DEBUG")

	v, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return v
}

func Run() (context.Context, context.CancelFunc) {
	logLevel := zap.InfoLevel
	if isDebug() {
		logLevel = zap.DebugLevel
	}

	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(logLevel)

	l, err := cfg.Build(
		zap.WithCaller(false),
		zap.AddStacktrace(zap.PanicLevel),
	)
	if err != nil {
		panic(fmt.Errorf("error initializing logger: %w", err))
	}

	zap.ReplaceGlobals(l)

	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
}
