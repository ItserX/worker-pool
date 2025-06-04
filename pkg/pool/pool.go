package pool

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const jobDuration = 50

type WorkerPool struct {
	workerCount int
	jobChan     chan string
	addChan     chan struct{}
	deleteChan  chan struct{}
	wg          *sync.WaitGroup
	mu          *sync.Mutex
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	wp := &WorkerPool{
		jobChan:    make(chan string),
		addChan:    make(chan struct{}),
		deleteChan: make(chan struct{}),
		wg:         &sync.WaitGroup{},
		mu:         &sync.Mutex{},
		ctx:        ctx,
		cancel:     cancel,
	}

	go wp.observeWorkers()

	for i := 0; i < numWorkers; i++ {
		wp.AddWorker()
	}

	return wp
}

func (wp *WorkerPool) observeWorkers() {
	for {
		select {
		case <-wp.ctx.Done():
			return

		case <-wp.addChan:
			wp.mu.Lock()
			wp.workerCount++
			id := wp.workerCount
			wp.wg.Add(1)
			wp.mu.Unlock()

			go wp.startWorker(id)
			fmt.Printf("worker %d added\n", id)
		}
	}
}

func (wp *WorkerPool) imitateJob(id int, data string) {
	fmt.Printf("worker %d processing job %s\n", id, data)
	time.Sleep(time.Millisecond * jobDuration)
}

func (wp *WorkerPool) startWorker(id int) {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.ctx.Done():
			fmt.Printf("worker %d stopped\n", id)
			return

		case data := <-wp.jobChan:
			wp.imitateJob(id, data)

		case <-wp.deleteChan:
			wp.mu.Lock()
			wp.workerCount--
			wp.mu.Unlock()
			fmt.Printf("worker %d deleted\n", id)
			return
		}
	}
}

func (wp *WorkerPool) AddWorker() {
	select {
	case <-wp.ctx.Done():
		fmt.Println("cannot add worker: pool is shutting down")
	default:
		wp.addChan <- struct{}{}
	}
}

func (wp *WorkerPool) DeleteWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if wp.workerCount <= 0 {
		fmt.Println("cannot delete worker: no workers")
		return
	}

	select {
	case <-wp.ctx.Done():
		fmt.Println("cannot delete worker: pool is shutting down")
	default:
		wp.deleteChan <- struct{}{}
	}
}

func (wp *WorkerPool) SendJob(data string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if wp.workerCount <= 0 {
		fmt.Println("cannot send job: no workers")
		return
	}

	select {
	case <-wp.ctx.Done():
		fmt.Println("cannot send job: pool is shutting down")
	default:
		wp.jobChan <- data
	}
}

func (wp *WorkerPool) GracefulShutdown() {
	wp.cancel()
	wp.wg.Wait()
	fmt.Println("worker pool shutdown complete")
}
