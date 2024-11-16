package assetfs

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/lorenzogood/blog/internal/assets"
)

type PermanentFS struct {
	AssetFS AssetFS

	baseNames map[string]string
	webBase   string
}

func newPermanent(a *AssetFS, webBase string) (*PermanentFS, error) {
	baseNames := make(map[string]string)
	for file := range a.Files() {
		baseName, err := getBaseName(file)
		if err != nil {
			return nil, fmt.Errorf("error getting basename for file: %w", err)
		}

		baseNames[baseName] = file
	}

	return &PermanentFS{
		AssetFS:   *a,
		webBase:   webBase,
		baseNames: baseNames,
	}, nil
}

func NewPermanent(f fs.FS, base, webBase string) (*PermanentFS, error) {
	a, err := New(f, base)
	if err != nil {
		return nil, err
	}

	return newPermanent(a, webBase)
}

func (p *PermanentFS) GetLink(name string) (string, error) {
	fullname, ok := p.baseNames[name]
	if !ok {
		return "", assets.ErrNotFound
	}

	return fmt.Sprintf("%s/%s", p.webBase, fullname), nil
}

func getBaseName(s string) (string, error) {
	front, ext, ok := strings.Cut(s, ".")
	if !ok {
		return "", fmt.Errorf("filename missing .ext suffix: %s", s)
	}

	name, _, ok := strings.Cut(front, "~")
	if !ok {
		return "", fmt.Errorf("filename missing ~{hash}: %s", s)
	}

	return fmt.Sprintf("%s.%s", name, ext), nil
}
