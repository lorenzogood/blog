package blog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/lorenzogood/blog/internal/markdown"
)

type Feedpost struct {
	Title     string
	Published time.Time
	Content   string
}

type feedpost struct {
	Title     string `toml:"title"`
	Content   string `toml:"content"`
	Published string `toml:"published"`
}

func ParseFeedposts(pr markdown.Renderer, path string) ([]Feedpost, error) {
	path = filepath.Clean(path)

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}

	s := struct {
		Feedposts []feedpost `toml:"post"`
	}{}
	if _, err := toml.NewDecoder(bytes.NewReader(f)).Decode(&s); err != nil {
		return nil, fmt.Errorf("error reading feedpost file %s: %w", path, err)
	}

	posts := make([]Feedpost, len(s.Feedposts))
	for v, p := range s.Feedposts {
		content, err := pr.RenderMarkdown([]byte(p.Content))
		if err != nil {
			return nil, fmt.Errorf("error rendering content for feedpost %s: %w", p.Title, err)
		}

		published, err := time.Parse(time.DateTime, p.Published)
		if err != nil {
			return nil, fmt.Errorf("error parsing published time for feedpost %s: %w", p.Title, err)
		}

		posts[v] = Feedpost{
			Title:     p.Title,
			Content:   content,
			Published: published,
		}
	}

	return posts, nil
}
