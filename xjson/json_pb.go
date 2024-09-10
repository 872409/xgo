package xjson

import (
	//"github.com/golang/protobuf/jsonpb"
	//"github.com/golang/protobuf/proto"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"strings"
	//"google.golang.org/protobuf/proto"
)

func PbUnmarshal[P proto.Message](from interface{}, pb P) (P, error) {
	var jsonRAW string
	switch from.(type) {
	case string:
		jsonRAW = from.(string)
		break
	default:
		jsonRAW = ToJsonString(from)
	}

	err := jsonpb.UnmarshalString(jsonRAW, pb)
	if err != nil && strings.Contains(err.Error(), "unknown field") {
		return pb, nil
	}
	return pb, err
}
