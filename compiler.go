package tmplbuild

type Compiler interface {
	Build(ctx *Context, files []string, symbols Symbols) error
}
