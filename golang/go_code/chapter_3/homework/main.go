package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

//1. 基于 errgroup 实现一个 http server 的启动和关闭 ，
//以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	ctx, done := context.WithCancel(context.Background())
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(getUserServer(gCtx, &wg))
	g.Go(getOrderServer(gCtx, &wg))
	g.Go(getProductServer(gCtx, &wg))

	go func() {
		<-gCtx.Done()
		done()
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		done()
	}()

	if err := g.Wait(); err != nil {
		fmt.Printf("origin eror: %T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack track: \n%+v\n", err)
		os.Exit(1)
	}
	fmt.Println("everything closed successfully")
}

func getUserServer(ctx context.Context, wg *sync.WaitGroup) func() error {
	return func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Duration(2 * time.Second))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, world!"))
		})

		server := &http.Server{Addr: ":6000", Handler: mux}
		errChan := make(chan error, 1)

		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			if err := server.Shutdown(shutCtx); err != nil {
				errChan <- errors.Wrap(err, "error shutting down the user server")
			}
			fmt.Println("the user server is closed")
			close(errChan)
			wg.Done()
		}()

		fmt.Println("the user server is starting")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return errors.Wrap(err, "err start in User Server")
		}
		fmt.Println("the user server is closing")
		err := <-errChan
		wg.Wait()
		return err
	}
}

func getOrderServer(ctx context.Context, wg *sync.WaitGroup) func() error {
	return func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("order, Hello, world!"))
		})
		server := &http.Server{Addr: ":7000", Handler: mux}
		errChan := make(chan error, 1)

		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(shutCtx); err != nil {
				errChan <- errors.Wrap(err, "error shutting down the hello order server")
			}
			fmt.Println("the order server is closed")
			close(errChan)
			wg.Done()
		}()

		fmt.Println("the order server is starting")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return errors.Wrap(err, "err start in order Server")
		}
		fmt.Println("the order server is closing")
		err := <-errChan
		wg.Wait()
		return err
	}
}

func getProductServer(ctx context.Context, wg *sync.WaitGroup) func() error {
	return func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`product, Hello, world!`))
		})
		server := &http.Server{Addr: ":8000", Handler: mux}
		errChan := make(chan error, 1)

		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(shutCtx); err != nil {
				errChan <- errors.Wrap(err, "error shutting down the hello product server")
			}
			fmt.Println("the product server is closed")
			close(errChan)
			wg.Done()
		}()

		fmt.Println("the product server is starting")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return errors.Wrap(err, "err start in order Server")
		}
		fmt.Println("the product server is closing")
		err := <-errChan
		wg.Wait()
		return err
	}
}
