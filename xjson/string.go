package xjson

import (
	"fmt"
	"strings"
)

type String string

func (f *String) String() string {
	if f == nil {
		return "0"
	}
	return fmt.Sprintf("%s", *f)
}

func (f *String) UnmarshalJSON(data []byte) error {
	floatStr := strings.Trim(string(data), "\"")

	if IsJsonNull(floatStr) || IsStrEmpty(floatStr) {
		*f = ""
		return nil
	}

	*f = String(floatStr)
	return nil
}
