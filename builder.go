package dike

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Builder struct {
	name  string
	short string
	d     map[string]*Rule
}

func NewBuilder(name, short string, d map[string]*Rule) *Builder {
	return &Builder{
		name:  name,
		short: short,
		d:     d,
	}
}

func (b *Builder) Build(filename string) error {
	var writer bytes.Buffer
	writer.WriteString(fmt.Sprintf("func (%s *%s) Check() error {\n", b.short, b.name))
	for k, v := range b.d {
		dc := v.Dc
		if dc == "" {
			dc = k
		}
		if v.Optional {
			str, err := optionalSolution(v.st, v.Kind, b.short, k)
			if err != nil {
				return err
			}
			writer.WriteString(str)
			if v.Kind == "struct" {
				err := subStructSolution(v.st.Type, v.sv, filename, b.short, k, &writer)
				if err != nil {
					return err
				}
				continue
			}

			if v.Kind == "ptr" && v.sv.Type().Elem().Kind() == reflect.Struct {
				err := subStructSolution(v.sv.Type().Elem(), v.sv.Elem(), filename, b.short, k, &writer)
				if err != nil {
					return err
				}
				continue
			}

			if v.Len != "" {
				str, err := lenSolution(v.Kind, b.short, k, dc, v.Len)
				if err != nil {
					return err
				}
				writer.WriteString(str)
				v.Size = nil
			}
			if len(v.Size) == 2 {
				str, err := sizeSolution(v.Kind, b.short, k, dc, v.Size)
				if err != nil {
					return err
				}
				writer.WriteString(str)
			}
			if len(v.Between) == 2 {
				str, err := betweenSolution(v.Kind, b.short, k, dc, v.Between)
				if err != nil {
					return err
				}
				writer.WriteString(str)
			}
			if v.Eq != "" {
				str, err := eqSolution(v.Kind, b.short, k, dc, v.Eq)
				if err != nil {
					return err
				}
				writer.WriteString(str)
				v.Neq = ""
				v.Lt = ""
				v.Lte = ""
				v.Gt = ""
				v.Gte = ""
			}
			if v.Neq != "" {
				str, err := neqSolution(v.Kind, b.short, k, dc, v.Neq)
				if err != nil {
					return err
				}
				writer.WriteString(str)
				v.Lt = ""
				v.Lte = ""
				v.Gt = ""
				v.Gte = ""
			}
			if v.Lt != "" {
				str, err := ltSolution(v.Kind, b.short, k, dc, v.Lt)
				if err != nil {
					return err
				}
				writer.WriteString(str)
				v.Lte = ""
				v.Gt = ""
				v.Gte = ""
			}
			if v.Lte != "" {
				str, err := lteSolution(v.Kind, b.short, k, dc, v.Lte)
				if err != nil {
					return err
				}
				writer.WriteString(str)
				v.Gt = ""
				v.Gte = ""
			}
			if v.Gt != "" {
				str, err := gtSolution(v.Kind, b.short, k, dc, v.Lte)
				if err != nil {
					return err
				}
				writer.WriteString(str)
				v.Gte = ""
			}
			if v.Gte != "" {
				str, err := gteSolution(v.Kind, b.short, k, dc, v.Lte)
				if err != nil {
					return err
				}
				writer.WriteString(str)
			}
			if v.Regexp != "" {
				str, err := regexpSolution(v.Kind, b.short, k, dc, v.Regexp)
				if err != nil {
					return err
				}
				writer.WriteString(str)
			}
			if len(v.In) > 0 {
				str, err := inSolution(v.Kind, b.short, k, dc, v.In)
				if err != nil {
					return err
				}
				writer.WriteString(str)
				v.NotIn = nil
			}
			if len(v.NotIn) > 0 {
				str, err := niSolution(v.Kind, b.short, k, dc, v.NotIn)
				if err != nil {
					return err
				}
				writer.WriteString(str)
			}
			writer.WriteString("}")
			continue
		}
		if v.Required {
			str, err := requiredSolution(v.st, v.Kind, b.short, k, dc)
			if err != nil {
				return err
			}
			writer.WriteString(str)
		}
		if v.Len != "" {
			str, err := lenSolution(v.Kind, b.short, k, dc, v.Len)
			if err != nil {
				return err
			}
			writer.WriteString(str)
			v.Size = nil
		}
		if len(v.Size) == 2 {
			str, err := sizeSolution(v.Kind, b.short, k, dc, v.Size)
			if err != nil {
				return err
			}
			writer.WriteString(str)
		}
		if len(v.Between) == 2 {
			str, err := betweenSolution(v.Kind, b.short, k, dc, v.Between)
			if err != nil {
				return err
			}
			writer.WriteString(str)
		}
		if v.Eq != "" {
			str, err := eqSolution(v.Kind, b.short, k, dc, v.Eq)
			if err != nil {
				return err
			}
			writer.WriteString(str)
			v.Neq = ""
			v.Lt = ""
			v.Lte = ""
			v.Gt = ""
			v.Gte = ""
		}
		if v.Neq != "" {
			str, err := neqSolution(v.Kind, b.short, k, dc, v.Neq)
			if err != nil {
				return err
			}
			writer.WriteString(str)
			v.Lt = ""
			v.Lte = ""
			v.Gt = ""
			v.Gte = ""
		}
		if v.Lt != "" {
			str, err := ltSolution(v.Kind, b.short, k, dc, v.Lt)
			if err != nil {
				return err
			}
			writer.WriteString(str)
			v.Lte = ""
			v.Gt = ""
			v.Gte = ""
		}
		if v.Lte != "" {
			str, err := lteSolution(v.Kind, b.short, k, dc, v.Lte)
			if err != nil {
				return err
			}
			writer.WriteString(str)
			v.Gt = ""
			v.Gte = ""
		}
		if v.Gt != "" {
			str, err := gtSolution(v.Kind, b.short, k, dc, v.Lte)
			if err != nil {
				return err
			}
			writer.WriteString(str)
			v.Gte = ""
		}
		if v.Gte != "" {
			str, err := gteSolution(v.Kind, b.short, k, dc, v.Lte)
			if err != nil {
				return err
			}
			writer.WriteString(str)
		}
		if v.Regexp != "" {
			str, err := regexpSolution(v.Kind, b.short, k, dc, v.Regexp)
			if err != nil {
				return err
			}
			writer.WriteString(str)
		}
		if len(v.In) > 0 {
			str, err := inSolution(v.Kind, b.short, k, dc, v.In)
			if err != nil {
				return err
			}
			writer.WriteString(str)
			v.NotIn = nil
		}
		if len(v.NotIn) > 0 {
			str, err := niSolution(v.Kind, b.short, k, dc, v.NotIn)
			if err != nil {
				return err
			}
			writer.WriteString(str)
		}
	}
	writer.WriteString(`return nil
}`)
	f, err := os.OpenFile(filename, os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	_, err = f.Write(writer.Bytes())
	_ = f.Close()
	return err
}

func requiredSolution(
	st reflect.StructField, kind string, shortName, fieldName, dc string) (string, error) {
	var str string
	switch kind {
	case "string":
		format := `if %s.%s ==""{
	return errors.New("%s can't be empty")
}`
		str = fmt.Sprintf(format, shortName, fieldName, dc)
	case "int8", "uint8",
		"int16", "uint16",
		"int32", "uint32",
		"int", "uint",
		"int64", "uint64",
		"float32", "float64":
		format := `if %s.%s ==0{
	return errors.New("%s can't be empty")
}`
		str = fmt.Sprintf(format, shortName, fieldName, dc)
	case "ptr":
		format := `if %s.%s ==nil{
	return errors.New("%s can't be empty")
}`

		str = fmt.Sprintf(format, shortName, fieldName, dc)
	case "slice", "array":
		format := `if len(%s.%s) ==0{
	return errors.New("%s can't be empty")
}`
		str = fmt.Sprintf(format, shortName, fieldName, dc)
	case "struct":
		format := `if reflect.DeepEqual(%s.%s,%s{}){
	return errors.New("%s can't be empty")
}`

		str = fmt.Sprintf(format, shortName, fieldName, st.Type.Name(), dc)
	default:
		return "", fmt.Errorf("%s kind unsupported", kind)
	}
	return str, nil
}

func optionalSolution(
	st reflect.StructField, kind string, shortName, fieldName string) (string, error) {
	var str string
	switch kind {
	case "string":
		format := `if %s.%s !=""{`
		str = fmt.Sprintf(format, shortName, fieldName)
	case "int8", "uint8",
		"int16", "uint16",
		"int32", "uint32",
		"int", "uint",
		"int64", "uint64",
		"float32", "float64":
		format := `if %s.%s !=0{`
		str = fmt.Sprintf(format, shortName, fieldName)
	case "ptr":
		format := `if %s.%s !=nil{`
		str = fmt.Sprintf(format, shortName, fieldName)
	case "slice", "array":
		format := `if len(%s.%s)!=0{`
		str = fmt.Sprintf(format, shortName, fieldName)
	case "struct":
		format := `if !reflect.DeepEqual(%s.%s,%s{}){`
		str = fmt.Sprintf(format, shortName, fieldName, st.Type.Name())
	default:
		return "", fmt.Errorf("%s kind unsupported", kind)
	}
	return str, nil
}

func lenSolution(
	kind string, shortName, fieldName, dc, length string) (string, error) {
	switch kind {
	case "int8", "uint8",
		"int16", "uint16",
		"int32", "uint32",
		"int", "uint",
		"int64", "uint64",
		"float32", "float64",
		"ptr", "struct":
		return "", fmt.Errorf("%s can't use the len condition", dc)
	default:
		format := `if len(%s.%s)!=%s{
	return errors.New("%s must have length %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, length, dc, length)
		return str, nil
	}
}

func sizeSolution(
	kind string, shortName, fieldName, dc string, size []string) (string, error) {
	switch kind {
	case "int8", "uint8",
		"int16", "uint16",
		"int32", "uint32",
		"int", "uint",
		"int64", "uint64",
		"float32", "float64",
		"ptr", "struct":
		return "", fmt.Errorf("%s can't use the size condition", dc)
	default:
		format := `if len(%s.%s)<%s || len(%s.%s)>%s{
	return errors.New("the length of %s must be between %s and %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, size[0], shortName, fieldName, size[1], dc, size[0], size[1])
		return str, nil
	}
}

func betweenSolution(
	kind string, shortName, fieldName, dc string, between []string) (string, error) {
	switch kind {
	case "string":
		format := `if %s.%s<"%s" || %s.%s>"%s"{
	return errors.New("the value of %s must be between %s and %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, between[0], shortName, fieldName, between[1], dc, between[0], between[1])
		return str, nil
	case "int8", "uint8",
		"int16", "uint16",
		"int32", "uint32",
		"int", "uint",
		"int64", "uint64",
		"float32", "float64":
		format := `if %s.%s<%s || %s.%s>%s{
	return errors.New("the value of %s must be between %s and %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, between[0], shortName, fieldName, between[1], dc, between[0], between[1])
		return str, nil
	default:
		return "", fmt.Errorf("%s can't use the bet condition", dc)
	}
}

func eqSolution(kind string, shortName, fieldName, dc, eq string) (string, error) {
	switch kind {
	case "string":
		format := `if %s.%s!="%s" {
	return errors.New("the value of %s must be %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, eq, dc, eq)
		return str, nil
	case "int8", "uint8",
		"int16", "uint16",
		"int32", "uint32",
		"int", "uint",
		"int64", "uint64",
		"float32", "float64":
		format := `if %s.%s!=%s {
	return errors.New("the value of %s must be %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, eq, dc, eq)
		return str, nil
	default:
		return "", fmt.Errorf("%s can't use the eq condition", dc)
	}
}

func neqSolution(kind string, shortName, fieldName, dc, neq string) (string, error) {
	switch kind {
	case "string":
		format := `if %s.%s=="%s" {
	return errors.New("the value of %s can't be %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, neq, dc, neq)
		return str, nil
	case "int8", "uint8",
		"int16", "uint16",
		"int32", "uint32",
		"int", "uint",
		"int64", "uint64",
		"float32", "float64":
		format := `if %s.%s==%s {
	return errors.New("the value of %s can't be %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, neq, dc, neq)
		return str, nil
	default:
		return "", fmt.Errorf("%s can't use the neq condition", dc)
	}
}

func ltSolution(kind string, shortName, fieldName, dc, lt string) (string, error) {
	switch kind {
	case "string":
		format := `if %s.%s>="%s" {
	return errors.New("the value of %s must be less than %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, lt, dc, lt)
		return str, nil
	case "int8", "uint8",
		"int16", "uint16",
		"int32", "uint32",
		"int", "uint",
		"int64", "uint64",
		"float32", "float64":
		format := `if %s.%s>=%s {
	return errors.New("the value of %s must be less than %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, lt, dc, lt)
		return str, nil
	default:
		return "", fmt.Errorf("%s can't use the lt condition", dc)
	}
}

func lteSolution(kind string, shortName, fieldName, dc, lte string) (string, error) {
	switch kind {
	case "string":
		format := `if %s.%s>"%s" {
	return errors.New("the value of %s must be less than or equal to %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, lte, dc, lte)
		return str, nil
	case "int8", "uint8",
		"int16", "uint16",
		"int32", "uint32",
		"int", "uint",
		"int64", "uint64",
		"float32", "float64":
		format := `if %s.%s>=%s {
	return errors.New("the value of %s must be less than or equal to %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, lte, dc, lte)
		return str, nil
	default:
		return "", fmt.Errorf("%s can't use the lte condition", dc)
	}
}

func gtSolution(kind string, shortName, fieldName, dc, gt string) (string, error) {
	switch kind {
	case "string":
		str := fmt.Sprintf(`if %s.%s<="%s" {
	return errors.New("the value of %s must be great than %s")
}`, shortName, fieldName, gt, dc, gt)
		return str, nil
	case "int8", "uint8",
		"int16", "uint16",
		"int32", "uint32",
		"int", "uint",
		"int64", "uint64",
		"float32", "float64":
		str := fmt.Sprintf(`if %s.%s<=%s {
	return errors.New("the value of %s must be great than %s")
}`, shortName, fieldName, gt, dc, gt)
		return str, nil
	default:
		return "", fmt.Errorf("%s can't use the lte condition", dc)
	}
}

func gteSolution(kind string, shortName, fieldName, dc, gte string) (string, error) {
	switch kind {
	case "string":
		format := `if %s.%s<"%s" {
	return errors.New("the value of %s must be great than or equal to %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, gte, dc, gte)
		return str, nil
	case "int8", "uint8",
		"int16", "uint16",
		"int32", "uint32",
		"int", "uint",
		"int64", "uint64",
		"float32", "float64":
		format := `if %s.%s<%s {
	return errors.New("the value of %s must be great than or equal to %s")
}`
		str := fmt.Sprintf(format, shortName, fieldName, gte, dc, gte)
		return str, nil
	default:
		return "", fmt.Errorf("%s can't use the gte condition", dc)
	}
}

func regexpSolution(
	kind string, shortName, fieldName, dc, reg string) (string, error) {
	switch kind {
	case "string":
		format := `if ok,_:= regexp.MatchString("%s",%s.%s);!ok{
	return errors.New("%s does not match the regular")
}`
		str := fmt.Sprintf(format, reg, shortName, fieldName, dc)
		return str, nil
	default:
		return "", fmt.Errorf("%s can't use the regexp condition", dc)
	}
}

func inSolution(
	kind string, shortName, fieldName, dc string, in []string) (string, error) {
	var inw strings.Builder
	switch kind {
	case "string":
		inw.WriteString("[]string{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`"%s"`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`"%s",`, v))
			}
		}
	case "int8":
		inw.WriteString("[]int8{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "uint8":
		inw.WriteString("[]uint8{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "int16":
		inw.WriteString("[]int16{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "uint16":
		inw.WriteString("[]uint16{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "int32":
		inw.WriteString("[]int32{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "uint32":
		inw.WriteString("[]uint32{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "int":
		inw.WriteString("[]int{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "uint":
		inw.WriteString("[]uint{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "int64":
		inw.WriteString("[]int64{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "uint64":
		inw.WriteString("[]uint64{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "float32":
		inw.WriteString("[]float32{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "float64":
		inw.WriteString("[]float64{")
		for k, v := range in {
			if k == len(in)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	default:
		return "", fmt.Errorf("%s can't use the in condition", dc)
	}
	inw.WriteString("}")
	ina := inw.String()
	format := `var %sInColl=%s
var is%sIn bool
for _,v:=range %sInColl{
	if v==%s.%s{
		is%sIn=true
	}
}
if !is%sIn {
	return errors.New("%s must be in the specified set")
}`
	str := fmt.Sprintf(
		format, dc, ina, fieldName, dc, shortName, fieldName, fieldName, fieldName, dc)
	return str, nil
}

func niSolution(
	kind string, shortName, fieldName, dc string, ni []string) (string, error) {
	var inw strings.Builder
	switch kind {
	case "string":
		inw.WriteString("[]string{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`"%s"`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`"%s",`, v))
			}
		}
	case "int8":
		inw.WriteString("[]int8{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "uint8":
		inw.WriteString("[]uint8{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "int16":
		inw.WriteString("[]int16{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "uint16":
		inw.WriteString("[]uint16{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "int32":
		inw.WriteString("[]int32{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "uint32":
		inw.WriteString("[]uint32{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "int":
		inw.WriteString("[]int{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "uint":
		inw.WriteString("[]uint{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "int64":
		inw.WriteString("[]int64{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "uint64":
		inw.WriteString("[]uint64{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "float32":
		inw.WriteString("[]float32{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	case "float64":
		inw.WriteString("[]float64{")
		for k, v := range ni {
			if k == len(ni)-1 {
				inw.WriteString(fmt.Sprintf(`%s`, v))
			} else {
				inw.WriteString(fmt.Sprintf(`%s,`, v))
			}
		}
	default:
		return "", fmt.Errorf("%s can't use the in condition", dc)
	}
	inw.WriteString("}")
	ina := inw.String()
	format := `var %sNiColl=%s
for _,v:=range %sNiColl{
	if v==%s.%s{
		return errors.New("%s can't be in the specified set")
	}
}`
	str := fmt.Sprintf(format, dc, ina, dc, shortName, fieldName, dc)
	return str, nil
}

func subStructSolution(
	st reflect.Type, sv reflect.Value, filename, shortName, fieldName string,
	writer *bytes.Buffer) error {
	sm := NewMatcherWithReflect(st, sv)
	result, err := sm.GetDefined("dike")
	if err != nil {
		return err
	}
	name, short := sm.GetStructName()
	sb := NewBuilder(name, short, result)
	err = sb.Build(filename)
	if err != nil {
		return err
	}
	format := `if err := %s.%s.Check();err !=nil {
	return err
}`
	str := fmt.Sprintf(format, shortName, fieldName)
	writer.WriteString(str)
	writer.WriteString("}")
	return nil
}
