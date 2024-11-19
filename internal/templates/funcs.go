package templates

import (
	"html/template"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/lorenzogood/blog/internal/assets"
)

func Funcs(a assets.LinkGetter) template.FuncMap {
	return template.FuncMap{
		"asset":     a.GetLink,
		"mark_safe": markSafe,
		"slugify":   slug.Make,
		"int":       strconv.Atoi,
	}
}

// Marks a string as html WITHOUT CLEANING IT. XSS!
func markSafe(s string) template.HTML {
	return template.HTML(s)
}
