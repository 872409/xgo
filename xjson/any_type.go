package xjson

import (
	"encoding/json"
	"errors"
)

type AnyType struct {
	Value interface{}
}

func NewAny(value any) AnyType {
	return AnyType{
		Value: value,
	}
}

func (t AnyType) Marshal() ([]byte, error) {
	return nil, errors.New("not implement")
}
func (t *AnyType) MarshalTo(data []byte) (n int, err error) {
	return 0, errors.New("not implement")
}
func (t *AnyType) Unmarshal(data []byte) error {
	return errors.New("not implement")
}
func (t *AnyType) Size() int {
	return -1
}

func (t AnyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}
func (t *AnyType) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Value)
}
