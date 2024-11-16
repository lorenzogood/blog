package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lorenzogood/blog/internal/assets"
	"github.com/lorenzogood/blog/internal/watcher"
	"go.uber.org/zap"
)

type Renderer struct {
	templ *template.Template
}

func New(base string, a assets.LinkGetter) (*Renderer, error) {
	base = filepath.Clean(base)

	t := template.New("")
	t.Funcs(Funcs(a))

	err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("filepath walk error: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		name := strings.TrimPrefix(path, base+"/")

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		if _, err := t.New(name).Parse(string(content)); err != nil {
			return fmt.Errorf("template render error for %s: %w", path, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Renderer{
		templ: t,
	}, nil
}

func NewWatched(base string, a assets.LinkGetter) (*Renderer, error) {
	logger := zap.L().Named("templates_watched")

	var t Renderer
	rebuildFunc := func() {
		logger.Info("building")
		templ, err := New(base, a)
		if err != nil {
			logger.Error("failed to rebuild templates", zap.Error(err))
			return
		}
		t = *templ
		logger.Debug("built")
	}

	rebuildFunc()

	if err := watcher.NewWatcher(base, 100*time.Millisecond, rebuildFunc); err != nil {
		return nil, fmt.Errorf("error watching template dir: %w", err)
	}

	return &t, nil
}

func (r *Renderer) Render(w io.Writer, name string, data any) error {
	return r.templ.ExecuteTemplate(w, name, data)
}
