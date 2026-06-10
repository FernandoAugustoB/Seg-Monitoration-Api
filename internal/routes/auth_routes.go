package routes

import (
	"Seg-Monitoration-Api/internal/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(s *services.UserService) *UserController {
	return &UserController{Service: s}
}

// SignupRequest define o corpo esperado para cadastro
type SignupRequest struct {
	Username string `json:"username" example:"operador_01"`
	Password string `json:"password" example:"senha_forte_123"`
}

func (c *UserController) AuthRoutes(r chi.Router) {
	//
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", c.SignupHandle)
		r.Post("/login", c.LoginHandle)
		// r.Post("/logout", c.LogoutHandle) -> Exemplo futuro
	})

	// Rotas de perfil que exigem autenticação
	r.Route("/user", func(r chi.Router) {
		// r.Use(middleware.Auth) -> Middleware aqui
		// r.Get("/me", c.GetProfileHandle)
		// r.Put("/update", c.UpdateHandle)
	})
}

// SignupHandle cria um novo usuário
// @Summary Cadastro de novo usuário
// @Description Cria uma conta e atribui o grupo padrão 'Membro' via transação
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body SignupRequest true "Dados de cadastro"
// @Success 201 {string} string "Usuário criado com sucesso"
// @Failure 400 {string} string "Dados inválidos"
// @Router /auth/signup [post]
func (c *UserController) SignupHandle(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	err := c.Service.Signup(r.Context(), req.Username, req.Password)
	if err != nil {
		// Aqui poderíamos tratar erros específicos (ex: usuário já existe)
		http.Error(w, "Erro ao criar usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Usuário criado com sucesso"))
}

// LoginHandle autentica e retorna o JWT
// @Summary Login do usuário
// @Description Valida credenciais e retorna JWT com permissões estilo Discord
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body SignupRequest true "Credenciais"
// @Success 200 {object} map[string]string "token"
// @Failure 401 {string} string "Credenciais inválidas"
// @Router /auth/login [post]
func (c *UserController) LoginHandle(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	token, err := c.Service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
