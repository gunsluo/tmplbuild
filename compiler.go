package tmplbuild

type Compiler interface {
	Build(ctx *Context, inputs []*Input, placeholders Placeholders) error
}
