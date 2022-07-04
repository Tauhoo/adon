package adon

import (
	"fmt"
	"reflect"
)

type Variable interface {
	GetValue() reflect.Value
}

type VariableStorage = Storage[Variable]

func NewVariableStorage() VariableStorage {
	return newStorage[Variable]()
}

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
	reflect.Float32: true,
	reflect.Float64: true,
	reflect.String:  true,
}

func (v variable) GetValue() reflect.Value {
	return reflect.Value(v)
}

func NewVariable(value reflect.Value) Variable {
	if !IsVariableKind(value.Kind()) {
		panic(fmt.Errorf("%w - got: %s", ErrInvalidValueKind, value.Kind().String()))
	}
	return variable(value)
}

func IsVariableKind(kind reflect.Kind) bool {
	_, isPrimitive := primitiveKind[kind]
	return isPrimitive
}

func ConvertVariableListToKindList(variableList []Variable) []reflect.Kind {
	kindList := []reflect.Kind{}
	for _, vavariable := range variableList {
		kindList = append(kindList, vavariable.GetValue().Kind())
	}
	return kindList
}
