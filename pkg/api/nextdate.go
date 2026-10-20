package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Izhdenev/Final-Project/pkg/db"
)

const layout = db.DateLayout

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	date, err := time.Parse(layout, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid dstart date: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(repeat), " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat rule format")
	}

	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid 'd' rule format: expected 'd <number>'")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("interval days must be an integer: %w", err)
		}
		if days <= 0 || days > 400 {
			return "", fmt.Errorf("interval days must be between 1 and 400")
		}

		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				return date.Format(layout), nil
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid 'y' rule format: expected just 'y'")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format(layout), nil
			}
		}

	default:
		return "", fmt.Errorf("unsupported repeat rule format: %s", rule)
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	var now time.Time
	if nowStr != "" {
		var err error
		now, err = time.Parse(layout, nowStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "parameter 'now' has invalid date format")
			return
		}
	} else {
		now = time.Now()
	}

	if dateStr == "" {
		writeError(w, http.StatusBadRequest, "required parameter 'date' is missing")
		return
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, next)
}
