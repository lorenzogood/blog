package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lorenzogood/blog/blog"
	"github.com/lorenzogood/blog/internal/config"
	"github.com/lorenzogood/blog/internal/startup"
	"github.com/lorenzogood/blog/website"
	"go.uber.org/zap"
)

type cfg struct {
	Development bool   `conf:"default:false"`
	ContentDir  string `conf:"required"`
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
	config.Parse("BLOG", &c)
	c.Web.Development = c.Development

	var b *blog.Blog
	var err error
	if c.Development {
		b, err = blog.NewWatched(c.ContentDir)
		if err != nil {
			return err
		}
	} else {
		b, err = blog.New(c.ContentDir)
		if err != nil {
			return fmt.Errorf("error building content: %w", err)
		}
	}

	if err := website.Run(ctx, c.Web, b); err != nil {
		return fmt.Errorf("error starting web service: %w", err)
	}

	return nil
}
