package xrequest

import (
	"fmt"
	"github.com/gogf/gf/v2/util/grand"
	"io"
	"net/http"
	"os"
	"path"
	"xgo/xio"
	"xgo/xutils"
)

func DownloadFileHash(url, dir string) (string, error) {
	return DownloadFile(url, dir, xutils.Md5Str(url))
}

func DownloadFileRand(url, dir string) (string, error) {
	var _fileName = grand.S(10) + xio.PathExt(url)
	return DownloadFileSave(url, path.Join(dir, _fileName))
}

func DownloadFile(url, dir, fileName string) (string, error) {
	var _fileName string
	if xio.PathExt(fileName) == "" {
		_fileName = fileName + xio.PathExt(url)
	}

	saveFilePath := fmt.Sprintf("%s/%s", dir, _fileName)
	return DownloadFileSave(url, saveFilePath)
}

func DownloadFileSave(url string, saveFilePath string) (string, error) {
	xio.MkdirAll(saveFilePath)
	v, err := http.Get(url)
	if err != nil {
		fmt.Printf("Http get [%v] failed! %v", url, err)
		return "", err
	}
	defer v.Body.Close()
	content, err := io.ReadAll(v.Body)
	if err != nil {
		fmt.Printf("Read http response failed! %v", err)
		return "", err
	}
	err = os.WriteFile(saveFilePath, content, 0666)
	if err != nil {
		fmt.Printf("Save to file failed! %v", err)
		return "", err
	}
	return saveFilePath, nil
}

func DownloadFileReader(url string) (io.ReadCloser, error) {
	v, err := http.Get(url)
	if err != nil {
		fmt.Printf("Http get [%v] failed! %v", url, err)
		return nil, err
	}
	return v.Body, nil
}
