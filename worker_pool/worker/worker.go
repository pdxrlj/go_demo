package worker

import (
	"woker_pool/job"
)

type StopSignal struct{}

type Worker struct {
	Receiver job.Chan
	Quit     chan StopSignal
}

func NewWorker() *Worker {
	return &Worker{
		Receiver: make(job.Chan),
		Quit:     make(chan StopSignal),
	}
}

func (w *Worker) Start(workerPool *Pool) {
	go func() {
		for {
			// 将当前的worker
			workerPool.WorkerQueue <- w
			select {
			// 从Receiver channel中取出任务
			case task := <-w.Receiver:
				// 运行任务
				task.RunTask()
			case <-w.Quit:
				return
			}
		}
	}()
}
