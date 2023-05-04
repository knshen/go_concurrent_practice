package main

import "sync"

type Loader struct {
	wg         *sync.WaitGroup
	downStream []*sync.WaitGroup // 下游依赖loader的wg

	prevLoaders []*Loader
	done        chan bool
	exec        Executor
}

func NewLoader(exec Executor) *Loader {
	return &Loader{
		downStream:  make([]*sync.WaitGroup, 0),
		prevLoaders: make([]*Loader, 0),
		exec:        exec,
	}
}
func (l *Loader) AddPreLoader(ldr *Loader) {
	l.prevLoaders = append(l.prevLoaders, ldr)
	if l.wg == nil {
		l.wg = &sync.WaitGroup{}
	}
	l.wg.Add(1)

	ldr.downStream = append(ldr.downStream, l.wg)
}

func (l *Loader) Execute() {
	defer func() {
		// 上游执行完成，释放下游loader的wg
		for _, wg := range l.downStream {
			wg.Done()
		}
		l.done <- true
	}()

	// 等待上游全部执行完成
	if len(l.prevLoaders) > 0 {
		l.wg.Wait()
	}

	// 执行本loader
	l.exec()

}
