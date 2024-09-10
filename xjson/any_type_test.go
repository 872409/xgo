package xjson

import (
	"encoding/json"
	"testing"
)

type TestMyRequest struct {
	UserInput AnyType `protobuf:"bytes,4,rep,name=user_input,proto3,customtype=AnyType" json:"user_input,omitempty"`
}

func TestAnyType_MarshalJSON(t *testing.T) {
	j := `{"user_input":1234}`
	inst := &TestMyRequest{}
	err := json.Unmarshal([]byte(j), inst)
	if err != nil {
		t.Errorf("json decode error, err=%+v", err)
		return
	}
	t.Logf("%+v", inst)
	str, err := json.Marshal(inst)
	if err != nil {
		t.Errorf("json encode error, err=%+v", err)
		return
	}
	t.Logf("json=%s", string(str))
}
