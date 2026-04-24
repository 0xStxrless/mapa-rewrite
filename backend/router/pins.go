package router

import (
	"encoding/json"
	"net/http"

	"github.com/0xstxrless/punkt/backend/internal/db"
)

func (h *App) GetPins(w http.ResponseWriter, r *http.Request) {
	pins, err := h.queries.ListPins(r.Context())
	if err != nil {
		h.logger.Println("Error fetching pins:", err)
		return
	}
	h.logger.Println("Fetched pins:", pins)
}

func (h *App) GetPin(w http.ResponseWriter, r *http.Request) {
	pinID, err := h.paramIDtoInt(r)
	if err != nil {
		http.Error(w, "Invalid pin ID", http.StatusNotFound)
		return
	}

	pin, err := h.queries.GetPin(r.Context(), int32(pinID))
	if err != nil {
		h.logger.Printf("Error fetching pin %d: %v\n", pinID, err)
		http.Error(w, "Pin not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pin)
}

func (h *App) UpdatePin(w http.ResponseWriter, r *http.Request) {
	pinID, err := h.paramIDtoInt(r)
	if err != nil {
		http.Error(w, "Invalid pin ID", http.StatusNotFound)
		return
	}

	var params db.UpdatePinParams
	params.ID = int32(pinID)

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pin, err := h.queries.UpdatePin(r.Context(), params)
	if err != nil {
		h.logger.Printf("Error updating pin %d: %v\n", pinID, err)
		http.Error(w, "Pin not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pin)
}

func (h *App) DeletePin(w http.ResponseWriter, r *http.Request) {
	pinID, err := h.paramIDtoInt(r)
	if err != nil {
		http.Error(w, "Invalid pin ID", http.StatusNotFound)
		return
	}

	if err := h.queries.DeletePin(r.Context(), int32(pinID)); err != nil {
		http.Error(w, "Pin Couldn't be deleted", http.StatusBadRequest)
		return
	}

}
