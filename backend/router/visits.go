package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *App) CreateVisit(w http.ResponseWriter, r *http.Request) {
	// description CAN BE empty

	var params db.CreateVisitParams
	var err error

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.logError("Invalid request body", w, r, http.StatusBadRequest, err)
		return
	}

	if params.PinID <= 0 {
		h.logError("Pin ID cannot be lower than 0", w, r, http.StatusBadRequest, nil)
		return
	}

	params.Name, err = h.sanitize(params.Name, CHARLIMIT)
	if err != nil {
		h.logError("Invalid name", w, r, http.StatusBadRequest, err)
		return
	}

	visit, err := h.Queries.CreateVisit(r.Context(), params)
	if err != nil {
		h.logError("Error creating visit", w, r, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(visit)
}

func (h *App) DeleteVisit(w http.ResponseWriter, r *http.Request) {
	strVisitID := chi.URLParam(r, "id")
	visitID, err := strconv.Atoi(strVisitID)
	if err != nil {
		h.logError("Invalid Visit ID", w, r, http.StatusBadRequest, err)
		return
	}

	if visitID < 0 {
		h.logError("Visit cannot be lower than 0", w, r, http.StatusBadRequest, err)
		return
	}

	err = h.Queries.DeleteVisit(r.Context(), int32(visitID))
	if err != nil {
		h.logError("Error deleting visit", w, r, http.StatusBadRequest, err)
		return
	}
}

func (h *App) GetVisitsByPin(w http.ResponseWriter, r *http.Request) {
	pin, err := h.paramIDtoInt(r)
	if err != nil {
		h.logError("Invalid Pin ID", w, r, http.StatusBadRequest, err)
		return
	}

	params, err := h.Queries.GetVisitsByPin(r.Context(), int32(pin))
	if err != nil {
		h.logError("Couldn't get visits", w, r, http.StatusBadRequest, err)
		return
	}

	json.NewEncoder(w).Encode(&params)
}

func (h *App) UpdateVisit(w http.ResponseWriter, r *http.Request) {
	var params db.UpdateVisitParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.logError("Invalid request body", w, r, http.StatusBadRequest, err)
		return
	}

	if params.ID <= 0 {
		h.logError("Visit ID cannot be lower than 0", w, r, http.StatusBadRequest, nil)
		return
	}

	if params.PinID <= 0 {
		h.logError("Pin ID cannot be lower than 0", w, r, http.StatusBadRequest, nil)
		return
	}

	if params.Name == "" {
		h.logError("Name cannot be empty", w, r, http.StatusBadRequest, nil)
		return
	}

	visit, err := h.Queries.UpdateVisit(r.Context(), params)
	if err != nil {
		h.logError("Couldn't update visit", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(visit)
}

func (h *App) GetVisitsByDateRange(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	start, err := h.sanitize(start, 10)
	if err != nil {
		h.logError("Invalid start date", w, r, http.StatusBadRequest, err)
		return
	}

	end, err = h.sanitize(end, 10)
	if err != nil {
		h.logError("Invalid end date", w, r, http.StatusBadRequest, err)
		return
	}

	if start == "" || end == "" {
		h.logError("start and end are required", w, r, http.StatusBadRequest, nil)
		return
	}

	// validate date format YYYY-MM-DD
	_, err = time.Parse("2006-01-02", start)
	if err != nil {
		h.logError("Invalid start date format, use YYYY-MM-DD", w, r, http.StatusBadRequest, err)
		return
	}

	_, err = time.Parse("2006-01-02", end)
	if err != nil {
		h.logError("Invalid end date format, use YYYY-MM-DD", w, r, http.StatusBadRequest, err)
		return
	}

	visits, err := h.Queries.GetVisitsByDateRange(r.Context(), db.GetVisitsByDateRangeParams{
		VisitedAt: pgtype.Timestamptz{Time: mustParseDate(start), Valid: true},
		VisitedAt_2:   pgtype.Timestamptz{Time: mustParseDateEnd(end), Valid: true},
	})
	if err != nil {
		h.logError("Failed to fetch visits", w, r, http.StatusInternalServerError, err)
		return
	}

	json.NewEncoder(w).Encode(&visits)
}

func mustParseDate(s string) time.Time {
    t, _ := time.Parse("2006-01-02", s)
    return t
}

func mustParseDateEnd(s string) time.Time {
    t, _ := time.Parse("2006-01-02", s)
    return t.Add(24*time.Hour - time.Second)
}