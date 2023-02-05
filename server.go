package main

import (
	"context"
	"net/http"
	"sync"
)

func initServer() (*http.Server, *sync.WaitGroup, error) {
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv, err := startHttpServer(httpServerExitDone)
	return srv, httpServerExitDone, err
}

func startHttpServer(wg *sync.WaitGroup) (*http.Server, error) {
	srv := &http.Server{Addr: ":3000"}
	http.Handle("/", http.FileServer(http.Dir(".")))

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return srv, nil
}

func closeServer(ctx context.Context, srv *http.Server, wg *sync.WaitGroup) error {
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	wg.Wait()

	return nil
}
