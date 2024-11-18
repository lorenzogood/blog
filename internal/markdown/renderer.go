package markdown

import (
	"bytes"
	"fmt"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

type Renderer interface {
	RenderMarkdown(markdown []byte) (string, error)
}

type renderer struct {
	m goldmark.Markdown
}

func (r *renderer) RenderMarkdown(markdown []byte) (string, error) {
	var b bytes.Buffer
	if err := r.m.Convert(markdown, &b); err != nil {
		return "", fmt.Errorf("goldmark error: %w", err)
	}

	return b.String(), nil
}

func NewPostRenderer() Renderer {
	md := goldmark.New(
		goldmark.WithExtensions(mathjax.MathJax),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	return &renderer{
		m: md,
	}
}
