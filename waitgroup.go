package main

import (
	"sync"
)

type Spawner interface {
	Run(task func())
}

type SafeWaitGroup interface {
	Spawner
	Wait()
}

type safeWaitGroupImpl struct {
	wg *sync.WaitGroup
}

func NewSafeWaitGroup() SafeWaitGroup {
	return &safeWaitGroupImpl{new(sync.WaitGroup)}
}

func (swg *safeWaitGroupImpl) Run(task func()) {
	swg.wg.Add(1)
	go func() {
		task()
		swg.wg.Add(-1)
	}()
}

func (swg *safeWaitGroupImpl) Wait() {
	swg.wg.Wait()
}

func RunGroup(taskRunner func(Spawner)) {
	swg := NewSafeWaitGroup()
	taskRunner(swg)
	swg.Wait()
}
