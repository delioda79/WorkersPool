package manager

import (
	"fmt"
	"testing"
)

func aTask(n int, ch chan int) Task {
	return func() {
		fmt.Println(n)
		ch <- n
	}
}

func TestManager(t *testing.T) {
	manager := Manager{}

	manager.Run()
	ch := make(chan int, 1)
	manager.Send(aTask(1, ch))
	manager.Send(aTask(1, ch))
	manager.Send(aTask(1, ch))
	manager.Send(aTask(1, ch))
	manager.Send(aTask(1, ch))
	manager.Send(aTask(1, ch))

	count := 0
	for {
		<-ch
		count++
		if count == 6 {
			break
		}
	}
}
