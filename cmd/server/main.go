package main

import (
	"Seg-Monitoration-Api/internal/routes"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	// 1. Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	// 2. Conecta ao Postgres (usando Pool de conexões para performance)
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco: %v", err)
	}
	defer dbPool.Close()

	// 3. Middlewares Básicos
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
	router := routes.SetupRoutes(dbPool, r)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/docs.json"), //The url pointing to API definition
	))

	// 4. Definição de Rotas

	log.Printf("Servidor rodando na porta %s", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}
