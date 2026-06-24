package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ema/fixture/backend/internal/auth"
)

type authHandler struct{}

// NewAuthHandler crea un handler para login.
func NewAuthHandler() Handler {
	return &authHandler{}
}

func (h *authHandler) Register(r Router) {
	r.Post("/api/auth/login", h.Login)
}

type loginRequest struct {
	PIN string `json:"pin"`
}

type loginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func getAdminPIN() string {
	if pin := os.Getenv("ADMIN_PIN"); pin != "" {
		return pin
	}
	return "1234"
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if req.PIN == "" {
		writeError(w, http.StatusBadRequest, "PIN requerido")
		return
	}

	if req.PIN != getAdminPIN() {
		writeError(w, http.StatusUnauthorized, "PIN incorrecto")
		return
	}

	token := auth.GenerarToken()
	writeJSON(w, http.StatusOK, loginResponse{
		Token: token,
		Role:  "admin",
	})
}
