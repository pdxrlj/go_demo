package main

import (
	"fmt"
	"runtime"
	"time"

	"woker_pool/worker"
)

type Task struct {
	Number int
}

func (t *Task) RunTask() {
	fmt.Println("This is task: ", t.Number)
	//设置个等待时间
	time.Sleep(1 * time.Second)
}

// 参考链接 http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/
func main() {
	poolNum := 10
	jobQueueNum := 100
	workerPool := worker.NewPool(poolNum, jobQueueNum)
	workerPool.Start()

	// 模拟百万请求
	dataNum := 100
	for i := 0; i < dataNum; i++ {
		t := Task{Number: i}
		workerPool.JobQueue <- &t
	}
	// 等待所有任务完成
	for {
		// 当前的goroutine数量因为poolNum+2
		fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
		time.Sleep(2 * time.Second)
	}
}
