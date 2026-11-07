package main

/*ROUGH LIST OF THINGS TO LEARN
- clean api design and project structure
- chi router
- DB connection and usage
- databases amd sql
- jwt authentication (header and cookie)
- logging and error handling
- basic redis caching
*/

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jusxtdev/cairn/internal/auth"
	"github.com/jusxtdev/cairn/internal/config"
	"github.com/jusxtdev/cairn/internal/httpx"
	"github.com/jusxtdev/cairn/internal/logger"
	"github.com/jusxtdev/cairn/internal/users"
)

func health(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSONResponse(w, http.StatusOK, httpx.Response{
		Status: true,
		Message: "API is healthy",
	})
}

func main(){
	// - - - initialize logging - - -  
	logfile, err := logger.InitLogger("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()

	slog.Info("application started")

	// - - - load the environment variables - - -  
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return
	}

	// - - - setup db - - -
	ctx := context.Background()
	slog.Debug("setting up db pool")
	db, err := pgxpool.New(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		slog.Error("db conenction failed", "error", err)
		return
	}
	defer db.Close()
	if err := db.Ping(ctx); err != nil {
		slog.Error("db conenction failed", "error", err)
		return
	}
	slog.Info("db connected successfully")

	// - - - setup router - - -
	r := chi.NewRouter()

	// modify the default WriteMethodNotAllowed
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteMethodNotAllowedResponse(w)
	})
	r.Use(middleware.Logger)
	r.Use(logger.LoggingMiddleware)

	// - - - initialize handlers - - -
	authHandler := auth.NewHandler(db, cfg)
	userHandler := users.NewHandler(db)

	// - - - register routes - - -
	// public routes
	r.Get("/health", health)
	r.Route("/auth", authHandler.Routes)

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return auth.AuthenticationMiddleware(authHandler.UserService, cfg, next)
		})
		r.Route("/users", userHandler.Routes)
	} )

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}