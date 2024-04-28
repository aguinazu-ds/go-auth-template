package main

import (
	"embed"
	"go-auth-template/db"
	"go-auth-template/handler"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

//go:embed public
var FS embed.FS

func main() {
	if err := initEverithing(); err != nil {
		panic(err)
	}

	router := chi.NewMux()

	// Routes
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/", handler.Make(handler.HandleHomeIndex))

	// Auth
	router.Get("/login", handler.Make(handler.HandleAuthLogin))
	router.Get("/signup", handler.Make(handler.HandleAuthSignup))
	router.Post("/signup", handler.Make(handler.HandleAuthSignupPost))

	port := os.Getenv("APP_PORT")
	slog.Info("Starting server on port " + port)
	http.ListenAndServe(port, router)
}

func initEverithing() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	return nil
}
