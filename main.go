package main

import (
	"embed"
	"go-auth-template/db"
	"go-auth-template/handler"
	"go-auth-template/pkg/authsession"
	"go-auth-template/pkg/mailer"
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
	router.Use(handler.WithUser)

	// Routes
	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/", handler.Make(handler.HandleHomeIndex))
	router.Get("/health-check", handler.Make(handler.HandleHealthCheck))

	// Auth
	router.Get("/login", handler.Make(handler.HandleAuthLogin))
	router.Post("/login", handler.Make(handler.HandleAuthLoginPost))
	router.Get("/signup", handler.Make(handler.HandleAuthSignup))
	router.Post("/signup", handler.Make(handler.HandleAuthSignupPost))
	router.Get("/activate", handler.Make(handler.HandleAuthActivate))
	router.Post("/logout", handler.Make(handler.HandleAuthLogoutPost))

	// Protected routes
	router.Group(func(c chi.Router) {
		c.Use(handler.WithAuth)
		c.Get("/account", handler.Make(handler.HandleUserAccount))
	})

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
	if err := mailer.Init(); err != nil {
		return err
	}
	authsession.Init()
	return nil
}
