package frontmatter

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

var ErrNotFound = errors.New("frontmatter not found")

// Parse frontmatter for s, return rest, error.
func Parse(i []byte, s any) ([]byte, error) {
	const delim = "+++"

	parts := strings.SplitN(string(i), delim, 3)
	if len(parts) != 3 {
		return nil, ErrNotFound
	}

	fraw := parts[1]
	if _, err := toml.NewDecoder(bytes.NewReader([]byte(fraw))).Decode(s); err != nil {
		return nil, fmt.Errorf("toml decode error: %w", err)
	}

	return []byte(parts[2]), nil
}
