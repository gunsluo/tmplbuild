package tmplbuild

type Context struct {
	Dst string
	Dir string

	Concurrent   int
	IgnorePrefix string
}
