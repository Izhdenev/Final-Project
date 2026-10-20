package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Izhdenev/Final-Project/pkg/db"
)

func writeTaskError(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "Задача не найдена")
		return
	}
	writeError(w, http.StatusInternalServerError, err.Error())
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	if _, err := strconv.ParseInt(task.ID, 10, 64); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный идентификатор задачи")
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

	if err := db.UpdateTask(&task); err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeTaskError(w, err)
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			writeTaskError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, struct{}{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}
