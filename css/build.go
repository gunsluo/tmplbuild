package css

import (
	"bytes"

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

	symbols[tmplbuild.CssMediaType] = symbol
	return nil
}

func (b *Compiler) build(ctx *tmplbuild.Context, input *tmplbuild.Input, symbols tmplbuild.Symbols) (*tmplbuild.Output, error) {
	data, err := b.rewriteData(input.Data, input.Base, symbols)
	if err != nil {
		return nil, err
	}
	input.Data = data

	output, err := b.Compiler.Write(ctx, input, true)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (b *Compiler) rewriteData(data []byte, base string, symbols tmplbuild.Symbols) ([]byte, error) {
	symbol, ok := symbols[tmplbuild.ImageMediaType]
	if !ok {
		return data, nil
	}

	placeholders, ok := symbol[base]
	if !ok {
		return data, nil
	}

	// sort keys fixes a bug that causes replacement issue
	// when keys have the same suffix
	keys := placeholders.SortKeys()
	for _, o := range keys {
		t, ok := placeholders[o]
		if !ok {
			continue
		}
		data = bytes.ReplaceAll(data, []byte(o), []byte(t))
	}

	return data, nil
}
