package other

import (
	"io"

	"github.com/gunsluo/tmplbuild"
	"github.com/gunsluo/tmplbuild/core"
)

type Compiler struct {
	core.Compiler
}

func (b *Compiler) Build(ctx *tmplbuild.Context, inputs []*tmplbuild.Input, placeholders tmplbuild.Placeholders) error {
	_, err := b.Compiler.Build(ctx, inputs, placeholders, b.build)
	if err != nil {
		return err
	}

	return nil
}

func (b *Compiler) build(ctx *tmplbuild.Context, input *tmplbuild.Input, placeholders tmplbuild.Placeholders) (string, string, error) {
	data, err := io.ReadAll(input.Reader)
	if err != nil {
		return "", "", err
	}

	origin, target, err := b.Compiler.WriteNotChange(ctx, input.Path, data)
	if err != nil {
		return "", "", err
	}

	return origin, target, nil
}
