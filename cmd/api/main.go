package main

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
	"github.com/jusxtdev/cairn/internal/httpx"
)

func health(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSONResponse(w, http.StatusOK, httpx.Response{
		Status: true,
		Message: "API is healthy",
	})
}

func main(){
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", health)

	http.ListenAndServe(":8080", r)
}