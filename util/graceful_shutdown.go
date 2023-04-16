package util

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Operation func(ctx context.Context) error

func GracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]Operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		fmt.Println("Shutting down...")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			fmt.Println("Timeout", timeout.Milliseconds(), "ms has been elapsed, force exit")
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				fmt.Printf("Cleaning up: %s\n", innerKey)
				if err := innerOp(ctx); err != nil {
					fmt.Println(innerKey, ": Clean up failed: ", err.Error())
					return
				}

				fmt.Println(innerKey, "was shutdown gracefully...")
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
