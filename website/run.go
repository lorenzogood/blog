package website

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lorenzogood/blog/blog"
	"github.com/lorenzogood/blog/internal/assets/assetfs"
	"github.com/lorenzogood/blog/internal/server"
	"github.com/lorenzogood/blog/internal/templates"
	"github.com/lorenzogood/blog/internal/web"
	"github.com/lorenzogood/blog/internal/web/mid"
	"github.com/lorenzogood/blog/public"
)

type Config struct {
	Addr        string `conf:"default:0.0.0.0:3000"`
	AssetDir    string `conf:"required"`
	TemplateDir string `conf:"required"`
	Development bool   `conf:"-"`
}

func Run(ctx context.Context, c Config, b *blog.Blog) error {
	p, err := assetfs.New(public.Public, "assets")
	if err != nil {
		return fmt.Errorf("error opening public assets: %w", err)
	}

	var a *assetfs.PermanentFS
	var t *templates.Renderer
	if c.Development {
		a, err = assetfs.NewWatchedPermanent(c.AssetDir, "/assets")
		if err != nil {
			return fmt.Errorf("error opening development assets: %w", err)
		}

		t, err = templates.NewWatched(c.TemplateDir, a)
		if err != nil {
			return fmt.Errorf("error opening development templates: %w", err)
		}
	} else {
		a, err = assetfs.NewPermanent(os.DirFS(c.AssetDir), ".", "/assets")
		if err != nil {
			return fmt.Errorf("error opening assets: %w", err)
		}

		t, err = templates.New(c.TemplateDir, a)
		if err != nil {
			return fmt.Errorf("error opening templates: %w", err)
		}
	}

	r := chi.NewMux()
	r.Use(mid.LogRecover)
	r.Use(middleware.Compress(5))
	r.Method(http.MethodGet, "/*", web.FileServer(p, web.WellKnownCacheHeader))
	r.Method(http.MethodGet, "/assets/*", web.PermanentFileServer(a))
	r.Method(http.MethodGet, "/", web.Handler(func(ctx *web.Ctx) error {
		data := struct {
			SinglePageData
			Posts []blog.Post
		}{
			SinglePageData: SinglePageData{
				Content: b.Index,
			},
			Posts: b.Posts,
		}

		return ctx.RespondTemplate(t, web.OK, "index.tmpl.html", data)
	}))
	r.Method(http.MethodGet, "/archive", web.Handler(func(ctx *web.Ctx) error {
		data := struct {
			SinglePageData
			Posts []blog.Post
		}{
			SinglePageData: SinglePageData{
				Title:       "Posts",
				Description: "Post Archive",
			},
			Posts: b.Posts,
		}

		return ctx.RespondTemplate(t, web.OK, "archive.tmpl.html", data)
	}))
	r.Method(http.MethodGet, "/posts/*", web.Handler(func(ctx *web.Ctx) error {
		slug := ctx.PathValue("*")

		for _, v := range b.Posts {
			if v.Slug == slug {
				data := struct {
					SinglePageData
					Date    time.Time
					Updated *time.Time
				}{
					SinglePageData: SinglePageData{
						Title:       v.Title,
						Description: v.Description,
						Content:     v.Content,
						ProseExtra:  true,
					},
					Date:    v.Date,
					Updated: v.Updated,
				}
				return ctx.RespondTemplate(t, web.OK, "post.tmpl.html", data)
			}
		}

		return ctx.RespondString(web.NotFound, "Post Not Found.")
	}))
	r.Method(http.MethodGet, "/feed", web.Handler(func(ctx *web.Ctx) error {
		data := struct {
			SinglePageData
			Posts []blog.Feedpost
		}{
			SinglePageData: SinglePageData{
				Title:       "Feed",
				Description: "Short Form Content Feed",
				ProseExtra:  true,
			},
			Posts: b.Feedposts,
		}

		return ctx.RespondTemplate(t, web.OK, "feed.tmpl.html", data)
	}))

	server.Serve(ctx, c.Addr, r)

	return nil
}
