package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gunsluo/tmplbuild"
	"github.com/gunsluo/tmplbuild/css"
	"github.com/gunsluo/tmplbuild/html"
	"github.com/gunsluo/tmplbuild/image"
	"github.com/gunsluo/tmplbuild/js"
	"github.com/gunsluo/tmplbuild/other"
)

func main() {
	var dst string
	var concurrent int
	var ignorePrefix string
	var replicaFiles string
	var basePath string

	flag.StringVar(&dst, "o", "dist", "is output dir")
	flag.IntVar(&concurrent, "c", 100, "is concurrent number")
	flag.StringVar(&ignorePrefix, "ignore-prefix", "", "ignore prefix for js css and image file")
	flag.StringVar(&replicaFiles, "replica-files", "", "the list of replica file")
	flag.StringVar(&basePath, "base-path", "", "the base path")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("missing source dir, for example: fsbuild -o dist src.")
		return
	}
	dir := args[0]
	repFiles := strings.Split(replicaFiles, ",")

	ts := newTasks()
	if err := ts.Walk(dir, basePath); err != nil {
		fmt.Printf("read: %v\n", err)
		return
	}

	ctx := &tmplbuild.Context{
		Dst:          dst,
		Dir:          dir,
		Concurrent:   concurrent,
		IgnorePrefix: ignorePrefix,
		ReplicaFiles: repFiles,
		BasePath:     basePath,
	}
	if err := ts.Build(ctx); err != nil {
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
			{tmplbuild.UnknowMediaType},
		},
	}

	return ts
}

func (ts *tasks) Walk(dir, basePath string) error {
	return filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if !strings.HasPrefix(path, basePath) {
				return nil
			}

			ts.AddFile(path)
			return nil
		})
}

func (ts *tasks) AddFile(path string) {
	mediaType := tmplbuild.GetMediaTypeByFilePath(path)

	t, ok := ts.all[mediaType]
	if !ok {
		compiler := newCompiler(mediaType)
		t = task{compiler: compiler}
	}

	t.files = append(t.files, path)
	ts.all[mediaType] = t
}

func (ts *tasks) Build(ctx *tmplbuild.Context) error {
	// run by priority
	symbols := tmplbuild.Symbols{}
	for _, mts := range ts.priority {
		tasks := []*task{}
		for _, mt := range mts {
			if t, ok := ts.all[mt]; ok {
				if len(t.files) > 0 {
					tasks = append(tasks, &t)
				}
			}
		}

		if len(tasks) == 0 {
			continue
		}

		// TODO: concurrent
		for _, t := range tasks {
			if err := t.Build(ctx, symbols); err != nil {
				return err
			}
		}
	}

	return nil
}

type task struct {
	compiler tmplbuild.Compiler
	files    []string
}

func (t *task) Build(ctx *tmplbuild.Context, symbols tmplbuild.Symbols) error {
	if t.compiler != nil {
		if err := t.compiler.Build(ctx, t.files, symbols); err != nil {
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
	case tmplbuild.UnknowMediaType:
		return &other.Compiler{}
	}

	return nil
}
