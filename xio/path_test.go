package xio

import (
	"fmt"
	"path"
	"testing"
)

func TestPathExt(t *testing.T) {
	fmt.Println(path.Dir("/saveFilePath/aa.png"), path.Dir("/aa/saveFilePath/"), path.Dir("../saveFilePath/"))
}
