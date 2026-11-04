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
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jusxtdev/cairn/internal/auth"
	"github.com/jusxtdev/cairn/internal/config"
	"github.com/jusxtdev/cairn/internal/httpx"
	"github.com/jusxtdev/cairn/internal/user"
)

func health(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSONResponse(w, http.StatusOK, httpx.Response{
		Status: true,
		Message: "API is healthy",
	})
}

func main(){
	// - - - load the environment variables - - -  
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// - - - setup db - - -
	ctx := context.Background()
	db, err := pgxpool.New(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(ctx); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// - - - setup router - - -
	r := chi.NewRouter()

	// modify the default WriteMethodNotAllowed
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteMethodNotAllowedResponse(w)
	})
	r.Use(middleware.Logger)

	// - - - initialize handlers - - -
	authHandler := auth.NewHandler(db, cfg)
	userHandler := user.NewHandler()

	// - - - register routes - - -
	// public routes
	r.Get("/health", health)
	r.Route("/auth", authHandler.Routes)

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return auth.AuthenticationMiddleware(db, next)
		})
		r.Route("/users", userHandler.Routes)
	} )

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}