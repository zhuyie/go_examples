package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var counter int64

func testFunc(ctx context.Context, duration time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	t := time.NewTicker(duration)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			atomic.AddInt64(&counter, 1)
		}
	}
}

func usage() {
	fmt.Fprint(os.Stderr, "Usage:\n")
	flag.PrintDefaults()
}

func startDebugHTTPService() {
	go http.ListenAndServe("localhost:12345", nil)
}

func main() {
	flag.Usage = usage
	numGoroutines := flag.Int("num_goroutines", 1, "number of goroutines")
	durationInMS := flag.Int("duration", 1, "ticker duration in Milliseconds")
	flag.Parse()
	fmt.Printf("numGoroutines = %v\n", *numGoroutines)
	fmt.Printf("durationInMS  = %v\n", *durationInMS)

	startDebugHTTPService()

	var wg sync.WaitGroup
	ctx, quitF := context.WithCancel(context.Background())
	duration := time.Duration(*durationInMS) * time.Millisecond
	for i := 0; i < *numGoroutines; i++ {
		wg.Add(1)
		go testFunc(ctx, duration, &wg)
	}

	qc := make(chan os.Signal)
	signal.Notify(qc, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-qc:
		fmt.Printf("Caught signal %v to quit\n", s)
		quitF()
	}

	wg.Wait()

	fmt.Printf("counter = %v\n", counter)

	return
}
