package router

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *App) paramIDtoInt(w http.ResponseWriter, r *http.Request) (int, error) {
	pinIDStr := chi.URLParam(r, "id")
	pinID, err := strconv.Atoi(pinIDStr)
	if err != nil {
		return 0, err
	}

	if pinID < 0 {
		return 0, err
	}
}
