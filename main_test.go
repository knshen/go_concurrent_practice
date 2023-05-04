package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestAsynGo(t *testing.T) {
	// test go pool

	pool := NewGoroutinePool(10)
	ctx := context.Background()

	pool.AsyncRun(ctx, biz)

	time.Sleep(time.Hour)
}

func TestAsynGoFull(t *testing.T) {
	pool := NewGoroutinePool(3)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		pool.AsyncRun(ctx, biz)
	}

	time.Sleep(time.Hour)
}

func TestGoWait(t *testing.T) {
	pool := NewGoroutinePool(10)
	ctx := context.Background()

	pool.RunWait(ctx, []Executor{
		biz,
		biz,
		biz,
	})

	fmt.Println("main finish")
}

func TestGoWaitFull(t *testing.T) {
	pool := NewGoroutinePool(3)
	ctx := context.Background()

	pool.RunWait(ctx, []Executor{
		biz,
		biz,
		biz,
		biz,
		biz,
	})

	fmt.Println("main finish")
}

func biz() error {
	fmt.Println("biz ok")
	time.Sleep(time.Second * 3)
	return nil
}
