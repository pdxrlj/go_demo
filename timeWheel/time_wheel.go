package main

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type TimeWheel struct {
	interval    time.Duration //时间的精度
	slotNums    int           //时间轮盘的齿轮数 走一圈的时间 interval*slotNums
	currentPos  int           //时间轮的当前的位置
	ticker      *time.Ticker
	slots       []*list.List //时间轮每一个位置对应的job
	taskRecords sync.Map
	isRunning   bool
}

type Job func(any) error

type Task struct {
	key       any //标识task对象，唯一
	interval  time.Duration
	createdAt time.Time
	pos       int //任务轮盘的位置
	circle    int //任务在轮盘走多少圈再执行
	job       Job //任务要执行的job
	times     int //任务要执行的次数
}

func (t *TimeWheel) Start() {
	t.ticker = time.NewTicker(t.interval)
	t.isRunning = true
	go t.run()
}

func (t *TimeWheel) run() {
	for {
		select {
		case <-t.ticker.C:
			t.checkAndRunTask()
		}
	}
}

func (t *TimeWheel) checkAndRunTask() {
	//获取当前轮盘位置的链表
	currentList := t.slots[t.currentPos]
	if currentList != nil {
		for item := currentList.Front(); item != nil; {
			task := item.Value.(*Task)
			if task.circle > 0 {
				//当前任务的循环次数还没有到执行的时候
				task.circle--
				item = item.Next()
				continue
			}

			//到时间要执行了
			if task.job != nil {
				go func() {
					if err := task.job(task.key); err != nil {
						log.Fatal(err)
					}
				}()
			}

			//任务执行完毕后，删除任务
			next := item.Next()
			t.taskRecords.Delete(task.key)
			currentList.Remove(item)
			item = next

			//判断是否是重复执行的任务
			//小于0表示一直执行
			if task.times < 0 {
				t.AddTask(task)
			}

			if task.times > 0 {
				task.times--
				t.AddTask(task)
			}
		}
	}

	//轮盘往前运行一步
	if t.currentPos == t.slotNums-1 {
		t.currentPos = 0
	} else {
		t.currentPos++
	}

}

func (t *TimeWheel) AddTask(task *Task) (int, int) {
	passedTime := time.Since(task.createdAt)
	delaySeconds := int(task.interval.Seconds())

	passedSeconds := int(passedTime.Seconds())
	intervalSeconds := int(t.interval.Seconds())
	circle := delaySeconds / intervalSeconds / t.slotNums
	pos := (t.currentPos + (delaySeconds-(passedSeconds%delaySeconds))/intervalSeconds) % t.slotNums
	if pos == t.currentPos && circle != 0 {
		circle--
	}

	return pos, circle
}

func (t *TimeWheel) Stop() {
	t.ticker.Stop()
	t.isRunning = false
}
