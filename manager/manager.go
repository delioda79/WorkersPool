package manager

import (
	"sync"
)

type Task func()

func runTask(task Task, done chan bool) {
	task()
	done <- true
}

// Manager represents our worker pool manager
type Manager interface {
	Run()
	Send(tsk Task)
}

// baseManager represents our worker pool manager
type baseManager struct {
	pool        []Task
	inCh        chan Task
	runners     int
	totRunners  int
	mutex       *sync.Mutex
	runnersChan chan bool
}

// NewManager returns a new manager
func NewManager(workers int) Manager {

	return &baseManager{
		pool:        []Task{},
		runners:     0,
		totRunners:  workers,
		mutex:       &sync.Mutex{},
		inCh:        make(chan Task),
		runnersChan: make(chan bool),
	}
}

// Run starts the exporter, it spins up a new coroutine
func (scr *baseManager) Run() {
	go scr.runWriter()
	go scr.runReader()
}

// Send sends a message to the exporter's loop
func (scr baseManager) Send(tsk Task) {
	scr.inCh <- tsk
	return
}

func (scr baseManager) popTask() (*Task, []Task) {

	if len(scr.pool) == 0 {
		return nil, scr.pool
	}
	if scr.runners >= scr.totRunners {
		return nil, scr.pool
	}

	head := &scr.pool[0]
	tail := scr.pool[1:]

	return head, tail
}
func (scr *baseManager) runWriter() {
	head, tail := scr.popTask()
	if head != nil {
		scr.pool = tail
		scr.runners++
		go runTask(*head, scr.runnersChan)
	}
}

func (scr *baseManager) runReader() {
	for {
		scr.mutex.Lock()
		select {
		case tsk := <-scr.inCh:
			scr.pool = append(scr.pool, tsk)
		case <-scr.runnersChan:
			scr.runners--
		}
		if scr.runners < scr.totRunners {
			scr.runWriter()
		}
		scr.mutex.Unlock()
	}
}
