package router

import (
	"encoding/json"
	"net/http"
)

func (h *App) GetWorkers(w http.ResponseWriter, r *http.Request) {
	workers, err := h.Queries.SelectWorkers(r.Context())
	if err != nil {
		h.logError("Error fetching workers:", w, r, http.StatusInternalServerError, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workers)
	}

}
