package core

import (
	"os"
	"path/filepath"
	"strings"

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
	origin string
	target string
	err    error
}

type BuildFunc func(ctx *tmplbuild.Context, input *tmplbuild.Input, placeholders tmplbuild.Placeholders) (string, string, error)

func (b *Compiler) Build(ctx *tmplbuild.Context, inputs []*tmplbuild.Input, placeholders tmplbuild.Placeholders, buildFunc BuildFunc) (tmplbuild.Placeholder, error) {
	// create dst dir
	if err := os.MkdirAll(ctx.Dst, 0755); err != nil {
		return nil, err
	}

	length := len(inputs)
	var concurrent int
	if ctx.Concurrent >= length {
		concurrent = length
	} else {
		concurrent = ctx.Concurrent
	}
	ch := make(chan req)
	done := make(chan resp, length)

	// set max worker number
	for w := 0; w < concurrent; w++ {
		go func() {
			for r := range ch {
				origin, target, err := buildFunc(r.ctx, r.input, placeholders)
				done <- resp{
					origin: origin,
					target: target,
					err:    err,
				}
			}
		}()
	}

	for _, input := range inputs {
		ch <- req{
			input: input,
			ctx: &tmplbuild.Context{
				Dst:          ctx.Dst,
				Dir:          ctx.Dir,
				IgnorePrefix: ctx.IgnorePrefix,
				ReplicaFiles: ctx.ReplicaFiles,
			},
		}
	}
	close(ch)

	// wait
	var err error
	placeholder := tmplbuild.Placeholder{}
	for i := 0; i < length; i++ {
		resp, ok := <-done
		if !ok {
			break
		}
		if resp.err != nil {
			err = resp.err
			break
		}

		placeholder[resp.origin] = resp.target
	}

	return placeholder, err
}

func (b *Compiler) Write(ctx *tmplbuild.Context, path string, data []byte) (string, string, error) {
	relativePath := strings.TrimPrefix(strings.TrimPrefix(path, ctx.Dir), "/")
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

	// save replica file
	for _, f := range ctx.ReplicaFiles {
		if f == origin {
			replicaPath := filepath.Join(ctx.Dst, relativePath)
			if err := os.WriteFile(replicaPath, data, 0644); err != nil {
				return "", "", err
			}
			break
		}
	}

	return origin, target, nil
}

//func (b *Compiler) saveReplicaFiles(, path string, data []byte) (string, string, error) {
//}

func (b *Compiler) WriteNotChange(ctx *tmplbuild.Context, path string, data []byte) (string, string, error) {
	relativePath := strings.TrimPrefix(strings.TrimPrefix(path, ctx.Dir), "/")
	outPath := filepath.Join(ctx.Dst, relativePath)

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return "", "", err
	}

	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return "", "", err
	}

	origin := relativePath
	if ctx.IgnorePrefix != "" {
		origin = strings.TrimPrefix(strings.TrimPrefix(origin, ctx.IgnorePrefix), "/")
	}

	return origin, origin, nil
}
