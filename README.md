# worker-pool
A flexible worker pool implementation that allows dynamic scaling of workers and graceful shutdown.

## Features

- Create a pool with an initial number of workers
- Dynamically add workers at runtime
- Remove workers at runtime
- Send jobs to the pool for processing
- Graceful shutdown that completes all pending jobs

## Installation

```sh
go get github.com/ItserX/worker-pool
```

## Example
```
func main() {
  // create a pool with 3 workers
	wp := pool.NewWorkerPool(3)

  // send 10 tasks
	for i := 1; i <= 10; i++ {
		task := fmt.Sprintf("task-%d", i)
		wp.SendJob(task)
	}

// add worker
	wp.AddWorker()

// send 10 tasks
	for i := 1; i <= 10; i++ {
		task := fmt.Sprintf("task-%d", i)
		wp.SendJob(task)
	}

  // delete 2 workers
	wp.DeleteWorker()
	wp.DeleteWorker()

  // send 10 tasks
	for i := 1; i <= 10; i++ {
		task := fmt.Sprintf("task-%d", i)
		wp.SendJob(task)
	}

  // gracefully shutdown
	wp.GracefulShutdown()
}
```

## API

### `NewWorkerPool(numWorkers int) *WorkerPool`  
Creates a new worker pool with the specified number of workers.

### `AddWorker()`  
Adds a new worker to the pool.

### `DeleteWorker()`  
Removes a worker from the pool. 

### `SendJob(data string)`  
Sends a job to the pool for processing.

### `GracefulShutdown()`  
Initiates a graceful shutdown.

## Testing
```sh
go test -v ./pkg/pool
```

## Run 
```sh
go run cmd/main.go
```
