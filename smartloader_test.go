package main

import (
	"fmt"
	"testing"
	"time"
)

func TestSmartLoader(t *testing.T) {
	sl := NewSmartLoader()
	sl.SetTimeout(time.Second * 1)

	l1 := NewLoader(func() error {
		fmt.Println("loader1")
		time.Sleep(time.Second)
		return nil
	})
	l2 := NewLoader(func() error {
		fmt.Println("loader2")
		time.Sleep(time.Second * 1)
		return nil
	})
	l3 := NewLoader(func() error {
		fmt.Println("loader3")
		time.Sleep(time.Second * 2)
		return nil
	})

	sl.AddLoaders(l1, l2, l3)
	l3.AddPreLoader(l1)
	l3.AddPreLoader(l2)

	// l1 ->
	// l2 -> l3

	t0 := time.Now()
	sl.Execute()
	dur := time.Since(t0)
	fmt.Println(dur)
}
