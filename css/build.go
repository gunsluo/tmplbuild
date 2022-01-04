package css

import (
	"bytes"
	"io"

	"github.com/gunsluo/tmplbuild"
	"github.com/gunsluo/tmplbuild/core"
)

type Compiler struct {
	core.Compiler
}

func (b *Compiler) Build(ctx *tmplbuild.Context, inputs []*tmplbuild.Input, placeholders tmplbuild.Placeholders) error {
	placeholder, err := b.Compiler.Build(ctx, inputs, placeholders, b.build)
	if err != nil {
		return err
	}

	placeholders[tmplbuild.CssMediaType] = placeholder
	return nil
}

func (b *Compiler) build(ctx *tmplbuild.Context, input *tmplbuild.Input, placeholders tmplbuild.Placeholders) (string, string, error) {
	buffer, err := io.ReadAll(input.Reader)
	if err != nil {
		return "", "", err
	}

	data, err := b.insteadOfPlaceholder(buffer, placeholders)
	if err != nil {
		return "", "", err
	}

	origin, target, err := b.Compiler.Write(ctx, input.Path, data)
	if err != nil {
		return "", "", err
	}

	return origin, target, nil
}

func (b *Compiler) insteadOfPlaceholder(data []byte, placeholders tmplbuild.Placeholders) ([]byte, error) {
	placeholder, ok := placeholders[tmplbuild.ImageMediaType]
	if !ok {
		return data, nil
	}

	for o, t := range placeholder {
		data = bytes.ReplaceAll(data, []byte(o), []byte(t))
	}

	return data, nil
}
