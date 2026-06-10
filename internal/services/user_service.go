package services

import (
	"Seg-Monitoration-Api/internal/auth"
	"Seg-Monitoration-Api/internal/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	DB *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{DB: db}
}

// Signup: Cria um novo usuário com senha protegida
func (s *UserService) Signup(ctx context.Context, username, password string) error {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	// Iniciando a Transação
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}
	// Garante que o rollback ocorra se houver erro; se houver commit, o rollback não faz nada
	defer tx.Rollback(ctx)

	// 1. Criar o Usuário
	var userID string
	err = tx.QueryRow(ctx, "INSERT INTO usuarios (username, password_hash) VALUES ($1, $2) RETURNING id", username, hashedPassword).Scan(&userID)
	if err != nil {
		return err
	}

	// 2. Associar ao grupo "Membro" (Supondo que você já tenha esse grupo criado no banco)
	// Buscamos o ID do grupo pelo nome para facilitar
	_, err = tx.Exec(ctx, `
        INSERT INTO usuario_grupos (usuario_id, grupo_id) 
        SELECT $1, id FROM grupos WHERE nome = 'Membro' LIMIT 1
    `, userID)
	if err != nil {
		return err
	}

	// Finaliza a transação com sucesso
	return tx.Commit(ctx)
}

// Login: Valida credenciais e retorna o Token com permissões (Estilo Discord)
func (s *UserService) Login(ctx context.Context, username, password string) (string, error) {
	// 1. Busca o usuário no banco pelo username
	user, err := repository.GetUserByUsername(ctx, s.DB, username)
	if err != nil {
		return "", errors.New("credenciais inválidas") // Erro genérico por segurança
	}

	// 2. Compara o hash da senha usando o Bcrypt
	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("credenciais inválidas")
	}

	// 3. Busca as permissões do usuário (Estilo Discord)
	permissions, err := repository.GetUserPermissions(ctx, s.DB, user.ID)
	if err != nil {
		return "", errors.New("erro ao carregar permissões do usuário")
	}

	// 4. Gera o JWT injetando as permissões nas Claims
	token, err := auth.GenerateJWT(user.ID, permissions)
	if err != nil {
		return "", errors.New("erro ao gerar token de acesso")
	}

	return token, nil
}

// UpdateAccount: Altera informações (ex: trocar senha ou username)
func (s *UserService) UpdateAccount(ctx context.Context, userID string, newUsername string, newPassword string) error {
	var hashed string
	if newPassword != "" {
		var err error
		hashed, err = auth.HashPassword(newPassword)
		if err != nil {
			return err
		}
	}
	return repository.UpdateUser(ctx, s.DB, userID, newUsername, hashed)
}

// DeleteAccount: Remove a conta
func (s *UserService) DeleteAccount(ctx context.Context, userID string) error {
	return repository.DeleteUser(ctx, s.DB, userID)
}
