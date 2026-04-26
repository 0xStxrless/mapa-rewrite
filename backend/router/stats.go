package router

import (
	"net/http"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

func (h *App) UpsertStreetworkStat(w http.ResponseWriter, r *http.Request) {
	var params db.UpsertStreetworkStatParams
	if err := h.jsonDec.Decode(&params); err != nil {
		h.logError("Invalid streetwork stat parameters", w, http.StatusBadRequest, err)
		return
	}

	var err error
	params.WorkerName, err = h.sanitize(params.WorkerName, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param worker_name", w, http.StatusBadRequest, err)
		return
	}

	params.Month, err = h.sanitize(params.Month, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param month", w, http.StatusBadRequest, err)
		return
	}

	stat, err := h.queries.UpsertStreetworkStat(r.Context(), params)
	if err != nil {
		h.logError("Couldn't upsert streetwork stat", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(stat)
}

func (h *App) GetStatsByMonth(w http.ResponseWriter, r *http.Request) {
	month, err := h.sanitize(chi.URLParam(r, "month"), CHARLIMIT)
	if err != nil {
		h.logError("Invalid month param", w, http.StatusBadRequest, err)
		return
	}

	stats, err := h.queries.GetStatsByMonth(r.Context(), month)
	if err != nil {
		h.logError("Couldn't fetch stats", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(stats)
}

func (h *App) GetStatsByWorker(w http.ResponseWriter, r *http.Request) {
	worker, err := h.sanitize(chi.URLParam(r, "worker"), CHARLIMIT)
	if err != nil {
		h.logError("Invalid worker param", w, http.StatusBadRequest, err)
		return
	}

	stats, err := h.queries.GetStatsByWorker(r.Context(), worker)
	if err != nil {
		h.logError("Couldn't fetch stats", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(stats)
}
