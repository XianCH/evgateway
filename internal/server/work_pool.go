package server

import (
	"fmt"
	"sync"
)

type Task func()

type WorkerPool struct {
	taskQueue chan Task
	wg        sync.WaitGroup
}

func NewWorkerPool(size int) *WorkerPool {
	return &WorkerPool{
		taskQueue: make(chan Task, 1000),
	}
}

func (wp *WorkerPool) Start(size int) {
	for i := 0; i < size; i++ {
		wp.wg.Add(1)
		go func(id int) {
			defer wp.wg.Done()
			for task := range wp.taskQueue {
				task()
			}
		}(i)
	}
	fmt.Println("[worker_pool] started with", size, "workers")
}

func (wp *WorkerPool) Submit(task Task) {
	wp.taskQueue <- task
}

func (wp *WorkerPool) Stop() {
	close(wp.taskQueue)
	wp.wg.Wait()
	fmt.Println("[worker_pool] stopped")
}
