package xio

import (
	"os"
	"path"
)

func MkdirAll(pathStr string, perm ...os.FileMode) error {
	_dir := path.Dir(pathStr)
	if _, err := os.Stat(_dir); os.IsNotExist(err) {
		_perm := os.FileMode(0777)
		if len(perm) > 0 {
			_perm = perm[0]
		}
		return os.MkdirAll(_dir, _perm)
	}
	return nil
}
