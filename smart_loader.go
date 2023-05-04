package main

import (
	"fmt"
	"time"
)

// run a DAG concurrently with timeout

// SmartLoader is a wrapper
type SmartLoader struct {
	loaders []*Loader
	timeout time.Duration

	done chan bool
}

func NewSmartLoader() *SmartLoader {
	return &SmartLoader{
		loaders: make([]*Loader, 0),
	}
}

func (l *SmartLoader) SetTimeout(d time.Duration) {
	l.timeout = d
}

func (l *SmartLoader) AddLoaders(lds ...*Loader) {
	l.loaders = append(l.loaders, lds...)
}

func (l *SmartLoader) Execute() {
	l.done = make(chan bool, len(l.loaders))
	for _, ldr := range l.loaders {
		ldr.done = l.done
	}

	//t1 := time.Now()
	for i, ldr := range l.loaders {
		go func(i int, ldr *Loader) {
			ldr.Execute()
		}(i, ldr)
	}
	//t2 := time.Now()
	//fmt.Println("before select:", t2.Second()-t1.Second())

	finishCount := 0
	for finishCount < len(l.loaders) {
		select {
		case <-l.done:
			finishCount++
		case <-time.After(l.timeout):
			//t3 := time.Now()
			//fmt.Println("smartLoader timeout!", t3.Second()-t2.Second())
			return
		}
	}

	fmt.Println("smartLoader success")
}
