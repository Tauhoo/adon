package adon

type ExecuteState uint16

const (
	ExecuteIdle ExecuteState = iota
	ExecuteError
	ExecuteRunning
	ExecuteDone
)

type StateEventListener = func(state ExecuteState, info any)

type StateEventPubliser interface {
	GetState() ExecuteState
	publish(state ExecuteState, info any)
	listen(state ExecuteState, fn StateEventListener)
}

type stateEventPubliser struct {
	state  ExecuteState
	mapper map[ExecuteState]([]StateEventListener)
}

func (s stateEventPubliser) GetState() ExecuteState {
	return s.state
}

func (s stateEventPubliser) publish(state ExecuteState, info any) {
	s.state = state
	listenerList, ok := s.mapper[state]
	if !ok {
		s.mapper[state] = []StateEventListener{}
		return
	}

	for _, listener := range listenerList {
		listener(state, info)
	}
}
func (s stateEventPubliser) listen(state ExecuteState, fn StateEventListener) {
	listenerList, ok := s.mapper[state]
	if ok {
		listenerList = append(listenerList, fn)
	} else {
		s.mapper[state] = []StateEventListener{fn}
	}
}

func NewStateEventPubliser() StateEventPubliser {
	return stateEventPubliser{
		state:  ExecuteIdle,
		mapper: map[ExecuteState][]StateEventListener{},
	}
}

type Executor interface {
	Execute(params ...Variable)
	Stop()
	GetFunction() Function
	GetState() ExecuteState
}

type executor struct {
	function           Function
	stateEventPubliser StateEventPubliser
	jobInstance        Job
	jobChannel         chan<- JobAction
}

func (e *executor) Execute(params ...Variable) {
	e.Stop()

	kindList := ConvertVariableListToKindList(params)
	if err := e.function.ValidateParams(kindList...); err != nil {
		e.stateEventPubliser.publish(ExecuteError, err)
		return
	}

	e.stateEventPubliser.publish(ExecuteRunning, e)
	fn := func() {
		result, err := e.function.Call(params...)
		if err != nil {
			e.stateEventPubliser.publish(ExecuteError, err)
		} else {
			e.stateEventPubliser.publish(ExecuteDone, result)
		}
	}

	e.jobChannel = e.jobInstance.Exec(fn)
}

func (e executor) Stop() {
	if e.stateEventPubliser.GetState() == ExecuteRunning {
		e.jobChannel <- JobStop
	}
}

func (e executor) GetFunction() Function {
	return e.function
}

func (e executor) GetState() ExecuteState {
	return e.stateEventPubliser.GetState()
}

func NewExecutor(fn Function, jobInstance Job) Executor {
	return &executor{
		function:           fn,
		stateEventPubliser: NewStateEventPubliser(),
		jobInstance:        jobInstance,
		jobChannel:         nil,
	}
}

type ExecutorStorage = Storage[Executor]

type executorStorage struct {
	storage Storage[Executor]
}

func (e executorStorage) stop(name string) {
	if result, ok := e.Find(name); ok {
		result.Value.Stop()
	}
}

func (e executorStorage) Set(record Record[Executor]) {
	e.stop(record.Name)
	e.storage.Set(record)
}
func (e executorStorage) Delete(name string) {
	e.stop(name)
	e.storage.Delete(name)
}
func (e executorStorage) Find(name string) (Record[Executor], bool) {
	return e.Find(name)
}
func (e executorStorage) GetList() []Record[Executor] {
	return e.GetList()
}
func (e executorStorage) GetByFilter(filter func(Record[Executor]) bool) []Record[Executor] {
	return e.storage.GetByFilter(filter)
}
func (e executorStorage) DeleteAll() {
	for _, record := range e.GetList() {
		record.Value.Stop()
	}
	e.storage.DeleteAll()
}

func NewExecutorStorage() ExecutorStorage {
	return executorStorage{
		storage: newStorage[Executor](),
	}
}
