package dike

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

var _ = fmt.Println

const gft = `package dike

import (
	"fmt"
	"reflect"
)

func Generator() error {
	t := reflect.TypeOf((*TMatcher)(nil)).Elem()
	matcher := NewMatcher(t)
	result, err := matcher.GetDefined("dike")
	if err != nil {
		return err
	}
	fmt.Println("result:", *result["Age"], *result["Name"])
	return nil
}`

type Input struct {
	StructName string
}

func (i *Input) CreateFuncGenerate() error {
	t, err := template.New("gft").Parse(gft)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, i)
	if err != nil {
		return err
	}
	filename := "generate.go"
	return os.WriteFile(filename, buf.Bytes(), 0666)
}
