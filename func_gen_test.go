package dike

import "testing"

func Test_CreateFuncGenerate(t *testing.T) {
	input := &Input{
		StructName: "TMatcher",
	}
	t.Log(input.CreateFuncGenerate())
}
