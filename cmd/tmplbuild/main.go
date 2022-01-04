package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gunsluo/tmplbuild"
	"github.com/gunsluo/tmplbuild/css"
	"github.com/gunsluo/tmplbuild/html"
	"github.com/gunsluo/tmplbuild/image"
	"github.com/gunsluo/tmplbuild/js"
)

func main() {
	var dst string
	var ignorePrefix string

	flag.StringVar(&dst, "o", "dist", "is output dir")
	flag.StringVar(&ignorePrefix, "ignore-prefix", "", "ignore prefix for js css and image file")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("missing source dir, for example: fsbuild -o dist src.")
		return
	}
	dir := args[0]

	ts := newTasks()
	if err := ts.Read(dir); err != nil {
		fmt.Printf("read: %v\n", err)
		return
	}

	if err := ts.Build(dst, dir, ignorePrefix); err != nil {
		fmt.Printf("build: %v\n", err)
	}
}

type tasks struct {
	all      map[tmplbuild.MediaType]task
	priority [][]tmplbuild.MediaType
}

func newTasks() *tasks {
	ts := &tasks{
		all: make(map[tmplbuild.MediaType]task),
		priority: [][]tmplbuild.MediaType{
			{tmplbuild.ImageMediaType},
			{tmplbuild.JsMediaType, tmplbuild.CssMediaType},
			{tmplbuild.HtmlMediaType},
		},
	}

	return ts
}

func (ts *tasks) Read(dir string) error {
	return filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}

			ts.AddInput(&tmplbuild.Input{
				Path:   path,
				Reader: f,
			})

			return nil
		})
}

func (ts *tasks) AddInput(input *tmplbuild.Input) {
	mediaType := tmplbuild.GetMediaTypeByFilePath(input.Path)

	t, ok := ts.all[mediaType]
	if !ok {
		compiler := newCompiler(mediaType)
		t = task{compiler: compiler}
	}

	t.inputs = append(t.inputs, input)
	ts.all[mediaType] = t
}

func (ts *tasks) Build(dst, dir, ignorePrefix string) error {
	ctx := &tmplbuild.Context{
		Dst:          dst,
		Dir:          dir,
		IgnorePrefix: ignorePrefix,
	}

	// run by priority
	placeholders := tmplbuild.Placeholders{}
	for _, mts := range ts.priority {
		tasks := []*task{}
		for _, mt := range mts {
			if t, ok := ts.all[mt]; ok {
				if len(t.inputs) > 0 {
					tasks = append(tasks, &t)
				}
			}
		}

		if len(tasks) == 0 {
			continue
		}

		// TODO: concurrent
		for _, t := range tasks {
			if err := t.Build(ctx, placeholders); err != nil {
				return err
			}
		}
	}

	return nil
}

type task struct {
	compiler tmplbuild.Compiler
	inputs   []*tmplbuild.Input
}

func (t *task) Build(ctx *tmplbuild.Context, placeholders tmplbuild.Placeholders) error {
	if t.compiler != nil {
		if err := t.compiler.Build(ctx, t.inputs, placeholders); err != nil {
			return err
		}
	}
	return nil
}

func newCompiler(mediaType tmplbuild.MediaType) tmplbuild.Compiler {
	switch mediaType {
	case tmplbuild.HtmlMediaType:
		return &html.Compiler{}
	case tmplbuild.CssMediaType:
		return &css.Compiler{}
	case tmplbuild.JsMediaType:
		return &js.Compiler{}
	case tmplbuild.ImageMediaType:
		return &image.Compiler{}
	}

	return nil
}
