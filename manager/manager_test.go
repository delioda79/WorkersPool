package manager

import (
	"fmt"
	"testing"
	"time"
)

func aTask(n int, ch chan int) Task {
	return func() {
		fmt.Println("Running")
		time.Sleep(time.Millisecond * 3)
		ch <- n
	}
}

func TestManager(t *testing.T) {
	nOTasks := 10000
	manager := NewManager(100)

	manager.Run()
	ch := make(chan int, 1)
	for i := 0; i < nOTasks; i++ {
		manager.Send(aTask(1, ch))
	}

	time.Sleep(time.Second * 5)

	for i := 0; i < nOTasks; i++ {
		manager.Send(aTask(1, ch))
	}

	count := 0
	for {
		<-ch
		fmt.Println("Received")
		count++
		if count == nOTasks*2 {
			break
		}
	}
}
