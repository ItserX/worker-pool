package pool

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestPoolCreation(t *testing.T) {
	pool := NewWorkerPool(3)
	defer pool.GracefulShutdown()

	timeout := time.After(100 * time.Millisecond)
	for {
		select {
		case <-timeout:
			t.Fatal("timeout waiting for workers to start")
		default:
			pool.mu.Lock()
			if pool.workerCount == 3 {
				pool.mu.Unlock()
				return
			}
			pool.mu.Unlock()
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func TestJobProcessing(t *testing.T) {
	pool := NewWorkerPool(1)
	defer pool.GracefulShutdown()

	var wg sync.WaitGroup
	wg.Add(10)

	done := make(chan struct{})

	for i := 0; i < 10; i++ {
		pool.SendJob(strconv.Itoa(i))
		wg.Done()
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for jobs to complete")
	}
}

func TestWorkerDeletion(t *testing.T) {
	pool := NewWorkerPool(2)
	defer pool.GracefulShutdown()

	pool.DeleteWorker()

	timeout := time.After(100 * time.Millisecond)
	for {
		select {
		case <-timeout:
			t.Fatal("timeout waiting for worker deletion")
		default:
			pool.mu.Lock()
			if pool.workerCount == 1 {
				pool.mu.Unlock()
				return
			}
			pool.mu.Unlock()
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func TestWorkerAddition(t *testing.T) {
	pool := NewWorkerPool(2)
	defer pool.GracefulShutdown()

	pool.AddWorker()

	timeout := time.After(100 * time.Millisecond)
	for {
		select {
		case <-timeout:
			t.Fatal("timeout waiting for worker addition")
		default:
			pool.mu.Lock()
			if pool.workerCount == 3 {
				pool.mu.Unlock()
				return
			}
			pool.mu.Unlock()
			time.Sleep(5 * time.Millisecond)
		}
	}
}
