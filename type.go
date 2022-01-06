package tmplbuild

import "path/filepath"

type MediaType = string

const (
	UnknowMediaType MediaType = ""
	HtmlMediaType   MediaType = "html"
	CssMediaType    MediaType = "css"
	JsMediaType     MediaType = "js"
	ImageMediaType  MediaType = "image"
)

func GetMediaTypeByFilePath(path string) MediaType {
	ext := filepath.Ext(path)
	switch ext {
	case ".htm", ".html", ".tpl":
		return HtmlMediaType
	case ".js":
		return JsMediaType
	case ".css":
		return CssMediaType
	case ".ico", ".png", ".jpg", ".jpeg", ".svg", ".gif", ".bmp":
		return ImageMediaType
	}

	return UnknowMediaType
}
