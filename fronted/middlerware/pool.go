package middlerware

import (
	"fmt"
	"sync"
	"time"
)

// 任务结构体
type Task struct {
	ID  int
	Job func()
}

// 协程池结构体
type Pool struct {
	taskQueue chan Task      //一个任务队列 taskQueue
	wg        sync.WaitGroup //一个 WaitGroup wg
}

// 创建协程池并运行
// NewPool 函数用于创建一个协程池，参数 numWorkers 指定了协程池中的工作协程数量。
// 在 NewPool 函数中，会初始化 taskQueue 通道，并启动指定数量的工作协程。
func NewPool(numWorkers int) *Pool {
	p := &Pool{
		taskQueue: make(chan Task),
	}

	p.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		//通过 go p.worker() 语句启动了指定数量的工作协程，这些工作协程会立即执行 worker 函数。
		//在 worker 函数中，通过 for task := range p.taskQueue 循环
		//工作协程会不断地从任务队列 taskQueue 中取出任务并执行。
		//如果任务队列为空，则工作协程会阻塞在取任务的操作上，直到有新的任务到来或者任务队列被关闭
		//因此初始时，这些协程会被阻塞
		go p.worker(i)
	}

	return p
}

// 添加任务到协程池
func (p *Pool) AddTask(task Task) {
	p.taskQueue <- task
}

// 工作协程
func (p *Pool) worker(workerID int) {
	//对于每个协程，不断从taskQueue取得任务
	for task := range p.taskQueue {
		fmt.Printf("%d new start task %d\n", workerID, task.ID)
		task.Job()
		time.Sleep(4 * time.Second)
		fmt.Printf("finished task %d\n", task.ID)
	}
	p.wg.Done()
}

// 等待所有任务完成
func (p *Pool) Wait() {
	close(p.taskQueue)
	p.wg.Wait()
}

//func main() {
//	// 创建一个协程池，设置工作协程数为3
//	pool := NewPool(10)
//
//	startTime := time.Now()
//	// 添加任务到协程池
//	for i := 0; i < 100; i++ {
//		taskID := i
//		task := Task{
//			ID: taskID,
//			Job: func() {
//				time.Sleep(time.Second)
//				fmt.Printf("Task %d is running\n", taskID)
//			},
//		}
//		pool.AddTask(task)
//	}
//
//	// 等待所有任务完成
//	pool.Wait()
//
//	endTime := time.Now()
//
//	// 计算处理时间并打印
//	processingTime := endTime.Sub(startTime)
//	fmt.Printf("Processing time: %v\n", processingTime)
//}
