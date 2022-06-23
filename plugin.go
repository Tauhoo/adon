package adon

import (
	"fmt"
	goplugin "plugin"
	"reflect"
)

type plugin struct {
	name            string
	functionStorage Storage[function]
	variableStorage Storage[variable]
}

func (p plugin) GetName() string {
	return p.name
}

func (p plugin) GetVariableStorage() Storage[variable] {
	return p.variableStorage
}

func (p plugin) GetFunctionStorage() Storage[function] {
	return p.functionStorage
}

func GetValueMapFromGoPlugin(goPlugin goplugin.Plugin) map[string]reflect.Value {
	iter := reflect.ValueOf(goPlugin).FieldByName("syms").MapRange()
	valueMap := map[string]reflect.Value{}
	for iter.Next() {
		key := iter.Key().String()
		s, err := goPlugin.Lookup(key)
		if err != nil {
			panic(fmt.Errorf("%w - look up for symbol fail symbol: %s", err, key))
		}
		value := reflect.ValueOf(s)
		valueMap[key] = value
	}
	return valueMap
}

func NewPlugin(name string, valueMap map[string]reflect.Value) plugin {
	functionStorage := NewStorage[function]()
	variableStorage := NewStorage[variable]()
	for k, v := range valueMap {
		switch {
		case IsFunctionKind(v.Kind()):
			functionStorage.Set(Record[function]{
				name:  k,
				value: NewFunction(v),
			})
		case IsVariableKind(v.Kind()):
			variableStorage.Set(Record[variable]{
				name:  k,
				value: NewVariable(v),
			})
		}
	}

	return plugin{
		functionStorage: functionStorage,
		variableStorage: variableStorage,
		name:            name,
	}
}
