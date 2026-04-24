package router

import (
	"net/http"

	"github.com/0xstxrless/punkt/backend/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.AppHandler) http.Handler {
	r := chi.NewRouter()

	// r.Post("/login/auth") // implement this when you get just all of mock data lol

	r.Group(func(r chi.Router) {
		// r.Use(authMiddleware) // Apply auth middleware to all routes in this group
		r.Get("/pins", h.GetPins)
		r.Get("/pin/{id}", h.GetPin)

	})
}
