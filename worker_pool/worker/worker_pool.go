package worker

import (
	"woker_pool/job"
)

type Pool struct {
	Size        int          // number of workers
	JobQueue    job.Chan     // job queue
	WorkerQueue chan *Worker // worker queue
}

func NewPool(poolSize, jobQueueLen int) *Pool {
	return &Pool{
		Size:        poolSize,                     // 协程池大小
		JobQueue:    make(job.Chan, jobQueueLen),  // 任务队列的大小
		WorkerQueue: make(chan *Worker, poolSize), // poolSize个协程
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.Size; i++ {
		worker := NewWorker()
		worker.Start(p)
	}

	go func() {
		for {
			select {
			// 从任务队列中取出任务
			case task := <-p.JobQueue:
				// 从协程池中取出一个协程,如果没有会阻塞等待
				worker := <-p.WorkerQueue
				// 将任务交给worker协程
				worker.Receiver <- task
			}
		}
	}()
}
