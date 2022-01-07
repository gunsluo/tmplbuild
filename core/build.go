package core

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gunsluo/tmplbuild"
	"github.com/gunsluo/tmplbuild/hash"
)

type Compiler struct {
}

type req struct {
	path string
	ctx  *tmplbuild.Context
}

type resp struct {
	output *tmplbuild.Output
	err    error
}

type BuildFunc func(ctx *tmplbuild.Context, input *tmplbuild.Input, symbols tmplbuild.Symbols) (*tmplbuild.Output, error)

func (b *Compiler) Build(ctx *tmplbuild.Context, files []string, symbols tmplbuild.Symbols, buildFunc BuildFunc) (tmplbuild.Symbol, error) {
	// create dst dir
	if err := os.MkdirAll(ctx.Dst, 0755); err != nil {
		return nil, err
	}

	length := len(files)
	var concurrent int
	if ctx.Concurrent >= length {
		concurrent = length
	} else {
		concurrent = ctx.Concurrent
	}
	ch := make(chan req)
	done := make(chan resp, length)

	// set read & write worker number
	for w := 0; w < concurrent; w++ {
		go func() {
			for r := range ch {
				input, err := b.Read(r.ctx, r.path)
				if err != nil {
					done <- resp{
						err: err,
					}
					return
				}
				output, err := buildFunc(r.ctx, input, symbols)
				done <- resp{
					output: output,
					err:    err,
				}
			}
		}()
	}

	for _, path := range files {
		ch <- req{
			path: path,
			ctx: &tmplbuild.Context{
				Dst:          ctx.Dst,
				Dir:          ctx.Dir,
				IgnorePrefix: ctx.IgnorePrefix,
				ReplicaFiles: ctx.ReplicaFiles,
				BasePath:     ctx.BasePath,
			},
		}
	}
	close(ch)

	// wait
	var err error
	symbol := tmplbuild.Symbol{}
	for i := 0; i < length; i++ {
		resp, ok := <-done
		if !ok {
			break
		}
		if resp.err != nil {
			err = resp.err
			break
		}

		placeholders, ok := symbol[resp.output.Base]
		if !ok {
			placeholders = tmplbuild.Placeholders{}
		}
		placeholders[resp.output.Origin] = resp.output.Target
		symbol[resp.output.Base] = placeholders
	}

	return symbol, err
}

func (b *Compiler) Read(ctx *tmplbuild.Context, path string) (*tmplbuild.Input, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var relativePath, base string
	if ctx.BasePath != "" {
		relativePath = strings.TrimPrefix(strings.TrimPrefix(path, ctx.BasePath), "/")
		if idx := strings.Index(relativePath, "/"); idx > 0 {
			base = relativePath[:idx]
			relativePath = relativePath[idx+1:]
		}
	} else {
		relativePath = strings.TrimPrefix(strings.TrimPrefix(path, ctx.Dir), "/")
	}

	input := &tmplbuild.Input{
		Base: base,
		Path: path,
		Data: data,

		RelativePath: relativePath,
	}

	return input, nil
}

func (b *Compiler) Write(ctx *tmplbuild.Context, input *tmplbuild.Input, enableHash bool) (*tmplbuild.Output, error) {
	var outPath, outDir, origin, target string

	if input.Base == "" {
		outDir = ctx.Dst
	} else {
		outDir = filepath.Join(strings.Replace(ctx.BasePath, ctx.Dir, ctx.Dst, 1), input.Base)
	}

	if enableHash {
		hashRelativePath := hash.GenerateName(input.RelativePath, input.Data)
		outPath = filepath.Join(outDir, hashRelativePath)

		origin = input.RelativePath
		target = hashRelativePath
		if ctx.IgnorePrefix != "" {
			origin = strings.TrimPrefix(strings.TrimPrefix(origin, ctx.IgnorePrefix), "/")
			target = strings.TrimPrefix(strings.TrimPrefix(target, ctx.IgnorePrefix), "/")
		}
	} else {
		outPath = filepath.Join(outDir, input.RelativePath)

		origin = input.RelativePath
		target = input.RelativePath
		if ctx.IgnorePrefix != "" {
			origin = strings.TrimPrefix(strings.TrimPrefix(origin, ctx.IgnorePrefix), "/")
		}
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return nil, err
	}

	if err := os.WriteFile(outPath, input.Data, 0644); err != nil {
		return nil, err
	}

	if enableHash {
		// replica file
		for _, f := range ctx.ReplicaFiles {
			if f == origin {
				replicaPath := filepath.Join(outDir, input.RelativePath)
				if err := os.WriteFile(replicaPath, input.Data, 0644); err != nil {
					return nil, err
				}
				break
			}
		}
	}

	return &tmplbuild.Output{
		Base:   input.Base,
		Origin: origin,
		Target: target,
	}, nil
}
