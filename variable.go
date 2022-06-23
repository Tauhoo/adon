package adon

import (
	"reflect"
)

type variable reflect.Value

func (v variable) GetValue() reflect.Value {
	return reflect.Value(v)
}
