package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type authResponse struct {
	Message string `json:"message"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterAuthRoutes(mux *chi.Mux) {
	mux.HandleFunc("/auth/register", registerHandler)
	mux.HandleFunc("/auth/logout", logoutHandler)
}

// @Summary Registra novo usuário
// @Description Cria uma nova conta de usuário a partir das credenciais fornecidas
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Dados de registro"
// @Success 201 {object} authResponse
// @Failure 400 {object} authResponse
// @Failure 405 {object} authResponse
// @Router /auth/register [post]
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, authResponse{Message: "method not allowed"})
		return
	}

	// TODO: create user account
	writeJSON(w, http.StatusCreated, authResponse{Message: "registration successful"})
}

// @Summary Encerra a sessão do usuário
// @Description Invalida o token ou revoga a sessão do usuário autenticado
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} authResponse
// @Failure 401 {object} authResponse
// @Failure 405 {object} authResponse
// @Router /auth/logout [post]
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, authResponse{Message: "method not allowed"})
		return
	}

	// TODO: revoke session or token
	writeJSON(w, http.StatusOK, authResponse{Message: "logout successful"})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
