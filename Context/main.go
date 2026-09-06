package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker: Received cancellation signal, exiting...")
			return
		default:
			fmt.Println("Worker: Working...")
			time.Sleep(1 * time.Second)
		}
	}
}

func withCancel() {
	ctx, cancel := context.WithCancel(context.Background())

	go worker(ctx)

	time.Sleep(3 * time.Second)
	cancel() // stop worker

	time.Sleep(1 * time.Second)
}

func withTimeout() {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()

	go worker(ctx)

	time.Sleep(4 * time.Second)

}

func withParent() {
	parent, cancel := context.WithCancel(context.Background())

	child := context.WithValue(parent, "user", "harshal")

	go worker(parent)
	go worker(child)

	time.Sleep(2 * time.Second)
	cancel() // stops BOTH

	time.Sleep(1 * time.Second)

}

func workerWithTicker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker with Ticker: Received cancellation signal, exiting...")
			return
		case t := <-ticker.C:
			fmt.Println("Worker with Ticker: Tick at", t)
		}
	}
}

func withTicker() {

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	go workerWithTicker(ctx)
	time.Sleep(4 * time.Second)

}

func main() {
	//withCancel()
	//withTimeout()
	withTicker()
}
