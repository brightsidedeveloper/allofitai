package main

import (
	"allofitai/handler"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
		return
	}

	router := chi.NewMux()

	router.Get("/", handler.HandleHomeIndex)

	port := os.Getenv("PORT")
	slog.Info("Starting server", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	return godotenv.Load()
}
