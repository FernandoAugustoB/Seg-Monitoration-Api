package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(db *pgxpool.Pool, r *chi.Mux) *chi.Mux {
	RegisterAuthRoutes(r)

	return r
}
