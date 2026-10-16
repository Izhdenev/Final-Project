package api

import (
	"net/http"

	"github.com/Izhdenev/Final-Project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
