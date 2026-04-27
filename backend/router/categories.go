package router

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/go-chi/chi/v5"
)

const CHARLIMIT = 64

func (h *App) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Queries.ListCategories(r.Context())
	if err != nil {
		h.logError("Couldn't fetch categories", w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *App) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var params db.CreateCategoryParams
	var err error
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.logError("Invalid Category Parameters", w, r, http.StatusBadRequest, err)
		return
	}

	params.Name, err = h.sanitize(params.Name, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param name", w, r, http.StatusBadRequest, err)
		return
	}

	params.Color, err = h.sanitize(params.Color, CHARLIMIT)
	if err != nil {
		h.logError("Invalid param name", w, r, http.StatusBadRequest, err)
		return
	}

	if !IsValidHex(params.Color) {
		h.logError("Invalid HEX color", w, r, http.StatusBadRequest, nil)
		return
	}

	category, err := h.Queries.CreateCategory(r.Context(), params)
	if err != nil {
		h.logError("Couldn't create category", w, r, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *App) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	category, err := h.sanitize(chi.URLParam(r, "category"), CHARLIMIT)
	if err != nil {
		h.logError("Invalid category name", w, r, http.StatusBadRequest, err)
		return
	}

	h.Queries.DeleteCategory(r.Context(), category)
}

func (h *App) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: implement it
}

func IsValidHex(color string) bool {
	digits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	chars := []string{"a", "b", "c", "d", "e", "f", "A", "B", "C", "D", "E", "F"}

	if len(color) != 7 {
		return false
	}

	if color[0] != '#' {
		return false
	}

	for _, c := range color[1:] {
		ch := string(c)
		if !slices.Contains(digits, ch) && !slices.Contains(chars, ch) {
			return false
		}
	}

	return true
}
