package routes

import (
	"Seg-Monitoration-Api/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func SetupRoutes(db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Logger)    // Log de requisições no terminal
	r.Use(middleware.Recoverer) // Impede que a API caia se houver um panic
	r.Use(middleware.RealIP)    // Captura o IP real (importante para auditoria)

	// Serviços
	userService := services.NewUserService(db)

	// Controladores
	userController := NewUserController(userService)

	// Rotas
	userController.AuthRoutes(r)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/docs.json"), //The url pointing to API definition
	))

	return r
}
