package adon

import (
	"io/ioutil"
	"path"
)

type PluginManager interface {
	GetPluginStorage() PluginStorage
	LoadPluginFromFile(path string) error
	LoadPluginFromFolder(path string) error
}

type pluginManager struct {
	jobInstance   Job
	pluginStorage PluginStorage
}

func (p pluginManager) GetPluginStorage() PluginStorage {
	return p.pluginStorage
}

func (p pluginManager) LoadPluginFromFile(path string) error {
	plugin, err := NewPluginFromFile(p.jobInstance, path)
	if err != nil {
		return err
	}

	p.pluginStorage.Set(Record[Plugin]{
		Name:  plugin.GetName(),
		Value: plugin,
	})

	return nil
}

func (p pluginManager) LoadPluginFromFolder(folderPath string) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := p.LoadPluginFromFile(path.Join(folderPath, file.Name())); err != nil {
			p.GetPluginStorage().DeleteAll()
			return err
		}
	}

	return nil
}

func NewPluginManager(jobInstance Job) PluginManager {
	return pluginManager{
		jobInstance:   jobInstance,
		pluginStorage: NewPluginStorage(),
	}
}
