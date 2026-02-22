package handler

import (
	"GOLANG/internal/usecase"
	"GOLANG/pkg/modules"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	switch r.Method {
	case http.MethodGet:
		if len(pathParts) > 1 {
			id, _ := strconv.Atoi(pathParts[1])
			user, err := h.uc.GetByID(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			json.NewEncoder(w).Encode(user)
		} else {
			users, _ := h.uc.GetAll()
			json.NewEncoder(w).Encode(users)
		}

	case http.MethodPost:
		var u modules.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := h.uc.Create(u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": id})

	case http.MethodPut:
		var u modules.User
		json.NewDecoder(r.Body).Decode(&u)
		if err := h.uc.Update(u); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		if len(pathParts) > 1 {
			id, _ := strconv.Atoi(pathParts[1])
			affected, err := h.uc.Delete(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]int64{"rows_affected": affected})
		}
	}
}

func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "UP"}`))
}
