package xjson

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

var JSONT = jsoniter.Config{
	//关闭对HTML字符串编码
	EscapeHTML: false,
	//json为简单字符串即不进行unescape（do not unescape object field）
	ObjectFieldMustBeSimpleString: true,
	//是否开启Number承载数据（整形、浮点型）
	UseNumber: true,
}.Froze()

func ConvertToMap(data interface{}) (Map, error) {
	//fmt.Println("ConvertToMap data", data)
	if data == nil {
		return Map{}, nil
	}
	var dataJSON []byte
	switch data.(type) {
	case Map:
		return data.(Map), nil
	case string:
		//fmt.Println("string", data, len(data.(string)))
		if len(data.(string)) > 0 {
			dataJSON = []byte(data.(string))
		}
		break
	case []byte:
		dataJSON = data.([]byte)
		break
	default:
		jsonStr, err := JSONT.Marshal(&data)
		if err != nil {
			return nil, err
		}
		dataJSON = jsonStr
	}
	if len(dataJSON) == 0 {
		return nil, fmt.Errorf("data is empty")
	}

	var bodyMap Map
	err := JSONT.Unmarshal(dataJSON, &bodyMap)

	if err != nil {
		return nil, err
	}

	return bodyMap, nil
}
