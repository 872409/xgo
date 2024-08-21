package xutils

import (
	"encoding/json"
	"net/http"
	"os"
)

func GetFileContentTypeWithPath(filePath string) (string, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return "", err
	}
	return GetFileContentType(file)
}

func GetFileContentType(output *os.File) (string, error) {

	// to sniff the content type only the first
	// 512 bytes are used.

	buf := make([]byte, 512)

	_, err := output.Read(buf)

	if err != nil {
		return "", err
	}

	// the function that actually does the trick
	contentType := http.DetectContentType(buf)

	return contentType, nil
}

func SaveFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0777)
}

func JSONMarshalToFile(path string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return SaveFile(path, data)
}

func JSONUnmarshalFromFile(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}
