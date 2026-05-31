package api

import (
	"encoding/json"
	"net/http"
	"os"
	"sprint13-14/pkg/auth"
)

type SignInRequest struct {
	Password string `json:"password"`
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	var req SignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if req.Password != expectedPassword {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Неверный пароль"})
		return
	}

	token, err := auth.GenerateToken()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
