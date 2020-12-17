package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"log"
	"golang.org/x/sync/errgroup"
)

func serveApp(ctx context.Context) error {
	s := &http.Server{Addr: "127.0.0.1:8081", Handler: http.DefaultServeMux}
	go func() {
		<-ctx.Done()
		shutdownctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		s.Shutdown(shutdownctx)
		log.Printf("Shutdown serveApp")
	}()
	log.Printf("Start serveApp")
	return s.ListenAndServe()
}

func serveDebug(ctx context.Context) error {
	sync := &http.Server{Addr: "127.0.0.1:8081", Handler: http.DefaultServeMux}
	go func() {
		<-ctx.Done()
		shutdownctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		s.Shutdown(shutdownctx)
		log.Printf("Shutdown serveDebug")
	}()
	log.Printf("Start serveDebug")
	return s.ListenAndServe()
}

func processSignal(ctx context.Context, ch chan os.Signal) error {

	//loop ch
	for {
		select {
		case <-ctx.Done():
			return nil
		case s := <-ch:
			log.Printf("get a signal %s", s.String())
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP:
				return errors.New("Close by signal " + s.String())
			case syscall.SIGHUP:
			default:
				return errors.New("Undefined signal")
			}
		}
	}
}

func main() {

	// Remark
	// done := make(chan error, 2)
	// stop := make(chan struct{})
	// go func() {
	// 	done <- serveDebug(stop)
	// }()
	// go func() {
	// 	done <- serveApp(stop)
	// }()
	// var stopped bool
	// for i :=0; i < cap(done); i++ {
	// 	if err := <-done; err !=nil {
	// 		fmt.Println("error: %v", err)
	// 	}
	// 	if !stopped {
	// 		stopped = true
	// 		close(stop)
	// 	}
	// }

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return serveApp(ctx)
	})
	g.Go(func() error {
		return serveDebug(ctx)
	})

	ch := make(chan os.Signal, 1)
	//将输入信号转发到ch
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)

	g.Go(func() error {
		return processSignal(ctx, ch)
	})
	if err := g.Wait(); err != nil {
		log.Printf("Server Error:%v\n", err)
	}
	
}