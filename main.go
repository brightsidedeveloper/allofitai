package main

import (
	"allofitai/handler"
	"allofitai/pkg/sb"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

//go:embed public
var FS embed.FS

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
		return
	}

	router := chi.NewMux()
	router.Use(handler.WithUser)

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))

	router.Get("/", handler.Make(handler.HandleHomeIndex))

	router.Get("/signin", handler.Make(handler.HandleAuthIndex))
	router.Post("/signin", handler.Make(handler.HandleSignInCreate))

	port := os.Getenv("PORT")
	slog.Info("Starting server", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return sb.Init()
}
