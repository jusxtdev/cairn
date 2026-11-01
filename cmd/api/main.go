package main

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main(){
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", health)

	http.ListenAndServe(":8080", r)
}