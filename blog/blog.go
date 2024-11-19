package blog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lorenzogood/blog/internal/markdown"
	"github.com/lorenzogood/blog/internal/watcher"
	"go.uber.org/zap"
)

type Blog struct {
	Index     string
	About     string
	Posts     []Post
	Feedposts []Feedpost
}

// Render a single markdown file, without frontmatter.
// Useful for Index, about page etc.
func renderFile(r markdown.Renderer, path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error opening file %s: %w", path, err)
	}

	html, err := r.RenderMarkdown(f)
	if err != nil {
		return "", fmt.Errorf("markdown render error on file %s: %w", path, err)
	}

	return html, nil
}

func New(base string) (*Blog, error) {
	base = filepath.Clean(base)

	pr := markdown.NewPostRenderer()

	var b Blog
	var err error

	b.Index, err = renderFile(pr, filepath.Join(base, "index.md"))
	if err != nil {
		return nil, err
	}

	b.About, err = renderFile(pr, filepath.Join(base, "about.md"))
	if err != nil {
		return nil, err
	}

	b.Posts, err = ParsePosts(pr, filepath.Join(base, "posts"))
	if err != nil {
		return nil, err
	}

	b.Feedposts, err = ParseFeedposts(pr, filepath.Join(base, "feedposts.toml"))
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func NewWatched(base string) (*Blog, error) {
	logger := zap.L().Named("blog_watched")

	var a Blog
	rebuildFunc := func() {
		logger.Info("building")
		b, err := New(base)
		if err != nil {
			logger.Error("failed to rebuild blog", zap.Error(err))
			return
		}
		a = *b
		logger.Debug("built")
	}

	rebuildFunc()

	if err := watcher.NewWatcher(base, 100*time.Millisecond, rebuildFunc); err != nil {
		return nil, fmt.Errorf("error watching blog dir: %w", err)
	}

	return &a, nil
}
