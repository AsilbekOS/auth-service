package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	AuthService *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{
		AuthService: authService,
	}
}

func (h *Handler) GenerateTokensHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	accessToken, refreshToken, err := h.AuthService.GenerateTokens(userID)
	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newAccessToken, newRefreshToken, err := h.AuthService.RefreshToken(request.RefreshToken)
	if err != nil {
		http.Error(w, "Error refreshing token", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
