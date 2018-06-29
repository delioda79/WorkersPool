# WorkersPool

A simple library to handle asynchronous tasks running with a limit of concurrent processes

In order to use the library first thing is to initialise the manager:

``` go

manager := NewManager(10)

```

Which creates a manager with maximum 10 concurrent processes. Each process is essentially a go routine which runs a task popped
from the task pool ( a queue).

The next step is now to run the manager:

``` go

manager.Run()

```

Which lets the manager wait for tasks to be sent. The run function starts two goroutines, so is non-blcoking and we can write any code
after that.

Let's now add some tasks to the pool. In order to do that we use the method ```Send``` on the manager and we send a task.
A task is a defined type which goes back to a simple:

``` go

type Task func()

```

Any function without arguments and returning nothing is a task. Let's see how this can be used to send any task.

``` go

func aTask(name string, age int) Task {
    return func() {
        fmt.Printf("Hi %s you are %d years old", name, age)
        time.Sleep(3000000)
    }
}

```

This function accepts some parameters and returns a Task (function with no params and returning nothing).
Let's now add tasks to the pool:

``` go
manager.Send(aTask("John", 32))
manager.Send(atask("Tom", 22))
manager.Send(atask("Mark", 48))

```

As soon as we send a Task, the task gets stored in the internal pool.
The firts task, when called, will print out:

```  bash
Hi John you are 32 years old

```

And so on. As soon as a task gets added to the pool it is then popped and a go routine gets created running that specific task.
This happens only if the number of tasks running is less or equal to the maximum number of concurrent processes set when initialising the manager.

If the number of task running is already equal to the maximum allowed the manager waits for a running process to end and then pops a new task from
the pool. 