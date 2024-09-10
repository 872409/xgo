package xjson

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Person struct {
	Name string
	Util Float64
}

func TestFloat64_MarshalJSON(t *testing.T) {
	jsonStr := `{"Name":"aaaaa0","Util":988.098801}`
	var o Person
	json.Unmarshal([]byte(jsonStr), &o)
	fmt.Println(o, o.Util)
}

func TestFloat64_12(t *testing.T) {
	var agentIds []int64
	FromJson("[1,2,3]", &agentIds)
	fmt.Println(agentIds)
}
