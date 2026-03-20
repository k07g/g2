package main

import (
	"log"
	"net/http"
	"os"

	"github.com/k07g/g2/handler"
	"github.com/k07g/g2/infrastructure/inmemory"
	"github.com/k07g/g2/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	repo := inmemory.NewTaskRepository()
	uc := usecase.NewTaskUseCase(repo)

	mux := http.NewServeMux()
	h := handler.NewTaskHandler(uc)
	h.Register(mux)

	log.Printf("server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
