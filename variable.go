package adon

import (
	"fmt"
	"reflect"
)

type variable reflect.Value

var primitiveKind = map[reflect.Kind]bool{
	reflect.Bool:    true,
	reflect.Int:     true,
	reflect.Int8:    true,
	reflect.Int16:   true,
	reflect.Int32:   true,
	reflect.Int64:   true,
	reflect.Uint:    true,
	reflect.Uint8:   true,
	reflect.Uint16:  true,
	reflect.Uint32:  true,
	reflect.Uint64:  true,
	reflect.Uintptr: true,
	reflect.Float32: true,
	reflect.Float64: true,
	reflect.String:  true,
}

func (v variable) GetValue() reflect.Value {
	return reflect.Value(v)
}

func NewVariable(value reflect.Value) variable {
	if IsVariableKind(value.Kind()) {
		panic(fmt.Errorf("%w - want: %s, got: %s", ErrInvalidValueKind, reflect.Func.String(), value.Kind().String()))
	}
	return variable(value)
}

func IsVariableKind(kind reflect.Kind) bool {
	_, isPrimitive := primitiveKind[kind]
	return isPrimitive
}
