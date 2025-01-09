package auth_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwt "github.com/alroymuhammad/go-go-manager/internal/appjwt"
	"github.com/alroymuhammad/go-go-manager/internal/usecase/auth_usecase"
)

type AuthHandler struct {
	Service *auth_usecase.AuthService
}

func NewAuthHandler(service *auth_usecase.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Action   string `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	switch req.Action {
	case "create":
		h.handleCreate(w, req.Email, req.Password)
	case "login":
		h.handleLogin(w, req.Email, req.Password)
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
	}
}

func (h *AuthHandler) handleCreate(w http.ResponseWriter, email, password string) {
	user, err := h.Service.Register(email, password)
	if err != nil {
		switch err {
		case auth_usecase.ErrUserExists:
			http.Error(w, "Email already exists", http.StatusConflict)
		case auth_usecase.ErrInvalidInput:
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
		default:
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	token, err := jwt.GenerateJWT(fmt.Sprintf("%d", user.ID))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusCreated, map[string]string{
		"email": user.Email,
		"token": token,
	})
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, email, password string) {
	token, err := h.Service.Login(email, password)
	if err != nil {
		switch err {
		case auth_usecase.ErrInvalidCredentials:
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		case auth_usecase.ErrInvalidInput:
			http.Error(w, "Invalid input parameters", http.StatusBadRequest)
		default:
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, map[string]string{
		"email": email,
		"token": token,
	})
}

func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
