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

	router.Group(func(router chi.Router) {
		router.Get("/", handler.Make(handler.RenderHome))
		router.Get("/signin", handler.Make(handler.RenderSignIn))
		router.Get("/create", handler.Make(handler.RenderCreate))
	})

	router.Group(func(router chi.Router) {
		router.Use(handler.RequireAuth)
		router.Get("/settings", handler.Make(handler.RenderSettings))
	})

	router.Group(func(router chi.Router) {
		router.Get("/oauth/google", handler.Make(handler.SignInWithGoogle))
		router.Post("/signin", handler.Make(handler.SignIn))
		router.Post("/create", handler.Make(handler.Create))
		router.Post("/logout", handler.Make(handler.Logout))
		router.Get("/auth/callback", handler.Make(handler.AuthCallback))
	})

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
