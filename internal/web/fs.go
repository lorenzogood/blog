package web

import (
	"errors"
	"fmt"

	"github.com/lorenzogood/blog/internal/assets/assetfs"
)

const (
	ImmutableCacheHeader = "public, max-age=31536000, immutable"
	WellKnownCacheHeader = "public, max-age=86400, immutable"
)

func FileServer(f *assetfs.AssetFS, cacheControl string) Handler {
	return fileServer(f, cacheControl, true)
}

func fileServer(f *assetfs.AssetFS, cacheControl string, sendEtag bool) Handler {
	return func(ctx *Ctx) error {
		path := ctx.PathValue("*")
		file, err := f.GetFile(path)
		if err != nil {
			if errors.Is(err, assetfs.ErrNotFound) {
				return ctx.RespondString(NotFound, "Not Found.")
			}

			return fmt.Errorf("http fileserver error: %w", err)
		}

		if sendEtag {
			if ctx.Request().Header.Get("If-None-Match") == file.Info.Etag {
				ctx.Header().Set("Cache-Control", cacheControl)
				ctx.SetStatus(NotModified)
				return nil
			}
		}

		ctx.Header().Set("Content-Type", file.Info.Mime)
		if sendEtag {
			ctx.Header().Set("Etag", file.Info.Etag)
		}
		ctx.Header().Set("Cache-Control", cacheControl)

		return ctx.Respond(OK, file.Content)
	}
}

func PermanentFileServer(f *assetfs.PermanentFS) Handler {
	return fileServer(&f.AssetFS, ImmutableCacheHeader, false)
}
