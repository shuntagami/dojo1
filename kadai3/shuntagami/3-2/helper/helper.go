package helper

import (
	"path"
	"path/filepath"
	"strings"
)

// FileNameFromURL returns a file name from download URL
func FileNameFromURL(url string) string {
	filename := path.Base(url)

	// remove query parameters if exists any
	index := strings.IndexRune(filename, '?')
	if index != -1 {
		filename = filename[:index]
	}

	return filename
}

// FileNameAndExt returns file name and it's extension
func FileNameAndExt(fileName string) (string, string) {
	ext := filepath.Ext(fileName)
	return strings.TrimSuffix(fileName, ext), ext
}
