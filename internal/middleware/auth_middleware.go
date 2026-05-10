// internal/middleware/auth.go
package middleware

import (
	"Seg-Monitoration-Api/internal/auth"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func EnsurePermission(requiredPermission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Pega o token do Header (Bearer <token>)
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token não fornecido", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			// 2. Faz o parse e valida o JWT
			claims := &auth.Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte("sua_chave_secreta"), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			// 3. Verifica se a permissão necessária está no array de permissões do usuário
			hasPerm := false
			for _, p := range claims.Permissions {
				if p == requiredPermission {
					hasPerm = true
					break
				}
			}

			if !hasPerm {
				http.Error(w, "Acesso negado: permissão insuficiente", http.StatusForbidden)
				return
			}

			// 4. Se chegou aqui, injeta o UserID no contexto para uso futuro
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
