package router

import (
	"net/http"
	"strconv"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

func (h *App) ListAppUpdates(w http.ResponseWriter, r *http.Request) {
	updates, err := h.queries.ListAppUpdates(r.Context())
	if err != nil {
		h.logError("Couldn't fetch app updates", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(updates)
}

func (h *App) GetLatestAppUpdate(w http.ResponseWriter, r *http.Request) {
	update, err := h.queries.GetLatestAppUpdate(r.Context())
	if err != nil {
		h.logError("Couldn't fetch latest app update", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(update)
}

func (h *App) GetUnviewedUpdates(w http.ResponseWriter, r *http.Request) {
	userIDStr, err := h.sanitize(chi.URLParam(r, "user_id"), CHARLIMIT)
	if err != nil {
		h.logError("Invalid user ID", w, http.StatusBadRequest, err)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID < 0 {
		h.logError("Invalid user ID", w, http.StatusBadRequest, err)
		return
	}

	updates, err := h.queries.GetUnviewedUpdates(r.Context(), int32(userID))
	if err != nil {
		h.logError("Couldn't fetch unviewed updates", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(updates)
}

func (h *App) MarkUpdateViewed(w http.ResponseWriter, r *http.Request) {
	var params db.MarkUpdateViewedParams
	if err := h.jsonDec.Decode(&params); err != nil {
		h.logError("Invalid parameters", w, http.StatusBadRequest, err)
		return
	}

	entry, err := h.queries.MarkUpdateViewed(r.Context(), params)
	if err != nil {
		h.logError("Couldn't mark update as viewed", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(entry)
}
