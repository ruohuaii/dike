package dike

import (
	"reflect"
	"testing"
)

type TMatcher struct {
	Name string `dike:"re;size:1,32;bet:1,34;in:1,2,3,4;ni:5,6,7,8"`
	Age  uint8  `dike:"re;size:1,32;bet:1,37;in:1,2,3,4;ni:5,6,7,8"`
}

func Test_GetRelation(t *testing.T) {
	ft := reflect.TypeOf((*TMatcher)(nil)).Elem()
	matcher := NewMatcher(ft)
	result, err := matcher.GetDefined("dike")
	if err != nil {
		t.Fatal("error:", err)
	}
	t.Log("result:", *result["Age"], *result["Name"])
}
