package router

import (
	"encoding/json"
	"net/http"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

func (h *App) UpsertStreetworkStat(w http.ResponseWriter, r *http.Request) {
	var params db.UpsertStreetworkStatParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.logError("Invalid streetwork stat parameters", w, r, http.StatusBadRequest, err)
		return
	}

	var err error
	params.WorkerName, err = h.sanitize(params.WorkerName, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param worker_name", w, r, http.StatusBadRequest, err)
		return
	}

	params.Month, err = h.sanitize(params.Month, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param month", w, r, http.StatusBadRequest, err)
		return
	}

	stat, err := h.Queries.UpsertStreetworkStat(r.Context(), params)
	if err != nil {
		h.logError("Couldn't upsert streetwork stat", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stat)
}

func (h *App) GetStatsByMonth(w http.ResponseWriter, r *http.Request) {
	month, err := h.sanitize(chi.URLParam(r, "month"), CHARLIMIT)
	if err != nil {
		h.logError("Invalid month param", w, r, http.StatusBadRequest, err)
		return
	}

	stats, err := h.Queries.GetStatsByMonth(r.Context(), month)
	if err != nil {
		h.logError("Couldn't fetch stats", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *App) GetStatsByWorker(w http.ResponseWriter, r *http.Request) {
	worker, err := h.sanitize(chi.URLParam(r, "worker"), CHARLIMIT)
	if err != nil {
		h.logError("Invalid worker param", w, r, http.StatusBadRequest, err)
		return
	}

	stats, err := h.Queries.GetStatsByWorker(r.Context(), worker)
	if err != nil {
		h.logError("Couldn't fetch stats", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
