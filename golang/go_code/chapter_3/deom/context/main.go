package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	//simple1()
	//simple2()
	simple3()
}

func simple3() {
	ctx, done := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	defer done()

	select {
	case <-ctx.Done():
		fmt.Println("Done")
		return
	case <-time.After(2 * time.Second):
		fmt.Println("over")
	}
}

func simple() {
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("chan return")
					return
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			fmt.Println("break")
			break
		}
	}
}

func simple1() {
	t := time.Now().Add(1 * time.Microsecond)
	ctx, cancel := context.WithDeadline(context.Background(), t)
	defer cancel()

	select {
	case <-time.After(10 * time.Microsecond):
		fmt.Println("oversleep")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}

func simple2() {

	ctx, done := context.WithCancel(context.Background())
	g, gctx := errgroup.WithContext(ctx)

	// goroutine to check for signals to gracefully finish all functions
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			fmt.Printf("Received signal: %s\n", sig)
			done()
		case <-gctx.Done():
			fmt.Printf("closing signal goroutine\n")
			return gctx.Err()
		}

		return nil
	})

	// just a ticker every 2s
	g.Go(func() error {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Printf("ticker 2s ticked\n")
				// testcase what happens if an error occured
				//return fmt.Errorf("test error ticker 2s")
			case <-gctx.Done():
				fmt.Printf("closing ticker 2s goroutine\n")
				return gctx.Err()
			}
		}
	})

	// just a ticker every 1s
	g.Go(func() error {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Printf("ticker 1s ticked\n")
			case <-gctx.Done():
				fmt.Printf("closing ticker 1s goroutine\n")
				return gctx.Err()
			}
		}
	})

	// force a stop after 60s
	time.AfterFunc(20*time.Second, func() {
		fmt.Printf("force finished after 60s")
		done()
	})

	// wait for all errgroup goroutines
	err := g.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Print("context was canceled")
		} else {
			fmt.Printf("received error: %v", err)
		}
	} else {
		fmt.Println("finished clean")
	}
}
