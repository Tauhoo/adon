package adon

import (
	"io/ioutil"
)

type PluginManager interface {
	GetPluginStorage() PluginStorage
	LoadPluginFromFile(path string) error
	LoadPluginFromFolder(path string) error
}

type pluginManager struct {
	pluginStorage PluginStorage
}

func (p pluginManager) GetPluginStorage() PluginStorage {
	return p.pluginStorage
}

func (p pluginManager) LoadPluginFromFile(path string) error {
	plugin, err := NewPluginFromFile(path)
	if err != nil {
		return err
	}

	p.pluginStorage.Set(Record[Plugin]{
		Name:  plugin.GetName(),
		Value: plugin,
	})

	return nil
}

func (p pluginManager) LoadPluginFromFolder(path string) error {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := p.LoadPluginFromFile(file.Name()); err != nil {
			p.GetPluginStorage().DeleteAll()
			return err
		}
	}

	return nil
}

func NewPluginManager() PluginManager {
	return pluginManager{
		pluginStorage: NewPluginStorage(),
	}
}