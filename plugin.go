package adon

import (
	"path/filepath"
	goplugin "plugin"
	"reflect"
)

type Plugin interface {
	GetName() string
	GetVariableStorage() VariableStorage
	GetExecutorStorage() ExecutorStorage
}

type PluginStorage = Storage[Plugin]

func NewPluginStorage() PluginStorage {
	return newStorage[Plugin]()
}

type plugin struct {
	name            string
	executorStorage ExecutorStorage
	variableStorage VariableStorage
}

func (p plugin) GetName() string {
	return p.name
}

func (p plugin) GetVariableStorage() VariableStorage {
	return p.variableStorage
}

func (p plugin) GetExecutorStorage() ExecutorStorage {
	return p.executorStorage
}

func NewPluginFromFile(jobInstance Job, path string) (Plugin, error) {
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

	return NewPlugin(jobInstance, filepath.Base(path), valueRecords), nil
}

func NewPlugin(jobInstance Job, name string, recordList []Record[reflect.Value]) Plugin {
	executorStorage := NewExecutorStorage()
	variableStorage := NewVariableStorage()
	for _, record := range recordList {
		switch {
		case IsFunctionKind(record.Value.Kind()):
			executorStorage.Set(Record[Executor]{
				Name:  record.Name,
				Value: NewExecutor(NewFunction(record.Value), jobInstance),
			})
		case IsVariableKind(record.Value):
			variableStorage.Set(Record[Variable]{
				Name:  record.Name,
				Value: NewVariableFromPointer(record.Value),
			})
		}
	}

	return plugin{
		executorStorage: executorStorage,
		variableStorage: variableStorage,
		name:            name,
	}
}
