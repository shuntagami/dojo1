package helper

import (
	"io/fs"
	"path/filepath"

	"github.com/pkg/errors"
)

/*
Fullpath returns a full path of specific directory under sample.
@example /workspace/sample/sample2
*/
func Fullpath(dirname string) (fullpath string, err error) {
	err = filepath.WalkDir(filepath.Join("/workspace", "sample"), func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed filepath.WalkDir")
		}
		if info.IsDir() && filepath.Base(path) == dirname {
			fullpath = path
			return err
		}
		return err
	})
	return fullpath, err
}
