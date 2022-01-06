package tmplbuild

type Compiler interface {
	Build(ctx *Context, files []string, placeholders Placeholders) error
}
