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
type Manager struct {
	pool        []Task
	inCh        chan Task
	runners     int
	totRunners  int
	mutex       *sync.Mutex
	runnersChan chan bool
}

// NewManager returns a new manager
func NewManager(workers int) *Manager {

	return &Manager{
		pool:        []Task{},
		runners:     0,
		totRunners:  workers,
		mutex:       &sync.Mutex{},
		inCh:        make(chan Task),
		runnersChan: make(chan bool),
	}
}

// Run starts the exporter, it spins up a new coroutine
func (scr *Manager) Run() {
	go scr.runWriter()
	go scr.runReader()
}

// Send sends a message to the exporter's loop
func (scr Manager) Send(tsk Task) {
	scr.inCh <- tsk
	return
}

func (scr Manager) popTask() (*Task, []Task) {

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
func (scr *Manager) runWriter() {
	for {
		scr.mutex.Lock()
		head, tail := scr.popTask()
		if head != nil {
			scr.pool = tail
			scr.runners++
			go runTask(*head, scr.runnersChan)
		}
		scr.mutex.Unlock()
	}
}

func (scr *Manager) runReader() {
	for {
		select {
		case tsk := <-scr.inCh:
			scr.mutex.Lock()
			scr.pool = append(scr.pool, tsk)
			scr.mutex.Unlock()
		case <-scr.runnersChan:
			scr.mutex.Lock()
			scr.runners--
			scr.mutex.Unlock()
		}
	}
}
