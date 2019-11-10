package main

import (
	"sync"
)

type SafeWaitGroup interface {
	Run(task func ())
	Wait()
}

type safeWaitGroupImpl struct {
	wg *sync.WaitGroup
}

func NewSafeWaitGroup() SafeWaitGroup {
	return &safeWaitGroupImpl{new(sync.WaitGroup)}
}

func (swg *safeWaitGroupImpl) Run(task func ()) {
	swg.wg.Add(1)
	go func() {
		task()
		swg.wg.Add(-1)
	}()
}

func (swg *safeWaitGroupImpl) Wait() {
	swg.wg.Wait()
}
