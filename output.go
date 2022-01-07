package tmplbuild

type Placeholders map[string]string

type Symbol map[string]Placeholders

type Symbols map[MediaType]Symbol

type Input struct {
	Path string
	Base string
	Data []byte

	RelativePath string
}

type Output struct {
	Base   string
	Origin string
	Target string
}
