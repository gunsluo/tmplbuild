package image

import (
	"github.com/gunsluo/tmplbuild"
	"github.com/gunsluo/tmplbuild/core"
)

type Compiler struct {
	core.Compiler
}

func (b *Compiler) Build(ctx *tmplbuild.Context, files []string, symbols tmplbuild.Symbols) error {
	symbol, err := b.Compiler.Build(ctx, files, symbols, b.build)
	if err != nil {
		return err
	}

	symbols[tmplbuild.ImageMediaType] = symbol
	return nil
}

func (b *Compiler) build(ctx *tmplbuild.Context, input *tmplbuild.Input, symbols tmplbuild.Symbols) (*tmplbuild.Output, error) {
	output, err := b.Compiler.Write(ctx, input, true)
	if err != nil {
		return nil, err
	}

	return output, nil
}
