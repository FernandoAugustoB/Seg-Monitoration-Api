package models

import "time"

type User struct {
    ID           string    `json:"id"`
    Username     string    `json:"username"`
    PasswordHash string    `json:"-"` // O "-" garante que a senha nunca saia no JSON da API
    CreatedAt    time.Time `json:"created_at"`
}