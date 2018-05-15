package main

import (
	"fmt"
	"time"

	"github.com/delioda79/WorkersPool/manager"
)

func aTask(n int, ch chan int) manager.Task {
	return func() {
		fmt.Println("Running", n)
		time.Sleep(300000000)
		ch <- n
	}
}

func main() {
	manager := manager.Manager{}

	manager.Run(4)
	ch := make(chan int, 1)
	time.Sleep(3000000000)
	start := time.Now()
	nOfTasks := 20
	for i := 0; i < nOfTasks; i++ {
		manager.Send(aTask(i, ch))
	}

	count := 0
	for {
		<-ch
		fmt.Println("Received")
		count++
		if count == nOfTasks {
			break
		}
	}
	end := time.Now()

	fmt.Println("It took ", end.Sub(start).Seconds(), " seconds")
}
