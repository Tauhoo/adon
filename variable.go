package adon

import "reflect"

type Kind reflect.Kind
type Value reflect.Value

type VariableSignature struct {
	kind Kind
	name string
}

func (vs VariableSignature) GetName() string {
	return vs.name
}

func (vs VariableSignature) GetKind() Kind {
	return vs.kind
}

type Variable struct {
	signature VariableSignature
	value     Value
}

func (v Variable) GetValue() Value {
	return v.value
}

func (v Variable) GetVariableSignature() VariableSignature {
	return v.signature
}

type VariableManager interface {
	Set(value Variable)
	Delete(name string)
	Find(name string) (Variable, bool)
}

type variableManager struct {
	variableMap map[string]Variable
}

func (vm *variableManager) Set(value Variable) {
	vm.variableMap[value.GetVariableSignature().GetName()] = value
}

func (vm *variableManager) Delete(name string) {
	delete(vm.variableMap, name)
}

func (vm *variableManager) Find(name string) (Variable, bool) {
	value, ok := vm.variableMap[name]
	return value, ok
}

func (vm *variableManager) GetVariableList() []Variable {
	variableList := []Variable{}
	for _, v := range vm.variableMap {
		variableList = append(variableList, v)
	}
	return variableList
}

func NewVariableManager() VariableManager {
	return &variableManager{
		variableMap: map[string]Variable{},
	}
}
