package manager

import (
	"fmt"
	"testing"
	"time"
)

func aTask(n int, ch chan int) Task {
	return func() {
		fmt.Println("Running")
		time.Sleep(3000000)
		ch <- n
	}
}

func TestManager(t *testing.T) {
	manager := Manager{}

	manager.Run(10)
	ch := make(chan int, 1)
	for i := 0; i < 10; i++ {
		manager.Send(aTask(1, ch))
	}

	count := 0
	for {
		<-ch
		fmt.Println("Received")
		count++
		if count == 6 {
			break
		}
	}
}
