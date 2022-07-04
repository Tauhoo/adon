package adon

type JobAction uint16

const (
	JobStop JobAction = iota
)

type Job interface {
	Start()
	Stop()
	Exec(fn func()) chan<- JobAction
}

type job struct {
	functionChanel  chan func()
	terminateChanel chan bool
}

func (j job) Stop() {
	j.terminateChanel <- false
}

func (j job) Exec(fn func()) chan<- JobAction {
	jobActionChanel := make(chan JobAction)
	wrapFunction := func() {
		doneChannel := make(chan bool)
		go func() {
			fn()
			doneChannel <- true
		}()

		for {
			select {
			case <-doneChannel:
				return
			case <-jobActionChanel:
				return
			}
		}
	}
	j.functionChanel <- wrapFunction
	return jobActionChanel
}

func (j job) Start() {
	go func() {
		for {
			select {
			case <-j.terminateChanel:
				return
			case fn := <-j.functionChanel:
				go fn()
			}
		}
	}()
}

func NewJob() Job {
	return job{
		functionChanel:  make(chan func()),
		terminateChanel: make(chan bool),
	}
}
