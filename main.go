package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k07g/g2/handler"
	"github.com/k07g/g2/infrastructure/inmemory"
	"github.com/k07g/g2/usecase"
)

func newServer(port string) *http.Server {
	repo := inmemory.NewTaskRepository()
	uc := usecase.NewTaskUseCase(repo)

	mux := http.NewServeMux()
	handler.NewTaskHandler(uc).Register(mux)

	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
}

func run(srv *http.Server, quit <-chan os.Signal) error {
	errCh := make(chan error, 1)
	go func() {
		log.Printf("server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-quit:
	}

	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("server stopped")
	return nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	if err := run(newServer(port), quit); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
