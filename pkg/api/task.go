package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Izhdenev/Final-Project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	writeJSON(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if _, err := strconv.ParseInt(task.ID, 10, 64); err != nil {
		writeJSON(w, map[string]string{"error": "некорректный идентификатор задачи"})
		return
	}

	title := strings.TrimSpace(task.Title)
	if title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}
	task.Title = title

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	writeJSON(w, struct{}{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, struct{}{})
		return
	}

	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, struct{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, struct{}{})
}
