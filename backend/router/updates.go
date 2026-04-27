package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

func (h *App) ListAppUpdates(w http.ResponseWriter, r *http.Request) {
	updates, err := h.Queries.ListAppUpdates(r.Context())
	if err != nil {
		h.logError("Couldn't fetch app updates", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updates)
}

func (h *App) GetLatestAppUpdate(w http.ResponseWriter, r *http.Request) {
	update, err := h.Queries.GetLatestAppUpdate(r.Context())
	if err != nil {
		h.logError("Couldn't fetch latest app update", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(update)
}

func (h *App) GetUnviewedUpdates(w http.ResponseWriter, r *http.Request) {
	userIDStr, err := h.sanitize(chi.URLParam(r, "user_id"), CHARLIMIT)
	if err != nil {
		h.logError("Invalid user ID", w, r, http.StatusBadRequest, err)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID < 0 {
		h.logError("Invalid user ID", w, r, http.StatusBadRequest, err)
		return
	}

	updates, err := h.Queries.GetUnviewedUpdates(r.Context(), int32(userID))
	if err != nil {
		h.logError("Couldn't fetch un viewed updates", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updates)
}

func (h *App) MarkUpdateViewed(w http.ResponseWriter, r *http.Request) {
	var params db.MarkUpdateViewedParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.logError("Invalid parameters", w, r, http.StatusBadRequest, err)
		return
	}

	entry, err := h.Queries.MarkUpdateViewed(r.Context(), params)
	if err != nil {
		h.logError("Couldn't mark update as viewed", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}
