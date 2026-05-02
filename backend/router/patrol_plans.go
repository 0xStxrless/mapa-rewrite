package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

func (h *App) CreatePatrolPlan(w http.ResponseWriter, r *http.Request) {
	var params db.CreatePatrolPlanParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.logError("Invalid patrol plan parameters", w, r, http.StatusBadRequest, err)
		return
	}

	var err error
	params.Name, err = h.sanitize(params.Name, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param name", w, r, http.StatusBadRequest, err)
		return
	}

	params.Date, err = h.sanitize(params.Date, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param date", w, r, http.StatusBadRequest, err)
		return
	}

	plan, err := h.Queries.CreatePatrolPlan(r.Context(), params)
	if err != nil {
		h.logError("Couldn't create patrol plan", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

func (h *App) GetPatrolPlan(w http.ResponseWriter, r *http.Request) {
	planID, err := h.paramIDtoInt(r)
	if err != nil {
		h.logError("Invalid patrol plan ID", w, r, http.StatusBadRequest, err)
		return
	}

	plan, err := h.Queries.GetPatrolPlan(r.Context(), int32(planID))
	if err != nil {
		h.logError("Patrol plan not found", w, r, http.StatusNotFound, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

func (h *App) ListPatrolPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.Queries.ListPatrolPlans(r.Context())
	if err != nil {
		h.logError("Couldn't fetch patrol plans", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plans)
}

func (h *App) DeletePatrolPlan(w http.ResponseWriter, r *http.Request) {
	planID, err := h.paramIDtoInt(r)
	if err != nil {
		h.logError("Invalid patrol plan ID", w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.Queries.DeletePatrolPlan(r.Context(), int32(planID)); err != nil {
		h.logError("Couldn't delete patrol plan", w, r, http.StatusInternalServerError, err)
		return
	}
}

func (h *App) GetPatrolPlanWithPins(w http.ResponseWriter, r *http.Request) {
	planID, err := h.paramIDtoInt(r)
	if err != nil {
		h.logError("Invalid patrol plan ID", w, r, http.StatusBadRequest, err)
		return
	}

	rows, err := h.Queries.GetPatrolPlanWithPins(r.Context(), int32(planID))
	if err != nil {
		h.logError("Couldn't fetch patrol plan pins", w, r, http.StatusInternalServerError, err)
		return
	}
	if rows == nil {
		h.logError("Couldn't fetch any patrol plan for this pin", w, r, http.StatusBadRequest, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func (h *App) AddPinToPatrolPlan(w http.ResponseWriter, r *http.Request) {
	planID, err := h.paramIDtoInt(r)
	if err != nil {
		h.logError("Invalid patrol plan ID", w, r, http.StatusBadRequest, err)
		return
	}

	var params db.AddPinToPatrolPlanParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.logError("Invalid parameters", w, r, http.StatusBadRequest, err)
		return
	}

	if params.PinID <= 0 {
		h.logError("Pin ID cannot be lower than 0", w, r, http.StatusBadRequest, nil)
		return
	}

	// params.SortOrder this needs to be implemented

	params.PatrolPlanID = int32(planID)

	entry, err := h.Queries.AddPinToPatrolPlan(r.Context(), params)
	if err != nil {
		h.logError("Couldn't add pin to patrol plan", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

func (h *App) RemovePinFromPatrolPlan(w http.ResponseWriter, r *http.Request) {
	planID, err := h.paramIDtoInt(r)
	if err != nil {
		h.logError("Invalid patrol plan ID", w, r, http.StatusBadRequest, err)
		return
	}

	pinIDStr, err := h.sanitize(chi.URLParam(r, "pin_id"), CHARLIMIT)
	if err != nil {
		h.logError("Invalid pin ID", w, r, http.StatusBadRequest, err)
		return
	}

	pinID, err := strconv.Atoi(pinIDStr)
	if err != nil || pinID < 0 {
		h.logError("Invalid pin ID", w, r, http.StatusBadRequest, err)
		return
	}

	params := db.RemovePinFromPatrolPlanParams{
		PatrolPlanID: int32(planID),
		PinID:        int32(pinID),
	}

	if err := h.Queries.RemovePinFromPatrolPlan(r.Context(), params); err != nil {
		h.logError("Couldn't remove pin from patrol plan", w, r, http.StatusInternalServerError, err)
		return
	}
}
