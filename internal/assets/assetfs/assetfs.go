package assetfs

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"maps"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	"go.uber.org/zap"
)

var ErrNotFound = errors.New("file not found")

type Info struct {
	Etag string
	Mime string
}

type AssetFS struct {
	base       fs.FS
	info       map[string]Info
	pathPrefix string
	webBase    string
}

func New(f fs.FS, base string) (*AssetFS, error) {
	logger := zap.L().Named("assetfs_build")
	info := make(map[string]Info)
	err := fs.WalkDir(f, base, func(path string, d fs.DirEntry, err error) error {
		filename := strings.TrimPrefix(path, base+"/")

		if err != nil {
			return fmt.Errorf("dir walk error: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		_, ext, found := strings.Cut(filepath.Base(path), ".")
		if !found {
			return fmt.Errorf("filename %s in invalid format (missing .type suffix)", path)
		}

		bytes, err := fs.ReadFile(f, path)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		hash := md5.Sum(bytes)
		etag := hex.EncodeToString(hash[:])
		mime := getFiletype(ext, bytes)

		logger.Debug(
			"added file",
			zap.String("path", filename),
			zap.String("hash", etag),
			zap.String("mime", mime),
		)

		info[filename] = Info{
			Etag: etag,
			Mime: mime,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	pathPrefix := base + "/"
	if base == "." {
		pathPrefix = ""
	}

	return &AssetFS{
		info:       info,
		base:       f,
		pathPrefix: pathPrefix,
	}, nil
}

type File struct {
	Info    Info
	Content []byte
}

func (a *AssetFS) GetFile(path string) (File, error) {
	m, ok := a.info[path]
	if !ok {
		return File{}, ErrNotFound
	}

	c, err := fs.ReadFile(a.base, a.pathPrefix+path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return File{}, ErrNotFound
		}

		return File{}, err
	}

	return File{
		Info:    m,
		Content: c,
	}, nil
}

func (a *AssetFS) Files() iter.Seq[string] {
	return maps.Keys(a.info)
}

func getFiletype(extension string, buf []byte) string {
	switch extension {
	case "txt":
		return "text/plain"
	case "css":
		return "text/css"
	case "js":
		return "application/javascript"
	case "json":
		return "application/json"
	case "webmanifest":
		return "application/json"
	case "ico":
		return "image/vnd.microsoft.icon"
	}

	kind, _ := filetype.Match(buf)
	if kind == filetype.Unknown {
		return "text/plain"
	}

	return kind.MIME.Value
}
