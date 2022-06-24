package adon

type Record[T any] struct {
	value T
	name  string
}

type Storage[T any] interface {
	Set(record Record[T])
	Delete(name string)
	Find(name string) (Record[T], bool)
	GetList() []Record[T]
	GetByFilter(filter func(Record[T]) bool) []Record[T]
}

type storage[T any] struct {
	valueMap map[string]Record[T]
}

func (vm *storage[T]) Set(record Record[T]) {
	vm.valueMap[record.name] = record
}

func (vm *storage[T]) Delete(name string) {
	delete(vm.valueMap, name)
}

func (vm *storage[T]) Find(name string) (Record[T], bool) {
	value, ok := vm.valueMap[name]
	return value, ok
}

func (vm *storage[T]) GetList() []Record[T] {
	list := []Record[T]{}
	for _, v := range vm.valueMap {
		list = append(list, v)
	}
	return list
}

func (vm *storage[T]) GetByFilter(filter func(Record[T]) bool) []Record[T] {
	list := []Record[T]{}
	for _, v := range vm.valueMap {
		if filter(v) {
			list = append(list, v)
		}
	}
	return list
}

func newStorage[T any]() Storage[T] {
	return &storage[T]{
		valueMap: map[string]Record[T]{},
	}
}
