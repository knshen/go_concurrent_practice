package main

import (
	"context"
	"fmt"
	"sync"
)

type Executor func() error

// goroutine池
type GoroutinePool struct {
	cap     int64
	workers chan bool
}

func NewGoroutinePool(cap int64) *GoroutinePool {
	return &GoroutinePool{
		cap:     cap,
		workers: make(chan bool, cap),
	}
}

func (p *GoroutinePool) Enter(ctx context.Context) {
	p.workers <- true
	// block here if workers are full
}

func (p *GoroutinePool) Leave(ctx context.Context) {
	<-p.workers
}

func (p *GoroutinePool) IsFull() bool {
	return len(p.workers) == int(p.cap)
}

func (p *GoroutinePool) AsyncRun(ctx context.Context, exec Executor) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recover:", err)
			}
		}()

		p.Enter(ctx)
		defer p.Leave(ctx)

		eErr := exec()
		if eErr != nil {
			fmt.Println("exec err:", eErr)
		}
	}()
}

func (p *GoroutinePool) RunWait(ctx context.Context, execs []Executor) {
	var wg sync.WaitGroup

	for _, exec := range execs {
		wg.Add(1)

		// 冗余代码
		go func() {
			defer func() {
				wg.Done()
				if err := recover(); err != nil {
					fmt.Println("recover:", err)
				}
			}()
			p.Enter(ctx)
			defer p.Leave(ctx)

			eErr := exec()
			if eErr != nil {
				fmt.Println("exec err:", eErr)
			}
		}()
	}

	wg.Wait()

}