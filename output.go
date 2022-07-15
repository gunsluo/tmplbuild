package tmplbuild

import "sort"

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

func (ps Placeholders) SortKeys() []string {
	keys := byLength{}

	for k := range ps {
		keys = append(keys, k)
	}

	sort.Sort(keys)
	return keys
}

type byLength []string

func (n byLength) Len() int           { return len(n) }
func (n byLength) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n byLength) Less(i, j int) bool { return len(n[i]) > len(n[j]) }
