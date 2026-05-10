// internal/repository/user_repository.go
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserHash(ctx context.Context, db *pgxpool.Pool, username string) (string, error) {
	var hash string
	query := "SELECT password_hash FROM usuarios WHERE username = $1 LIMIT 1"

	// O pgx usa prepared statements automaticamente aqui ($1)
	err := db.QueryRow(ctx, query, username).Scan(&hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}
