package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (a *AppConfig) routes() http.Handler {
	// create mux
	mux := chi.NewRouter()

	// create cors middleware
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"POST", "GET", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// create middleware for testing back-end
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/handle", a.HandleAuthentication)

	// create get request for get serve broker
	mux.Post("/serve-broker", a.ServeBroker)

	// return mux
	return mux
}
