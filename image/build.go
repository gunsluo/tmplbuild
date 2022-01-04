package image

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gunsluo/tmplbuild"
	"github.com/gunsluo/tmplbuild/hash"
)

type Compiler struct {
}

type req struct {
	ctx   *tmplbuild.Context
	input *tmplbuild.Input
}

type resp struct {
	placeholder tmplbuild.Placeholder
	err         error
}

type bOne struct {
	origin string
	target string
}

func (b *Compiler) Build(ctx *tmplbuild.Context, inputs []*tmplbuild.Input, placeholders tmplbuild.Placeholders) error {
	// create dst dir
	if err := os.MkdirAll(ctx.Dst, 0755); err != nil {
		return err
	}

	done := make(chan resp)
	ch := make(chan req)
	aggrch := make(chan bOne)

	go func() {
		placeholder := tmplbuild.Placeholder{}
		for v := range aggrch {
			placeholder[v.origin] = v.target
		}

		done <- resp{placeholder: placeholder}
	}()

	length := len(inputs)
	wg := sync.WaitGroup{}
	wg.Add(length)
	go func() {
		for r := range ch {
			go func(r req) {
				origin, target, err := b.buildOne(r.ctx, r.input)
				if err != nil {
					done <- resp{err: err}
					return
				}
				aggrch <- bOne{
					origin: origin,
					target: target,
				}
				wg.Done()
			}(r)
		}
		wg.Wait()
		close(aggrch)
	}()

	for _, input := range inputs {
		ch <- req{
			input: input,
			ctx: &tmplbuild.Context{
				Dst:          ctx.Dst,
				Dir:          ctx.Dir,
				IgnorePrefix: ctx.IgnorePrefix,
			},
		}
	}
	close(ch)

	// wait
	resp := <-done
	if resp.err == nil {
		placeholders[tmplbuild.ImageMediaType] = resp.placeholder
	}
	return resp.err
}

func (b *Compiler) buildOne(ctx *tmplbuild.Context, input *tmplbuild.Input) (string, string, error) {
	data, err := io.ReadAll(input.Reader)
	if err != nil {
		return "", "", err
	}

	relativePath := strings.TrimPrefix(strings.TrimPrefix(input.Path, ctx.Dir), "/")
	hashRelativePath := hash.GenerateName(relativePath, data)
	outPath := filepath.Join(ctx.Dst, hashRelativePath)

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return "", "", err
	}

	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return "", "", err
	}

	origin := relativePath
	target := hashRelativePath
	if ctx.IgnorePrefix != "" {
		origin = strings.TrimPrefix(strings.TrimPrefix(origin, ctx.IgnorePrefix), "/")
		target = strings.TrimPrefix(strings.TrimPrefix(target, ctx.IgnorePrefix), "/")
	}

	return origin, target, nil
}
