package lang

import (
	"fmt"
	"gojs/ast"
	"strconv"
)

type ValueType int

const (
	ValueTypeUndefined ValueType = iota
	ValueTypeNull
	ValueTypeStr
	ValueTypeInt
	ValueTypeBool
	ValueTypeObj
)

type Value struct {
	Type ValueType
	Str  string
	Int  int
	Bool bool
	Obj  Object
}

func (v Value) String() string {
	if v.Type == ValueTypeStr {
		return v.Str
	}

	if v.Type == ValueTypeInt {
		return fmt.Sprintf("%d", v.Int)
	}

	if v.Type == ValueTypeBool {
		return fmt.Sprintf("%t", v.Bool)
	}

	if v.Type == ValueTypeObj {
		return "[object]"
	}

	if v.Type == ValueTypeUndefined {
		return "undefined"
	}

	if v.Type == ValueTypeNull {
		return "null"
	}

	return "[ILLEGAL]"
}

func NewUndefined() Value {
	return Value{Type: ValueTypeUndefined}
}

func NewNull() Value {
	return Value{Type: ValueTypeNull}
}

func NewStr(str string) Value {
	return Value{Type: ValueTypeStr, Str: str}
}

func NewInt(val int) Value {
	return Value{Type: ValueTypeInt, Int: val}
}

func NewBool(val bool) Value {
	return Value{Type: ValueTypeBool, Bool: val}
}

func NewObj(obj Object) Value {
	return Value{Type: ValueTypeObj, Obj: obj}
}

type Object interface {
	_Object()
	GetProperty(name string) Value
	SetProperty(name string, value Value)
}
type JsObject struct {
	Storage map[string]Value
}

func (j *JsObject) GetProperty(name string) Value {
	return j.Storage[name]
}

func (j *JsObject) SetProperty(name string, value Value) {
	j.Storage[name] = value
}

func (j *JsObject) _Object() {}

type Function struct {
	Name       string
	Body       ast.Statement
	Parameters []ast.Identifier
}

func (f *Function) GetProperty(name string) Value {
	panic("nope")
}

func (f *Function) SetProperty(name string, value Value) {
	panic("nope")
}

func (f *Function) _Object() {}

type NativeFunction struct {
	Function func(values ...Value)
}

func (f *NativeFunction) GetProperty(name string) Value {
	panic("nope")
}

func (f *NativeFunction) SetProperty(name string, value Value) {
	panic("nope")
}

func (f *NativeFunction) _Object() {}

type Array struct {
	Store []Value
}

func (a *Array) GetProperty(name string) Value {
	idx, err := strconv.Atoi(name)
	if err != nil {
		panic(err)
	}

	return a.Store[idx]
}

func (a *Array) SetProperty(name string, value Value) {
	idx, err := strconv.Atoi(name)
	if err != nil {
		panic(err)
	}
	a.Store[idx] = value
}

func (a *Array) _Object() {}
