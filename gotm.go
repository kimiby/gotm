package gotm

import (
	"errors"
	"fmt"
)

type Gotm struct {
	JobQueue      chan Job
	Workers       map[string][]WorkerIteraface
	WorkersQueues map[string]chan chan Job
}

func Create(q int) (Gotm, error) {
	var gotm Gotm
	var err error

	if q == 0 {
		err = errors.New("Invalid job queue length")
	} else {
		gotm = Gotm{
			JobQueue:      make(chan Job, q),
			Workers:       make(map[string][]WorkerIteraface),
			WorkersQueues: make(map[string]chan chan Job),
		}
	}

	return gotm, err
}

func (gotm *Gotm) Register(w WorkerIteraface) {
	w.Init()
	class := w.GetType()
	if _, ok := gotm.WorkersQueues[class]; !ok {
		gotm.WorkersQueues[class] = make(chan chan Job)
	}

	go func() {
		for {
			gotm.WorkersQueues[class] <- w.GetJobs()

			select {
			case job := <-w.GetJobs():
				w.Do(&job)
			case <-w.GetQuit():
				return
			}
		}
	}()

	gotm.Workers[class] = append(gotm.Workers[class], w)
}

func (gotm *Gotm) Dispatch() {
	go func() {
		for {
			select {
			case job := <-gotm.JobQueue:
				go func() {
					worker := <-gotm.WorkersQueues[job.Type]
					worker <- job
				}()
			}
		}
	}()
}

type WorkerIteraface interface {
	Init()
	GetType() string
	GetJobs() chan Job
	GetQuit() chan bool
	Do(job *Job)
}

type Worker struct {
	Type string
	Jobs chan Job
	Quit chan bool
}

func (w *Worker) GetJobs() chan Job {
	return w.Jobs
}

func (w *Worker) GetQuit() chan bool {
	return w.Quit
}

func (w *Worker) Do(job *Job) {
	fmt.Println("Custom worker 'Do' method not implemented")
}

func (w *Worker) Init() {
	w.Jobs = make(chan Job)
	w.Quit = make(chan bool)
}

func (w *Worker) GetType() string {
	return w.Type
}

type Job struct {
	Type string
	Args []interface{}
}
