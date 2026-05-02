// this package is responsible for all CRUD operations on DB
package router

import (
	"log/slog"
	"net/http"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

type App struct {
	Queries *db.Queries
	Logger  *slog.Logger
}

func NewRouter(h *App) http.Handler {
	r := chi.NewRouter()

	r.Post("/login", h.Login)

	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)

		r.Post("/auth/change-password", h.ChangePassword)

		// pins
		r.Get("/pins", h.GetPins)                                // valid
		r.Get("/pin/{id}", h.GetPin)                             // valid
		r.Post("/pins", h.CreatePin)                             // valid
		r.Put("/pin/{id}", h.UpdatePin)                          // valid
		r.Delete("/pin/{id}", h.DeletePin)                       // valid
		r.Get("/pins/category/{category}", h.ListPinsByCategory) // valid

		// categories
		r.Get("/categories", h.ListCategories)             // valid
		r.Post("/categories", h.CreateCategory)            // valid
		r.Delete("/category/{category}", h.DeleteCategory) // valid
		r.Put("/category", h.UpdateCategory)               // valid (doesn't have category in url since the name is the identifier)

		// visits
		r.Post("/visits", h.CreateVisit)            // valid
		r.Delete("/visit/{id}", h.DeleteVisit)      // valid
		r.Get("/visits/pin/{id}", h.GetVisitsByPin) // valid
		r.Put("/visit", h.UpdateVisit)              // valid (doesn't have id in url since the id is in the body)

		// patrol plans
		r.Get("/patrol-plans", h.ListPatrolPlans)                              // valid
		r.Get("/patrol-plan/{id}", h.GetPatrolPlan)                            // valid
		r.Get("/patrol-plan/{id}/pins", h.GetPatrolPlanWithPins)               // valid
		r.Post("/patrol-plans", h.CreatePatrolPlan)                            // valid
		r.Delete("/patrol-plan/{id}", h.DeletePatrolPlan)                      // valid
		r.Post("/patrol-plan/{id}/pins", h.AddPinToPatrolPlan)                 //valid
		r.Delete("/patrol-plan/{id}/pins/{pin_id}", h.RemovePinFromPatrolPlan) // valid

		// stats
		r.Post("/stats", h.UpsertStreetworkStat)            // valid
		r.Get("/stats/month/{month}", h.GetStatsByMonth)    // valid
		r.Get("/stats/worker/{worker}", h.GetStatsByWorker) // valid
		r.Get("/all-stats", h.GetAllStats)                  // valid

		// app updates
		r.Get("/updates", h.ListAppUpdates)                        // not needed for now
		r.Get("/updates/latest", h.GetLatestAppUpdate)             // not needed for now
		r.Get("/updates/unviewed/{user_id}", h.GetUnviewedUpdates) // not needed for now
		r.Post("/updates/viewed", h.MarkUpdateViewed)              // not needed for now
	})

	return h.requestLogger(r)
}
