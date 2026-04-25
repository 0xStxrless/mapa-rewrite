package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

func (h *App) CreateVisit(w http.ResponseWriter, r *http.Request) {
	var params db.CreateVisitParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	visit, err := h.queries.CreateVisit(r.Context(), params)
	if err != nil {
		h.logError("Error creating visit", w, r, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(visit)
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

	h.queries.DeleteVisit(r.Context(), int32(visitID))
}

// TODO: add update visit
