package xjson

import (
	"strconv"
	"strings"
)

type Int64 int64

func (f *Int64) UnmarshalJSON(data []byte) error {
	floatStr := strings.Trim(string(data), "\"")

	if IsJsonNull(floatStr) || IsStrEmpty(floatStr) {
		*f = Int64(0)
		return nil
	}

	v, err := strconv.ParseFloat(floatStr, 64)
	//fmt.Println(v, err)
	if err == nil {
		*f = Int64(v)
		return nil
	}
	*f = Int64(0)
	return nil
}
