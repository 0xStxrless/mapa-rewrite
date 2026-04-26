package router

import (
	"errors"
	"html"
	"net/http"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
)

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

func (h *App) logError(msg string,
	w http.ResponseWriter,
	code int,
	err error,
) {
	http.Error(w, msg, code)
	h.logger.Printf("%s: %v\n", msg, err)
}

// this is general input sanitization agains xss, injection, html escape etc
// IT DOESN'T SANITIZE AGAINST SQLI, PATH TRAVERSALS AND COMMAND INJECTIONS
func (h *App) sanitize(s string, maxLen int) (string, error) {
	if len(s) > maxLen {
		return "", errors.New("input exceeds maximum lenght")
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
