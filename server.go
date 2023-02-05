package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

func initServer() (*http.Server, *sync.WaitGroup) {
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	return startHttpServer(httpServerExitDone), httpServerExitDone
}

func startHttpServer(wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: ":3000"}
	http.Handle("/", http.FileServer(http.Dir(".")))

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(fmt.Sprintf("ListenAndServe(): %v", err))
		}
	}()
	return srv
}

func closeServer(ctx context.Context, srv *http.Server, wg *sync.WaitGroup) {
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}

	wg.Wait()
}
