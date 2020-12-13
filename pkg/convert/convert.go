package convert

import "path/filepath"

type Converter interface {
	Convert(relativePath, baseFolder, outFolder string) error
}

func IsConvertible(path string) bool {
	return filepath.Ext(path) == ".flac"
}
