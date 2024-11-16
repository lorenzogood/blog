package assets

import "errors"

var ErrNotFound = errors.New("asset not found")

// Get a link for an asset from dist/*.
type LinkGetter interface {
	GetLink(path string) (string, error)
}
