package adon

import (
	"fmt"
	"reflect"
)

type Function interface {
	Call(params ...Variable) []Variable
}

type function struct {
	parameterSignatureList []VariableSignature
	resultSignatureList    []VariableSignature
	value                  Value
}

func (f function) ValidateArguments(params ...Variable) error {
	if len(params) != len(f.parameterSignatureList) {
		return fmt.Errorf("%w - argument length invalid want: %d, got: %d",
			ErrInvalidFunctionArguments,
			len(f.parameterSignatureList),
			len(params),
		)
	}

	for index, signature := range f.parameterSignatureList {
		if signature.GetKind() != params[index].GetVariableSignature().GetKind() {
			return fmt.Errorf("%w - argument kind invalid at index: %d, want: %v, got: %v",
				ErrInvalidFunctionArguments,
				index,
				signature.GetKind(),
				params[index].GetVariableSignature().GetKind(),
			)
		}
	}

	return nil
}

func (f function) Call(params ...Variable) ([]Variable, error) {
	if err := f.ValidateArguments(params...); err != nil {
		return nil, err
	}

	arguments := []reflect.Value{}
	for _, param := range params {
		arguments = append(arguments, reflect.Value(param.GetValue()))
	}

	callResultList := reflect.Value(f.value).Call(arguments)
	resultList := []Variable{}

	for index, signature := range f.resultSignatureList {
		variable := Variable{
			signature: signature,
			value:     Value(callResultList[index]),
		}

		resultList = append(resultList, variable)
	}

	return resultList, nil
}
