package router

import (
	"context"
	"errors"
	"html"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
)

type contextKey string

const loggerKey contextKey = "logger"

func (h *App) paramIDtoInt(r *http.Request) (int, error) {
	pinIDStr, err := h.sanitize(chi.URLParam(r, "id"), CHARLIMIT)
	if err != nil {
		return 0, err
	}
	pinID, err := strconv.Atoi(pinIDStr)
	if err != nil {
		return 0, err
	}

	if pinID < 0 {
		return 0, err
	}
	return pinID, nil
}

func (h *App) logError(msg string, w http.ResponseWriter, r *http.Request, code int, err error) {
	http.Error(w, msg, code)
	loggerFromCtx(r.Context()).Error(msg, "error", err, "http_status", code)
}

// this is general input sanitization against xss, injection, html escape etc
// IT DOESN'T SANITIZE AGAINST SQLI, PATH TRAVERSALS AND COMMAND INJECTIONS
func (h *App) sanitize(s string, maxLen int) (string, error) {
	if len(s) > maxLen {
		return "", errors.New("input exceeds maximum length")
	}

	if s == "" {
		return "", errors.New("empty string")
	}

	if !utf8.ValidString(s) {
		return "", errors.New("invalid UTF-8 encoding")
	}

	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
	s = html.EscapeString(s)
	s = strings.TrimSpace(s)

	return s, nil
}

func NewLogger() *slog.Logger {
	// JSON in prod (parseable by Datadog, Loki, CloudWatch etc)
	// Text in dev (human readable)
	var handler slog.Handler

	if os.Getenv("ENV") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo, // don't flood prod with debug
			AddSource: true,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	}

	return slog.New(handler)
}
func (h *App) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		log := h.Logger.With(
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"ip", r.RemoteAddr,
		)

		ctx := context.WithValue(r.Context(), loggerKey, log)

		start := time.Now()
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Info("request completed", "duration_ms", time.Since(start).Milliseconds())
		// add some spacing in logs for better readability
		println()
		println()
	})
}

func loggerFromCtx(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
