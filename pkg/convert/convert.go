package convert

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Converter interface {
	Convert(relativePath, baseFolder, outFolder string) error
}

func NewConverter(name string) (Converter, error) {
	switch name {
	case "lame":
		return NewLameConverter(), nil
	case "opus":
		return NewOpusConverter(), nil
	}
	return nil, fmt.Errorf("unknown converter: %q", name)
}

func IsConvertible(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".flac")
}
