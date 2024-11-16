package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lorenzogood/blog/internal/config"
	"github.com/lorenzogood/blog/internal/startup"
	"github.com/lorenzogood/blog/website"
	"go.uber.org/zap"
)

type cfg struct {
	Development bool `conf:"default:false"`
	Web         website.Config
}

func main() {
	ctx, cancel := startup.Run()
	defer cancel()

	if err := run(ctx); err != nil {
		zap.L().Error("failed to start application", zap.Error(err))
		os.Exit(1)
	}

	<-ctx.Done()
}

func run(ctx context.Context) error {
	var c cfg
	config.Parse("NEWSLETTER", &c)
	c.Web.Development = c.Development

	if err := website.Run(ctx, c.Web); err != nil {
		return fmt.Errorf("error starting web service: %w", err)
	}

	return nil
}
