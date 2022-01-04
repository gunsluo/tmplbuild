package tmplbuild

import (
	"io"
	"path/filepath"
)

type Input struct {
	Path   string
	Reader io.Reader
}

func GetMediaTypeByFilePath(path string) MediaType {
	ext := filepath.Ext(path)
	switch ext {
	case ".html":
		return HtmlMediaType
	case ".js":
		return JsMediaType
	case ".css":
		return CssMediaType
	case ".ico":
		fallthrough
	case ".png":
		fallthrough
	case ".jpg":
		fallthrough
	case ".jpeg":
		fallthrough
	case ".svg":
		return ImageMediaType
	}

	return UnknowMediaType
}
