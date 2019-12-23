package main

type Task struct {
	f    func(...interface{})
	args []interface{}
}

func NewTask(f func(...interface{}), args ...interface{}) *Task {
	return &Task{
		f:    f,
		args: args,
	}
}

func (t *Task) Execute() {
	t.f(t.args)
}

type Pool struct {
	EntryChan chan *Task
	JobsChan  chan *Task
	WorkerNum int
}

func NewPool(workerNum int) *Pool {
	return &Pool{
		EntryChan: make(chan *Task, 50),
		JobsChan:  make(chan *Task, 50),
		WorkerNum: workerNum,
	}
}

func (p *Pool) Submit(t *Task) {
	p.EntryChan <- t
}

// Worker 协程池创建一个 Worker
func (p *Pool) Worker(id int) {
	// 永久从 JobsChan 中取任务，并执行
	for task := range p.JobsChan {
		task.Execute()
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.WorkerNum; i++ {
		go p.Worker(i)
	}
	for task := range p.EntryChan {
		p.JobsChan <- task
	}
}
