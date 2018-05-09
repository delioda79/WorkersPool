package manager

import (
	"fmt"
)

type Task func()

type Manager struct {
	pool []Task
	ch   chan Task
}

// Run starts the exporter, it spins up a new coroutine
func (scr *Manager) Run() {
	fmt.Println("Running the exporter")
	scr.ch = make(chan Task, 1)
	fmt.Println("Channel created", scr.ch)
	go scr.run()
}

// Send sends a message to the exporter's loop
func (scr Manager) Send(tsk Task) {
	scr.ch <- tsk
	return
}
func (scr Manager) run() {

	for {
		fmt.Println("Waiting for a message", scr.ch)
		msg := <-scr.ch
		fmt.Println("Message received, calling goroutine", msg)
		go msg()
	}
}
