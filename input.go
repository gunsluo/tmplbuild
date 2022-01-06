package tmplbuild

import (
	"io"
)

type Input struct {
	Path   string
	Reader io.Reader
}
