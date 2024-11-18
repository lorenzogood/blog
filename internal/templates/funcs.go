package templates

import (
	"html/template"

	"github.com/lorenzogood/blog/internal/assets"
)

func Funcs(a assets.LinkGetter) template.FuncMap {
	return template.FuncMap{
		"asset":     a.GetLink,
		"mark_safe": markSafe,
	}
}

// Marks a string as html WITHOUT CLEANING IT. XSS!
func markSafe(s string) template.HTML {
	return template.HTML(s)
}
