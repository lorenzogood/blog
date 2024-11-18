package blog

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/gosimple/slug"
	"github.com/lorenzogood/blog/internal/frontmatter"
	"github.com/lorenzogood/blog/internal/markdown"
	"go.uber.org/zap"
)

type Post struct {
	Title       string
	Description string
	Date        time.Time
	Updated     *time.Time
	Content     string
	Slug        string
}

type front struct {
	Title       string `toml:"title"`
	Description string `toml:"description"`
	Date        string `toml:"date"`
	Updated     string `toml:"updated,omitempty"`
}

func ParsePost(r markdown.Renderer, path string) (Post, error) {
	path = filepath.Clean(path)
	f, err := os.ReadFile(path)
	if err != nil {
		return Post{}, fmt.Errorf("error reading file %s: %w", path, err)
	}

	var fr front
	rest, err := frontmatter.Parse(f, &fr)
	if err != nil {
		return Post{}, fmt.Errorf("error parsing frontmatter for post %s: %w", path, err)
	}

	content, err := r.RenderMarkdown(rest)
	if err != nil {
		return Post{}, fmt.Errorf("markdown render error for post %s: %w", path, err)
	}

	date, err := time.Parse(time.DateTime, fr.Date)
	if err != nil {
		return Post{}, fmt.Errorf("error parsing date for post %s: %w", path, err)
	}

	var updated *time.Time
	if fr.Updated != "" {
		u, err := time.Parse(time.DateTime, fr.Date)
		if err != nil {
			return Post{}, fmt.Errorf("error parsing update date for post %s: %w", path, err)
		}

		updated = &u
	}

	slug := fmt.Sprintf("%d/%d/%s", date.UTC().Year(), date.UTC().Month(), slug.Make(fr.Title))

	p := Post{
		Title:       fr.Title,
		Description: fr.Description,
		Date:        date,
		Updated:     updated,
		Content:     content,
		Slug:        slug,
	}

	return p, nil
}

func ParsePosts(r markdown.Renderer, path string) ([]Post, error) {
	logger := zap.L().Named("post_build")
	path = filepath.Clean(path)
	var posts []Post
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir error: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		p, err := ParsePost(r, path)
		if err != nil {
			return err
		}

		logger.Debug(
			"added post",
			zap.String("title", p.Title),
			zap.String("slug", p.Slug),
		)

		posts = append(posts, p)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return posts, nil
}
