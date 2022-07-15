package html

import (
	"bytes"

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
	data, err := b.rewriteData(input.Data, input.Base, symbols)
	if err != nil {
		return nil, err
	}
	input.Data = data

	output, err := b.Compiler.Write(ctx, input, false)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (b *Compiler) rewriteData(data []byte, base string, symbols tmplbuild.Symbols) ([]byte, error) {
	placeholders := tmplbuild.Placeholders{}
	for _, symbol := range symbols {
		ps, ok := symbol[base]
		if !ok {
			continue
		}

		for k, v := range ps {
			placeholders[k] = v
		}
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
