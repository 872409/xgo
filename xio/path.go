package xio

import (
	"path"
	"strings"
)

func PathExt(pathStr string) string {
	ext := path.Ext(pathStr)
	if i := strings.Index(ext, "?"); i > -1 {
		ext = ext[0:i]
	}
	return ext
}
