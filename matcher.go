package dike

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

type Matcher struct {
	t reflect.Type
	v reflect.Value
}

type Rule struct {
	//Kind 字段类型
	Kind string
	//Required 必须包含
	Required bool
	//Optional 非必填
	Optional bool
	Eq       string
	Neq      string
	Lt       string
	Lte      string
	Gt       string
	Gte      string
	Between  []string
	In       []string
	NotIn    []string
	Size     []string
	Len      string
	Regexp   string
	//Dc 字段描述,例如json序列化时将其指定为dc
	Dc string

	st reflect.StructField
	sv reflect.Value
}

func NewMatcher(ptr interface{}) *Matcher {
	if ptr == nil {
		panic("NewMatcher received nil")
	}
	t := reflect.TypeOf(ptr)
	v := reflect.ValueOf(ptr)
	m := &Matcher{
		t: t,
		v: v,
	}
	err := m.check()
	if err != nil {
		panic(err)
	}
	m.t = m.t.Elem()
	m.v = m.v.Elem()
	return m
}

func NewMatcherWithReflect(t reflect.Type, v reflect.Value) *Matcher {
	return &Matcher{
		t: t,
		v: v,
	}
}

func (m *Matcher) GetStructName() (string, string) {
	name := m.t.Name()
	short := strings.ToLower(string(name[0]))
	return name, short
}

//GetDefined 获取定义的tag
func (m *Matcher) GetDefined(tag string) (map[string]*Rule, error) {
	rules := make(map[string]*Rule)
	for i := 0; i < m.t.NumField(); i++ {
		if !m.v.IsValid() {
			return nil, errors.New("invalid reflected value or field is not accessible")
		}
		st := m.t.Field(i)
		defined := st.Tag.Get(tag)
		sv := m.v.Field(i)
		kind := sv.Type().Kind().String()
		rule, err := m.analyze(kind, defined)
		if err != nil {
			return nil, err
		}
		rule.st = st
		rule.sv = sv
		rules[st.Name] = rule
	}
	return rules, nil
}

//analyze 解析tag信息
func (m *Matcher) analyze(kind, defined string) (*Rule, error) {
	var rule = &Rule{
		Kind: kind,
	}
	dvs := strings.Split(defined, ";")
	for _, v := range dvs {
		if len(v) == 0 {
			continue
		}
		//进行第二次拆分
		vvs := strings.Split(v, ":")
		if len(vvs) == 1 {
			err := m.checkSinglePar(rule, vvs[0])
			if err != nil {
				return nil, err
			}
		} else if len(vvs) == 2 {
			err := m.checkPairPar(rule, kind, vvs)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("parameter format error:%s", defined)
		}
	}
	return rule, nil
}

func (m *Matcher) checkSinglePar(
	rule *Rule, defined string) error {
	switch defined {
	case Required:
		rule.Required = true
	case Optional:
		rule.Optional = true
	default:
		return fmt.Errorf("unsupported parameter type %s", defined)
	}
	return nil
}

func (m *Matcher) checkPairPar(
	rule *Rule, kind string, vvs []string) error {
	var err error
	switch vvs[0] {
	case Eq:
		err = m.checkEq(rule, kind, vvs)
	case Neq:
		err = m.checkNeq(rule, kind, vvs)
	case Lt:
		err = m.checkLt(rule, kind, vvs)
	case Lte:
		err = m.checkLte(rule, kind, vvs)
	case Gt:
		err = m.checkGt(rule, kind, vvs)
	case Gte:
		err = m.checkGte(rule, kind, vvs)
	case Regexp:
		err = m.checkRegexp(rule, vvs)
	case Description:
		err = m.checkDc(rule, vvs)
	case Between:
		err = m.checkBetween(rule, kind, vvs)
	case Size:
		err = m.checkSize(rule, kind, vvs)
	case Len:
		err = m.checkLen(rule, kind, vvs)
	case In:
		err = m.checkIn(rule, kind, vvs)
	case NotIn:
		err = m.checkNotIn(rule, kind, vvs)
	default:
		return fmt.Errorf("unsupported parameter type %s", vvs[0])
	}
	return err
}

func (m *Matcher) checkIn(
	rule *Rule, kind string, vvs []string) error {
	bts := strings.Split(vvs[1], ",")
	for _, v := range bts {
		err := m.checkKindVal(kind, vvs[0], v)
		if err != nil {
			return err
		}
	}
	rule.In = bts
	return nil
}

func (m *Matcher) checkNotIn(
	rule *Rule, kind string, vvs []string) error {
	bts := strings.Split(vvs[1], ",")
	for _, v := range bts {
		err := m.checkKindVal(kind, vvs[0], v)
		if err != nil {
			return err
		}
	}
	rule.NotIn = bts
	return nil
}

func (m *Matcher) checkSize(
	rule *Rule, kind string, vvs []string) error {
	bts := strings.Split(vvs[1], ",")
	if len(bts) != 2 {
		return fmt.Errorf("the value defined by the %s parameter is in the wrong format", vvs[0])
	}
	for _, v := range bts {
		err := m.checkKindVal(kind, vvs[0], v)
		if err != nil {
			return err
		}
	}
	rule.Size = bts
	return nil
}

