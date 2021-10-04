package oslo

import (
	"runtime"
)

var (
	QueueSize = 1024
	queue     chan func()
	stop      chan struct{}
)

func Init() {
	if check() {
		return
	}

	runtime.LockOSThread()

	if queue == nil {
		queue = make(chan func(), QueueSize)
	}

	if stop == nil {
		stop = make(chan struct{})

		go func() {
			for {
				select {
				case f := <-queue:
					f()
				case <-stop:
					runtime.UnlockOSThread()
					return
				}
			}
		}()
	}
}

func check() bool {
	return queue != nil && stop != nil
}

func Stop() {
	if stop != nil {
		stop <- struct{}{}

		stop = nil
		queue = nil
	}
}

func Call(block bool, f func()) {
	if !check() {
		Init()
	}

	if block {
		done := make(chan struct{})

		queue <- func() {
			f()
			done <- struct{}{}
		}

		<-done

		return
	}

	queue <- f
}
