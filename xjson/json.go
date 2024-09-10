package xjson

import (
	"bytes"
	"encoding/json"
)

func ToJson(object interface{}) []byte {
	js, _ := json.Marshal(object)
	return js
}

func ToJsonString(object interface{}) string {
	js, _ := json.Marshal(object)
	return string(js)
}

func JSONPrettyFormat(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}

func FromJson(jsonStr string, object interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), object)
	return err
}
