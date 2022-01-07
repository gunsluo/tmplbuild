package other

import (
	"github.com/gunsluo/tmplbuild"
	"github.com/gunsluo/tmplbuild/core"
)

type Compiler struct {
	core.Compiler
}

func (b *Compiler) Build(ctx *tmplbuild.Context, files []string, symbols tmplbuild.Symbols) error {
	_, err := b.Compiler.Build(ctx, files, symbols, b.build)
	if err != nil {
		return err
	}

	return nil
}

func (b *Compiler) build(ctx *tmplbuild.Context, input *tmplbuild.Input, symbols tmplbuild.Symbols) (*tmplbuild.Output, error) {
	output, err := b.Compiler.Write(ctx, input, false)
	if err != nil {
		return nil, err
	}

	return output, nil
}
