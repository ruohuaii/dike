package dike

import (
	"reflect"
	"testing"
)

func Test_Builder(t *testing.T) {
	ft := reflect.TypeOf((*TBuilder)(nil)).Elem()
	matcher := NewMatcher(ft)
	result, err := matcher.GetDefined("dike")
	if err != nil {
		t.Fatal("error:", err)
	}
	name, short := matcher.GetStructName()
	builder := NewBuilder(name, short, result)
	err = builder.Build("1.go")
	t.Log("error:", err)
}
