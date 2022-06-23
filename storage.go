package adon

type Storage[T any] struct {
	valueMap map[string]T
}

func (vm *Storage[T]) Set(name string, value T) {
	vm.valueMap[name] = value
}

func (vm *Storage[T]) Delete(name string) {
	delete(vm.valueMap, name)
}

func (vm *Storage[T]) Find(name string) (T, bool) {
	value, ok := vm.valueMap[name]
	return value, ok
}

func (vm *Storage[T]) GetList() []T {
	list := []T{}
	for _, v := range vm.valueMap {
		list = append(list, v)
	}
	return list
}
