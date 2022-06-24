package adon

import (
	"path/filepath"
	goplugin "plugin"
	"reflect"
)

type Plugin interface {
	GetName() string
	GetVariableStorage() VariableStorage
	GetFunctionStorage() FunctionStorage
}

type PluginStorage = Storage[Plugin]

func NewPluginStorage() PluginStorage {
	return newStorage[Plugin]()
}

type plugin struct {
	name            string
	functionStorage FunctionStorage
	variableStorage VariableStorage
}

func (p plugin) GetName() string {
	return p.name
}

func (p plugin) GetVariableStorage() VariableStorage {
	return p.variableStorage
}

func (p plugin) GetFunctionStorage() FunctionStorage {
	return p.functionStorage
}

func NewPluginFromFile(path string) (Plugin, error) {
	goPlugin, err := goplugin.Open(path)
	if err != nil {
		return nil, err
	}

	iter := reflect.ValueOf(*goPlugin).FieldByName("syms").MapRange()
	valueRecords := []Record[reflect.Value]{}
	for iter.Next() {
		key := iter.Key().String()
		symbol, err := goPlugin.Lookup(key)
		if err != nil {
			return nil, err
		}
		valueRecords = append(valueRecords, Record[reflect.Value]{
			Name:  key,
			Value: reflect.ValueOf(symbol),
		})
	}

	return NewPlugin(filepath.Base(path), valueRecords), nil
}

func NewPlugin(name string, recordList []Record[reflect.Value]) Plugin {
	functionStorage := NewFunctionStorage()
	variableStorage := NewVariableStorage()
	for _, record := range recordList {
		switch {
		case IsFunctionKind(record.Value.Kind()):
			functionStorage.Set(Record[Function]{
				Name:  record.Name,
				Value: NewFunction(record.Value),
			})
		case IsVariableKind(record.Value.Kind()):
			variableStorage.Set(Record[Variable]{
				Name:  record.Name,
				Value: NewVariable(record.Value),
			})
		}
	}

	return plugin{
		functionStorage: functionStorage,
		variableStorage: variableStorage,
		name:            name,
	}
}
