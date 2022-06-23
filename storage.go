package adon

type Record[T any] struct {
	value T
	name  string
}

type Storage[T any] struct {
	valueMap map[string]Record[T]
}

func (vm *Storage[T]) Set(record Record[T]) {
	vm.valueMap[record.name] = record
}

func (vm *Storage[T]) Delete(name string) {
	delete(vm.valueMap, name)
}

func (vm *Storage[T]) Find(name string) (Record[T], bool) {
	value, ok := vm.valueMap[name]
	return value, ok
}

func (vm *Storage[T]) GetList() []Record[T] {
	list := []Record[T]{}
	for _, v := range vm.valueMap {
		list = append(list, v)
	}
	return list
}

func (vm *Storage[T]) GetByFilter(filter func(Record[T]) bool) []Record[T] {
	list := []Record[T]{}
	for _, v := range vm.valueMap {
		if filter(v) {
			list = append(list, v)
		}
	}
	return list
}

func NewStorage[T any]() Storage[T] {
	return Storage[T]{
		valueMap: map[string]Record[T]{},
	}
}
