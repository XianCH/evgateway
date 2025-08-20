package server

import (
	"sync"
	"testing"
	"time"
)

func TestWorkerPool_StartAndStop(t *testing.T) {
	wp := NewWorkerPool(5) // 修复：传递 size 参数
	wp.Start(5)

	if len(wp.taskQueue) != 0 {
		t.Errorf("expected task queue to be empty, got %d tasks", len(wp.taskQueue))
	}
	wp.Stop()
}

func TestWorkerPool_Submit(t *testing.T) {
	wp := NewWorkerPool(5) // 修复：传递 size 参数
	wp.Start(5)

	var wg sync.WaitGroup
	var executedTasks []int
	var mu sync.Mutex

	// Submit tasks
	for i := 0; i < 10; i++ {
		wg.Add(1)
		taskID := i
		wp.Submit(func() {
			defer wg.Done()
			mu.Lock()
			executedTasks = append(executedTasks, taskID)
			mu.Unlock()
		})
	}

	wg.Wait() // 等待所有任务完成
	wp.Stop()

	if len(executedTasks) != 10 {
		t.Errorf("expected 10 tasks to be executed, got %d", len(executedTasks))
	}
}

func TestWorkerPool_ConcurrentExecution(t *testing.T) {
	wp := NewWorkerPool(5) // 修复：传递 size 参数
	wp.Start(5)

	var mu sync.Mutex
	var executedTasks []int

	// Submit tasks
	for i := 0; i < 20; i++ {
		taskID := i
		wp.Submit(func() {
			time.Sleep(100 * time.Millisecond) // Simulate work
			mu.Lock()
			executedTasks = append(executedTasks, taskID)
			mu.Unlock()
		})
	}

	// Allow some time for tasks to complete
	time.Sleep(3 * time.Second)
	wp.Stop()

	// Verify all tasks were executed
	mu.Lock()
	if len(executedTasks) != 20 {
		t.Errorf("Expected 20 tasks to be executed, got %d", len(executedTasks))
	}
	mu.Unlock()
}
