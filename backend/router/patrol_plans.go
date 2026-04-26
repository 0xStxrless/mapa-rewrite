package router

import (
	"net/http"
	"strconv"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

func (h *App) CreatePatrolPlan(w http.ResponseWriter, r *http.Request) {
	var params db.CreatePatrolPlanParams
	if err := h.jsonDec.Decode(&params); err != nil {
		h.logError("Invalid patrol plan parameters", w, http.StatusBadRequest, err)
		return
	}

	var err error
	params.Name, err = h.sanitize(params.Name, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param name", w, http.StatusBadRequest, err)
		return
	}

	params.Date, err = h.sanitize(params.Date, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param date", w, http.StatusBadRequest, err)
		return
	}

	plan, err := h.queries.CreatePatrolPlan(r.Context(), params)
	if err != nil {
		h.logError("Couldn't create patrol plan", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(plan)
}

func (h *App) GetPatrolPlan(w http.ResponseWriter, r *http.Request) {
	planID, err := h.paramIDtoInt(r)
	if err != nil {
		h.logError("Invalid patrol plan ID", w, http.StatusBadRequest, err)
		return
	}

	plan, err := h.queries.GetPatrolPlan(r.Context(), int32(planID))
	if err != nil {
		h.logError("Patrol plan not found", w, http.StatusNotFound, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(plan)
}

func (h *App) ListPatrolPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.queries.ListPatrolPlans(r.Context())
	if err != nil {
		h.logError("Couldn't fetch patrol plans", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(plans)
}

func (h *App) DeletePatrolPlan(w http.ResponseWriter, r *http.Request) {
	planID, err := h.paramIDtoInt(r)
	if err != nil {
		h.logError("Invalid patrol plan ID", w, http.StatusBadRequest, err)
		return
	}

	if err := h.queries.DeletePatrolPlan(r.Context(), int32(planID)); err != nil {
		h.logError("Couldn't delete patrol plan", w, http.StatusInternalServerError, err)
		return
	}
}

func (h *App) GetPatrolPlanWithPins(w http.ResponseWriter, r *http.Request) {
	planID, err := h.paramIDtoInt(r)
	if err != nil {
		h.logError("Invalid patrol plan ID", w, http.StatusBadRequest, err)
		return
	}

	rows, err := h.queries.GetPatrolPlanWithPins(r.Context(), int32(planID))
	if err != nil {
		h.logError("Couldn't fetch patrol plan pins", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(rows)
}

func (h *App) AddPinToPatrolPlan(w http.ResponseWriter, r *http.Request) {
	planID, err := h.paramIDtoInt(r)
	if err != nil {
		h.logError("Invalid patrol plan ID", w, http.StatusBadRequest, err)
		return
	}

	var params db.AddPinToPatrolPlanParams
	if err := h.jsonDec.Decode(&params); err != nil {
		h.logError("Invalid parameters", w, http.StatusBadRequest, err)
		return
	}
	params.PatrolPlanID = int32(planID)

	entry, err := h.queries.AddPinToPatrolPlan(r.Context(), params)
	if err != nil {
		h.logError("Couldn't add pin to patrol plan", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.jsonEnc.Encode(entry)
}

func (h *App) RemovePinFromPatrolPlan(w http.ResponseWriter, r *http.Request) {
	planID, err := h.paramIDtoInt(r)
	if err != nil {
		h.logError("Invalid patrol plan ID", w, http.StatusBadRequest, err)
		return
	}

	pinIDStr, err := h.sanitize(chi.URLParam(r, "pin_id"), CHARLIMIT)
	if err != nil {
		h.logError("Invalid pin ID", w, http.StatusBadRequest, err)
		return
	}

	pinID, err := strconv.Atoi(pinIDStr)
	if err != nil || pinID < 0 {
		h.logError("Invalid pin ID", w, http.StatusBadRequest, err)
		return
	}

	params := db.RemovePinFromPatrolPlanParams{
		PatrolPlanID: int32(planID),
		PinID:        int32(pinID),
	}

	if err := h.queries.RemovePinFromPatrolPlan(r.Context(), params); err != nil {
		h.logError("Couldn't remove pin from patrol plan", w, http.StatusInternalServerError, err)
		return
	}
}
