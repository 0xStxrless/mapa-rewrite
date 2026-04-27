package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/0xstxrless/punkt/backend/internal/auth"
	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

type contextKey string

const userContextKey contextKey = "user"

type Claims struct {
	UserID int32  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token              string `json:"token"`
	MustChangePassword bool   `json:"must_change_password"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (h *App) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logError("Invalid request body", w, http.StatusBadRequest, err)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Run a dummy scrypt verify so timing is indistinguishable
			// from a real failed attempt — prevents user enumeration.
			auth.VerifyPassword(req.Password, "scrypt:deadbeef:deadbeef")
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		h.logError("DB error during login", w, http.StatusInternalServerError, err)
		return
	}

	if !auth.VerifyPassword(req.Password, user.PasswordHash) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Migrate legacy SHA-256 or scrypt hash to the current scrypt format.
	// SetPasswordHash also flips must_change_password = false, so only
	// migrate if the user isn't still on a seeded password.
	isLegacy := !strings.HasPrefix(user.PasswordHash, "scrypt:")
	if isLegacy {
		if newHash, err := auth.HashPassword(req.Password); err == nil {
			_ = h.queries.SetPasswordHash(r.Context(), db.SetPasswordHashParams{
				ID:           user.ID,
				PasswordHash: newHash,
			})
		}
	}

	_ = h.queries.UpdateLastLogin(r.Context(), user.ID)

	secret, err := auth.GetSessionSecret()
	if err != nil {
		h.logError("Session secret unavailable", w, http.StatusInternalServerError, err)
		return
	}

	token, err := auth.SignSession(auth.SessionPayload{
		UserID: user.ID,
		Email:  user.Email,
	}, secret)
	if err != nil {
		h.logError("Failed to sign session", w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{
		Token:              token,
		MustChangePassword: user.MustChangePassword.Bool && user.MustChangePassword.Valid,
	})
}

func (h *App) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logError("Invalid request body", w, http.StatusBadRequest, err)
		return
	}

	if len(req.NewPassword) < 8 {
		http.Error(w, "new password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	user, err := h.queries.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		h.logError("User not found", w, http.StatusNotFound, err)
		return
	}

	if !auth.VerifyPassword(req.CurrentPassword, user.PasswordHash) {
		http.Error(w, "invalid current password", http.StatusUnauthorized)
		return
	}

	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		h.logError("Failed to hash password", w, http.StatusInternalServerError, err)
		return
	}

	if err := h.queries.SetPasswordHash(r.Context(), db.SetPasswordHashParams{
		ID:           claims.UserID,
		PasswordHash: newHash,
	}); err != nil {
		h.logError("Failed to update password", w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AuthMiddleware validates the signed session token and injects claims into context.
func (h *App) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "missing or malformed token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		secret, err := auth.GetSessionSecret()
		if err != nil {
			http.Error(w, "server misconfiguration", http.StatusInternalServerError)
			return
		}

		payload, err := auth.VerifyAndParseSession(tokenStr, secret)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{
			UserID: payload.UserID,
			Email:  payload.Email,
		}
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClaimsFromContext(ctx context.Context) *Claims {
	c, _ := ctx.Value(userContextKey).(*Claims)
	return c
}
