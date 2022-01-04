package tmplbuild

type MediaType = string

const (
	UnknowMediaType MediaType = ""
	HtmlMediaType   MediaType = "html"
	CssMediaType    MediaType = "css"
	JsMediaType     MediaType = "js"
	ImageMediaType  MediaType = "image"
)
