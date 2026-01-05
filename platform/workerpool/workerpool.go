package workerpool

import "sync"

type Task func()

type WorkerPool struct {
	tasks      chan Task
	sems       sync.Map
	maxWorkers int
	once       sync.Once
}

func New(maxWorkers int, taskBuffer int) *WorkerPool {
	if taskBuffer <= 0 {
		taskBuffer = 10000
	}
	if maxWorkers <= 0 {
		maxWorkers = 100
	}

	return &WorkerPool{
		tasks:      make(chan Task, taskBuffer),
		maxWorkers: maxWorkers,
	}
}

func (p *WorkerPool) Start() {
	p.once.Do(func() {
		for i := 0; i < p.maxWorkers; i++ {
			go p.worker()
		}
	})
}

func (p *WorkerPool) SubmitWithKey(key string, maxConcurrentPerKey int, task Task) {
	if maxConcurrentPerKey < 1 {
		maxConcurrentPerKey = 1
	}

	semI, _ := p.sems.LoadOrStore(key, make(chan struct{}, maxConcurrentPerKey))
	sem := semI.(chan struct{})

	p.tasks <- func() {
		sem <- struct{}{}
		defer func() { <-sem }()

		defer func() {
			if r := recover(); r != nil {

			}
		}()

		task()
	}
}

func (p *WorkerPool) Submit(task Task) {
	p.tasks <- task
}

func (p *WorkerPool) worker() {
	for task := range p.tasks {
		task()
	}
}

func (p *WorkerPool) Stop() {
	close(p.tasks)
}
