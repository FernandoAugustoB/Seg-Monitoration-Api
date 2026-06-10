package routes

import (
	"Seg-Monitoration-Api/internal/auth"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func EnsurePermission(requiredPermission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Extrai o header Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token não fornecido", http.StatusUnauthorized)
				return
			}

			// O formato esperado é "Bearer <TOKEN>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]

			// 2. Valida o JWT e mapeia para as Claims que criamos
			claims := &auth.Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte("sua_chave_secreta"), nil // Mudar para os.Getenv("JWT_SECRET") depois
			})

			if err != nil || !token.Valid {
				http.Error(w, "Token inválido ou expirado", http.StatusUnauthorized)
				return
			}

			// 3. Checa a permissão (Estilo Discord)
			hasPermission := false
			for _, perm := range claims.Permissions {
				if perm == requiredPermission || perm == "admin.gerenciar" { // Admin pula a checagem
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				http.Error(w, "Acesso negado: permissão insuficiente", http.StatusForbidden)
				return
			}

			// 4. Injeta o ID do usuário no Contexto (Útil para saber quem criou um relatório, ex: r.Context().Value(UserIDKey))
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
