// internal/repository/user_repository.go
package repository

import (
	"context"

	"Seg-Monitoration-Api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserHash(ctx context.Context, db *pgxpool.Pool, username string) (string, error) {
	var hash string
	query := "SELECT password_hash FROM usuarios WHERE username = $1 LIMIT 1"

	err := db.QueryRow(ctx, query, username).Scan(&hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// GetUserPermissions busca todas as permissões únicas dos grupos do usuário
func GetUserPermissions(ctx context.Context, db *pgxpool.Pool, userID string) ([]string, error) {
	query := `
		SELECT DISTINCT p.nome 
		FROM permissoes p
		JOIN grupo_permissoes gp ON gp.permissao_id = p.id
		JOIN usuario_grupos ug ON ug.grupo_id = gp.grupo_id
		WHERE ug.usuario_id = $1
	`

	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permName string
		if err := rows.Scan(&permName); err != nil {
			return nil, err
		}
		permissions = append(permissions, permName)
	}

	// Se o slice vier vazio (nil), transformamos em um slice vazio instanciado [] para o JSON não quebrar
	if permissions == nil {
		permissions = []string{}
	}

	return permissions, nil
}

// CreateUser insere o novo usuário no Postgres
func CreateUser(ctx context.Context, db *pgxpool.Pool, username, hash string) error {
	_, err := db.Exec(ctx, "INSERT INTO usuarios (username, password_hash) VALUES ($1, $2)", username, hash)
	return err
}

// GetUserByUsername busca os dados básicos para o login
func GetUserByUsername(ctx context.Context, db *pgxpool.Pool, username string) (models.User, error) {
	var u models.User
	err := db.QueryRow(ctx, "SELECT id, username, password_hash FROM usuarios WHERE username = $1", username).Scan(&u.ID, &u.Username, &u.PasswordHash)
	return u, err
}

// DeleteUser remove o usuário (as relações em outras tabelas caem pelo ON DELETE CASCADE do SQL que criamos)
func DeleteUser(ctx context.Context, db *pgxpool.Pool, id string) error {
	_, err := db.Exec(ctx, "DELETE FROM usuarios WHERE id = $1", id)
	return err
}

func UpdateUser(ctx context.Context, db *pgxpool.Pool, userID string, username string, hash string) error {
	// Se ambos forem vazios, não há o que atualizar
	if username == "" && hash == "" {
		return nil
	}

	// Exemplo de query que lida com campos opcionais
	// Em um projeto real mais complexo, você poderia usar um Query Builder
	if username != "" && hash != "" {
		_, err := db.Exec(ctx, "UPDATE usuarios SET username = $1, password_hash = $2 WHERE id = $3", username, hash, userID)
		return err
	} else if username != "" {
		_, err := db.Exec(ctx, "UPDATE usuarios SET username = $1 WHERE id = $2", username, userID)
		return err
	} else {
		_, err := db.Exec(ctx, "UPDATE usuarios SET password_hash = $1 WHERE id = $2", hash, userID)
		return err
	}
}
