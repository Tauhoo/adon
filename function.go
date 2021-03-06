package adon

import (
	"fmt"
	"reflect"
)

type Function interface {
	GetValue() reflect.Value
	GetParamList() []reflect.Kind
	GetReturnList() []reflect.Kind
	ValidateParams(params ...reflect.Kind) error
	Call(params ...Variable) ([]Variable, error)
}

type FunctionStorage = Storage[Function]

func NewFunctionStorage() FunctionStorage {
	return newStorage[Function]()
}

type function reflect.Value

func (f function) GetValue() reflect.Value {
	return reflect.Value(f)
}

func (f function) GetParamList() []reflect.Kind {
	result := []reflect.Kind{}
	reflectType := f.GetValue().Type()
	for i := 0; i < reflectType.NumIn(); i++ {
		result = append(result, reflectType.In(i).Kind())
	}
	return result
}

func (f function) GetReturnList() []reflect.Kind {
	result := []reflect.Kind{}
	reflectType := f.GetValue().Type()
	for i := 0; i < reflectType.NumOut(); i++ {
		result = append(result, reflectType.Out(i).Kind())
	}
	return result
}

func (f function) ValidateParams(params ...reflect.Kind) error {
	expectedParams := f.GetParamList()
	if len(params) != len(expectedParams) {
		return fmt.Errorf("%w - argument length invalid want: %d, got: %d",
			ErrInvalidFunctionArguments,
			len(expectedParams),
			len(params),
		)
	}

	for index, kind := range expectedParams {
		if kind != params[index] {
			return fmt.Errorf("%w - argument kind invalid at index: %d, want: %v, got: %v",
				ErrInvalidFunctionArguments,
				index,
				kind,
				params[index],
			)
		}
	}

	return nil
}

func (f function) Call(params ...Variable) ([]Variable, error) {
	values := []reflect.Value{}
	for _, param := range params {
		values = append(values, param.GetValue())
	}
	kinds := []reflect.Kind{}
	for _, param := range values {
		kinds = append(kinds, param.Kind())
	}

	if err := f.ValidateParams(kinds...); err != nil {
		return nil, err
	}

	callResults := f.GetValue().Call(values)

	results := []Variable{}
	for _, result := range callResults {
		pointerReflectValue := reflect.New(result.Type())
		pointerReflectValue.Elem().Set(result)
		results = append(results, NewVariableFromPointer(pointerReflectValue))
	}

	return results, nil
}

func NewFunction(value reflect.Value) Function {
	if value.Kind() != reflect.Func {
		panic(fmt.Errorf("%w - want: %s, got: %s", ErrInvalidValueKind, reflect.Func.String(), reflect.Func.String()))
	}
	return function(value)
}

func IsFunctionKind(kind reflect.Kind) bool {
	return kind == reflect.Func
}
