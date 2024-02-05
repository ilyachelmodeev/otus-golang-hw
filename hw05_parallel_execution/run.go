package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func limitedIncrement(m *sync.Mutex, errorsTotal *int, maxErrorLimit int) bool {
	m.Lock()
	defer m.Unlock()
	*errorsTotal++
	if maxErrorLimit > 0 {
		return *errorsTotal < maxErrorLimit
	}
	return true
}

func worker(w *sync.WaitGroup, m *sync.Mutex, in chan Task, errorsTotal *int, maxErrorLimit int) {
	defer w.Done()
	for task := range in {
		if err := task(); err != nil {
			if !limitedIncrement(m, errorsTotal, maxErrorLimit) {
				return
			}
		}
	}
}

func Run(tasks []Task, n, m int) error {
	maxErrorLimit := m
	if m <= 0 {
		maxErrorLimit = -1
	}

	mut := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(n) // no more than n tasks

	in := make(chan Task, len(tasks))
	for i := 0; i < len(tasks); i++ {
		in <- tasks[i]
	}
	close(in)

	var errorsTotal int
	for i := 0; i < n; i++ {
		go worker(wg, mut, in, &errorsTotal, maxErrorLimit)
	}
	wg.Wait()

	var result error
	if maxErrorLimit > 0 && errorsTotal >= maxErrorLimit {
		result = ErrErrorsLimitExceeded
	}
	return result
}
