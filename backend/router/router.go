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

		// pins
		r.Get("/pins", h.GetPins)
		r.Get("/pin/{id}", h.GetPin)
		r.Post("/pins", h.CreatePin)
		r.Put("/pin/{id}", h.UpdatePin)
		r.Delete("/pin/{id}", h.DeletePin)
		r.Get("/pins/category/{category}", h.ListPinsByCategory)

		// categories
		r.Get("/categories", h.ListCategories)
		r.Post("/categories", h.CreateCategory)
		r.Delete("/category/{category}", h.DeleteCategory)
		r.Put("/category/{category}", h.UpdateCategory)

		// visits
		r.Post("/visits", h.CreateVisit)
		r.Delete("/visit/{id}", h.DeleteVisit)
		r.Get("/visits/pin/{id}", h.GetVisitsByPin)
		r.Put("/visit/{id}", h.UpdateVisit)

		// patrol plans
		r.Get("/patrol-plans", h.ListPatrolPlans)
		r.Get("/patrol-plan/{id}", h.GetPatrolPlan)
		r.Get("/patrol-plan/{id}/pins", h.GetPatrolPlanWithPins)
		r.Post("/patrol-plans", h.CreatePatrolPlan)
		r.Delete("/patrol-plan/{id}", h.DeletePatrolPlan)
		r.Post("/patrol-plan/{id}/pins", h.AddPinToPatrolPlan)
		r.Delete("/patrol-plan/{id}/pins/{pin_id}", h.RemovePinFromPatrolPlan)

		// stats
		r.Post("/stats", h.UpsertStreetworkStat)
		r.Get("/stats/month/{month}", h.GetStatsByMonth)
		r.Get("/stats/worker/{worker}", h.GetStatsByWorker)

		// app updates
		r.Get("/updates", h.ListAppUpdates)
		r.Get("/updates/latest", h.GetLatestAppUpdate)
		r.Get("/updates/unviewed/{user_id}", h.GetUnviewedUpdates)
		r.Post("/updates/viewed", h.MarkUpdateViewed)
	})

	return r
}
