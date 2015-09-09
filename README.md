# GOTM
Golang Task Manager

# Example

```go
package main

import (
	"fmt"
	"github.com/kimiby/gotm"
	"time"
)

type CustomWorker struct {
	gotm.Worker
}

type Node struct {
	id int
}

func (w CustomWorker) Do(job *gotm.Job) { fmt.Println("Do custom work") }

func main() {
	a := CustomWorker{gotm.Worker{Type: "default"}}
	g, _ := gotm.Create(10)

	g.Register(&a)
	g.Dispatch()

	for {
		j := gotm.Job{"default", args}
		g.JobQueue <- j
		g.JobQueue <- j
		g.JobQueue <- j
		time.Sleep(time.Second * 10)
	}
}
```

# License
The MIT License (MIT)
