package handlers

import (
	"Seg-Monitoration-Api/internal/repository"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary Realiza o login do usuário
// @Description Recebe credenciais e retorna um JWT com as permissões do usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body handlers.LoginRequest true "Credenciais de Login"
// @Success 200 {object} map[string]string "token: <jwt_token>"
// @Failure 401 {string} string "Credenciais inválidas"
// @Router /auth/login [post]
func LoginHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

		// 1. Decode do JSON (Otimizado para não ler mais que o necessário)
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// 2. Validação básica (Segurança de Input)
		if req.Username == "" || req.Password == "" {
			http.Error(w, "Usuário e senha são obrigatórios", http.StatusBadRequest)
			return
		}

		// 3. Busca no Banco de Dados (Repository)
		// Aqui chamaremos a função que busca o Hash da senha no Postgres
		_, err = repository.GetUserHash(r.Context(), db, req.Username)
		if err != nil {
			// Dica de Segurança: Não diga se o usuário existe ou não, use uma mensagem genérica
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}

		// 4. Verificação de Senha (Bcrypt) e Geração de Token (JWT)
		// [Implementaremos na sequência...]
	}
}