func (m *Matcher) checkLen(
	rule *Rule, kind string, vvs []string) error {
	err := m.checkKindVal(kind, vvs[0], vvs[1])
	if err != nil {
		return err
	}
	rule.Len = vvs[1]
	return nil
}

func (m *Matcher) checkBetween(
	rule *Rule, kind string, vvs []string) error {
	bts := strings.Split(vvs[1], ",")
	if len(bts) != 2 {
		return fmt.Errorf("the value defined by the %s parameter is in the wrong format", vvs[0])
	}
	for _, v := range bts {
		err := m.checkKindVal(kind, vvs[0], v)
		if err != nil {
			return err
		}
	}
	rule.Between = bts
	return nil
}

func (m *Matcher) checkEq(
	rule *Rule, kind string, vvs []string) error {
	err := m.checkKindVal(kind, vvs[0], vvs[1])
	if err != nil {
		return err
	}
	rule.Eq = vvs[1]
	return nil
}

func (m *Matcher) checkNeq(
	rule *Rule, kind string, vvs []string) error {
	err := m.checkKindVal(kind, vvs[0], vvs[1])
	if err != nil {
		return err
	}
	rule.Neq = vvs[1]
	return nil
}

func (m *Matcher) checkLt(
	rule *Rule, kind string, vvs []string) error {
	err := m.checkKindVal(kind, vvs[0], vvs[1])
	if err != nil {
		return err
	}
	rule.Lt = vvs[1]
	return nil
}

func (m *Matcher) checkLte(
	rule *Rule, kind string, vvs []string) error {
	err := m.checkKindVal(kind, vvs[0], vvs[1])
	if err != nil {
		return err
	}
	rule.Lte = vvs[1]
	return nil
}

func (m *Matcher) checkGt(
	rule *Rule, kind string, vvs []string) error {
	err := m.checkKindVal(kind, vvs[0], vvs[1])
	if err != nil {
		return err
	}
	rule.Gt = vvs[1]
	return nil
}

func (m *Matcher) checkGte(
	rule *Rule, kind string, vvs []string) error {
	err := m.checkKindVal(kind, vvs[0], vvs[1])
	if err != nil {
		return err
	}
	rule.Gte = vvs[1]
	return nil
}

func (m *Matcher) checkRegexp(
	rule *Rule, vvs []string) error {
	rule.Regexp = vvs[1]
	return nil
}

func (m *Matcher) checkDc(
	rule *Rule, vvs []string) error {
	rule.Dc = vvs[1]
	return nil
}

func (m *Matcher) checkKindVal(
	kind string, key, value string) error {
	switch kind {
	case reflect.Int8.String():
		d, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d < math.MinInt8 {
			return fmt.Errorf("the value defined by the %s parameter is less than the minimum value of the int8 type", key)
		}
		if d > math.MaxInt8 {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the int8 type", key)
		}
	case reflect.Uint8.String():
		d, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d > math.MaxUint8 {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the uint8 type", key)
		}

	case reflect.Int16.String():
		d, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d < math.MinInt16 {
			return fmt.Errorf("the value defined by the %s parameter is less than the minimum value of the int16 type", key)
		}
		if d > math.MaxInt16 {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the int16 type", key)
		}
	case reflect.Uint16.String():
		d, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d > math.MaxUint16 {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the uint16 type", key)
		}

	case reflect.Int32.String():
		d, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d < math.MinInt32 {
			return fmt.Errorf("the value defined by the %s parameter is less than the minimum value of the int32 type", key)
		}
		if d > math.MaxInt32 {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the int32 type", key)
		}
	case reflect.Uint32.String():
		d, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d > math.MaxUint32 {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the uint32 type", key)
		}

	case reflect.Int.String():
		d, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d < math.MinInt {
			return fmt.Errorf("the value defined by the %s parameter is less than the minimum value of the int type", key)
		}
		if d > math.MaxInt {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the int type", key)
		}
	case reflect.Uint.String():
		d, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d > math.MaxUint {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the uint type", key)
		}

	case reflect.Int64.String():
		d, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d < math.MinInt64 {
			return fmt.Errorf("the value defined by the %s parameter is less than the minimum value of the int64 type", key)
		}
		if d > math.MaxInt64 {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the int64 type", key)
		}
	case reflect.Uint64.String():
		d, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d > math.MaxUint64 {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the uint64 type", key)
		}

	case reflect.Float32.String():
		d, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d > math.MaxFloat32 {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the float32 type", key)
		}

	case reflect.Float64.String():
		d, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("the value defined by the %s parameter is in the wrong format", key)
		}
		if d > math.MaxFloat64 {
			return fmt.Errorf("the value defined by the %s parameter is greater than the maximum value of the float64 type", key)
		}
	}
	return nil
}

func (m *Matcher) check() error {
	vT := m.v.Type()
	if vT.Kind() == reflect.Ptr && vT.Elem().Kind() == reflect.Struct {
		return nil
	}
	return fmt.Errorf("please enter a pointer of type %s", m.t.Name())
}
