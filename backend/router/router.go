// this package is responsible for all CRUD operations on DB
package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

type App struct {
	queries *db.Queries
	logger  log.Logger
	jsonEnc json.Encoder
	jsonDec json.Decoder
}

func NewRouter(h *App) http.Handler {
	r := chi.NewRouter()

	// r.Post("/login/auth") // implement this when you get just all of mock data lol

	r.Group(func(r chi.Router) {
		// r.Use(authMiddleware) // Apply auth middleware to all routes in this group
		r.Get("/pins", h.GetPins)
		r.Get("/pin/{id}", h.GetPin)
	})

	return r
}
