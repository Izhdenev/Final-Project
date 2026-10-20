package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Izhdenev/Final-Project/pkg/db"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if strings.TrimSpace(task.Date) == "" {
		task.Date = now.Format(layout)
	}

	t, err := time.Parse(layout, task.Date)
	if err != nil {
		return fmt.Errorf("дата представлена в формате, отличном от 20060102")
	}

	repeat := strings.TrimSpace(task.Repeat)
	isPast := t.Before(nowDate)

	if repeat != "" {
		next, err := NextDate(now, task.Date, repeat)
		if err != nil {
			return fmt.Errorf("правило повторения указано в неправильном формате")
		}
		if isPast {
			task.Date = next
		}
	} else if isPast {
		task.Date = now.Format(layout)
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	title := strings.TrimSpace(task.Title)
	if title == "" {
		writeError(w, http.StatusBadRequest, "не указан заголовок задачи")
		return
	}
	task.Title = title

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка сохранения задачи в базе данных")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": fmt.Sprintf("%d", id)})
}
