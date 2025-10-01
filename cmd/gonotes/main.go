package main

import (
	"log/slog"
	"net/http"
	"os"

	"gonotes/internal/auth/authhandler"
	"gonotes/internal/auth/storage/authsqlite"
	"gonotes/internal/config"
	"gonotes/internal/lib/logger"
	"gonotes/internal/middleware"
	"gonotes/internal/notes/noteshandler"
	"gonotes/internal/notes/storage/notessqlite"
	"gonotes/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.New(cfg.Env)

	router := chi.NewRouter()

	db, err := storage.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.Any("err", err))
		os.Exit(1)
	}

	notesRepository := notessqlite.NewNoteRepository(db)
	userRepository := authsqlite.NewUserRepository(db)

	noteshandler := noteshandler.NewHandler(log, notesRepository)
	authhandler := authhandler.NewHandler(log, userRepository, cfg.SecretKey)

	authMW := middleware.NewAuthMiddleware(cfg.SecretKey)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:*", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Post("/auth/register", authhandler.Create)
	router.Post("/auth/login", authhandler.Login)

	router.Route("/notes", func(r chi.Router) {
		r.Use(authMW.Auth)
		r.Post("/", noteshandler.Create)
		r.Get("/{id}", noteshandler.Get)
		r.Delete("/{id}", noteshandler.Delete)
		r.Get("/", noteshandler.GetAll)
	})

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("starting server", "port", cfg.HTTPServer.Address)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server error", slog.Any("err", err))
	}

	log.Info("server stopped")
}
