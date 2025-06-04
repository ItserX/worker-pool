package main

import (
	"fmt"

	"github.com/ItserX/worker-pool/pkg/pool"
)

func main() {
	wp := pool.NewWorkerPool(3)

	for i := 1; i <= 10; i++ {
		task := fmt.Sprintf("task-%d", i)
		wp.SendJob(task)
	}

	wp.AddWorker()

	for i := 1; i <= 10; i++ {
		task := fmt.Sprintf("task-%d", i)
		wp.SendJob(task)
	}

	wp.DeleteWorker()
	wp.DeleteWorker()

	for i := 1; i <= 10; i++ {
		task := fmt.Sprintf("task-%d", i)
		wp.SendJob(task)
	}

	wp.GracefulShutdown()
}
