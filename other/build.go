package other

import (
	"github.com/gunsluo/tmplbuild"
	"github.com/gunsluo/tmplbuild/core"
)

type Compiler struct {
	core.Compiler
}

func (b *Compiler) Build(ctx *tmplbuild.Context, files []string, placeholders tmplbuild.Placeholders) error {
	_, err := b.Compiler.Build(ctx, files, placeholders, b.build)
	if err != nil {
		return err
	}

	return nil
}

func (b *Compiler) build(ctx *tmplbuild.Context, input *tmplbuild.Input, placeholders tmplbuild.Placeholders) (string, string, error) {
	origin, target, err := b.Compiler.WriteNotChange(ctx, input.Path, input.Data)
	if err != nil {
		return "", "", err
	}

	return origin, target, nil
}
