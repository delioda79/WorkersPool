package manager

import (
	"fmt"
)

type Task func()

type TaskRunner struct {
	ch chan Task
	n  int
}

func (tr TaskRunner) Run() {
	fmt.Printf("Runner %d up and running\n", tr.n)
	//tr.ch <- func() {}
	for {
		task := <-tr.ch
		//fmt.Printf("Runner %d running ", tr.n)
		task()
	}
}

type Manager struct {
	pool        []Task
	inCh, outCh chan Task
	runners     []TaskRunner
}

// Run starts the exporter, it spins up a new coroutine
func (scr *Manager) Run(nOfRunners int) {
	fmt.Println("Running the exporter")
	scr.inCh = make(chan Task, 1)
	fmt.Println("Channel created", scr.inCh)
	scr.outCh = make(chan Task, 1)
	fmt.Println("Channel created", scr.outCh)
	fmt.Printf("Create %d runners\n", nOfRunners)
	scr.runners = []TaskRunner{}
	for i := 0; i < nOfRunners; i++ {
		task := TaskRunner{
			ch: scr.outCh,
			n:  i,
		}
		scr.runners = append(scr.runners, task)
		go task.Run()
	}
	// i := 0
	// for {
	// 	<-scr.outCh
	// 	fmt.Println("Started")
	// 	i++
	// 	scr.outCh <- func() {}
	// 	if i == nOfRunners {
	// 		break
	// 	}
	// }
	go scr.runReader()
}

// Send sends a message to the exporter's loop
func (scr Manager) Send(tsk Task) {
	//if tsk == nil {
	//fmt.Println("Sent", tsk)
	//}
	scr.inCh <- tsk
	return
}

func (scr Manager) popTask() (*Task, []Task) {
	if len(scr.pool) == 0 {
		return nil, scr.pool
	}

	head := &scr.pool[0]
	tail := scr.pool[1:]

	return head, tail
}
func (scr *Manager) runWriter() {
	for {
		head, tail := scr.popTask()
		if head != nil {
			scr.pool = tail
			scr.outCh <- *head
		}
	}
}
func (scr *Manager) runReader() {
	go scr.runWriter()
	for {
		msg := <-scr.inCh
		scr.pool = append(scr.pool, msg)
	}
}
