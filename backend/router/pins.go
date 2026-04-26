package router

import (
	"encoding/json"
	"net/http"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

func (h *App) CreatePin(w http.ResponseWriter, r *http.Request) {
	var params db.CreatePinParams
	if err := h.jsonDec.Decode(&params); err != nil {
		h.logError("Invalid pin parameters", w, http.StatusBadRequest, err)
		return
	}

	var err error
	params.Title, err = h.sanitize(params.Title, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param title", w, http.StatusBadRequest, err)
		return
	}

	params.Category, err = h.sanitize(params.Category, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param category", w, http.StatusBadRequest, err)
		return
	}

	pin, err := h.queries.CreatePin(r.Context(), params)
	if err != nil {
		h.logError("Couldn't create pin", w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(pin)
}

func (h *App) GetPins(w http.ResponseWriter, r *http.Request) {
	pins, err := h.queries.ListPins(r.Context())
	if err != nil {
		h.logger.Println("Error fetching pins:", err)
		return
	}
	h.logger.Println("Fetched pins:", pins)

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(pins)
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
	h.jsonEnc.Encode(pin)
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
	h.jsonEnc.Encode(pin)
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

func (h *App) ListPinsByCategory(w http.ResponseWriter, r *http.Request) {
	category, err := h.sanitize(chi.URLParam(r, "category"), CHARLIMIT)
	if err != nil {
		h.logError("Invalid category name", w, http.StatusBadRequest, err)
		return
	}

	pins, err := h.queries.ListPinsByCategory(r.Context(), category)
	if err != nil {
		h.logError("Couldn't fetch pins by category", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(pins)
}
