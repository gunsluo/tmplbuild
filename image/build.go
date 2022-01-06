package image

import (
	"github.com/gunsluo/tmplbuild"
	"github.com/gunsluo/tmplbuild/core"
)

type Compiler struct {
	core.Compiler
}

func (b *Compiler) Build(ctx *tmplbuild.Context, files []string, placeholders tmplbuild.Placeholders) error {
	placeholder, err := b.Compiler.Build(ctx, files, placeholders, b.build)
	if err != nil {
		return err
	}

	placeholders[tmplbuild.ImageMediaType] = placeholder
	return nil
}

func (b *Compiler) build(ctx *tmplbuild.Context, input *tmplbuild.Input, placeholders tmplbuild.Placeholders) (string, string, error) {
	origin, target, err := b.Compiler.Write(ctx, input.Path, input.Data)
	if err != nil {
		return "", "", err
	}

	return origin, target, nil
}
