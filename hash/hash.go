package hash

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"strings"
)

func GenerateName(filename string, data []byte) string {
	sh := calcShortHash(data, 16)

	ext := filepath.Ext(filename)
	base := filepath.Base(filename)
	dir := filepath.Dir(filename)
	name := strings.TrimSuffix(base, ext)
	name += "-" + sh + ext

	return filepath.Join(dir, name)
}

func calcShortHash(data []byte, length int) string {
	hash := CalcHash(data)
	if len(hash) <= length {
		return hash
	}

	return hash[len(hash)-length:]
}

func CalcHash(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))
}
