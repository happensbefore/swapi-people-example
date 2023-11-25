package workerpool

import "sync"

type backgroundJob func()

type WorkerPool struct {
	jobs []backgroundJob
}

func New() *WorkerPool {
	return &WorkerPool{
		jobs: []backgroundJob{},
	}
}

func (pool *WorkerPool) AddBackgroundJob(job backgroundJob) {
	pool.jobs = append(pool.jobs, job)
}

func (pool *WorkerPool) Run() {
	var wg sync.WaitGroup

	for _, job := range pool.jobs {
		wg.Add(1)

		go func(job backgroundJob) {
			defer wg.Done()

			job()
		}(job)
	}

	wg.Wait()
}
