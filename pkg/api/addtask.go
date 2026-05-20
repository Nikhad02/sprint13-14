package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"sprint13-14/pkg/db"
	"strings"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func afterNow(now time.Time, t time.Time) bool {
	nowZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tZero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return nowZero.After(tZero)
}

func checkAndFixData(task *db.Task) error {
	now := time.Now()

	task.Repeat = strings.TrimSpace(task.Repeat)

	if task.Date == "" {
		task.Date = now.Format(DateLayout)
	}
	t, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		return errors.New("некорректный формат даты")
	}
	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(DateLayout)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return errors.New("некорректное правило повторения")
			}
			task.Date = next
		}
	} else if task.Repeat != "" {
		_, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return errors.New("некорректное правило повторения")
		}
	}
	return nil
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок"})
		return
	}
	err = checkAndFixData(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка сохранения в базу данных"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"id": id})
}
