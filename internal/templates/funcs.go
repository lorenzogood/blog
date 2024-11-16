package templates

import (
	"html/template"

	"github.com/lorenzogood/blog/internal/assets"
)

func Funcs(a assets.LinkGetter) template.FuncMap {
	return template.FuncMap{
		"asset": a.GetLink,
	}
}
