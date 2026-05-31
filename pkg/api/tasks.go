package api

import (
	"encoding/json"
	"net/http"
	"sprint13-14/pkg/db"
)

type TasksResp struct {
	Tasks []db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	searchParam := r.FormValue("search")
	tasks, err := db.Tasks(50, searchParam)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, TasksResp{Tasks: tasks})
}

func writeJson(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMsg})
}
