// cmd/server/main.go
package main

// Comando para automatizar o Swagger. Rode: go generate ./...
//go:generate swag init -g main.go -o ../../docs

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Seg-Monitoration-Api/internal/routes"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Carrega as variáveis do arquivo .env
	// No ambiente de produção do GitHub Codespaces ou Docker, as ENVs podem vir diretamente do sistema
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. Usando variáveis de ambiente do sistema.")
	}

	// 2. Valida se as variáveis essenciais existem
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Erro: A variável de ambiente DB_URL não foi definida.")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback padrão
	}

	// 3. Configura e conecta ao Pool do Postgres (pgxpool)
	// O Background context é usado aqui porque essa conexão dura o ciclo de vida inteiro da aplicação
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Não foi possível conectar ao pool do banco de dados: %v\n", err)
	}
	defer func() {
		log.Println("Fechando pool de conexões com o banco de dados...")
		dbPool.Close()
	}()

	// Testa a conexão imediatamente para não subir a API com o banco fora do ar
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Erro ao pingar o banco de dados: %v\n", err)
	}
	log.Println("Conexão com o PostgreSQL estabelecida com sucesso!")

	// 4. Inicializa o Hub de Rotas Componentizadas
	router := routes.SetupRoutes(dbPool)

	// 5. Configura o Servidor HTTP
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second, // Proteção contra ataques de conexões lentas (Slowloris)
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 6. Canal para escutar sinais de interrupção do Sistema Operacional (Garuda/Linux)
	shutdownChan := make(chan os.Signal, 1)
	// Escuta SIGINT (Ctrl+C) e SIGTERM (parada do container Docker)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Subimos o servidor em uma Goroutine separada para não bloquear a thread principal
	go func() {
		log.Printf("Servidor API operacional na porta %s...\n", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Falha crítica no servidor HTTP: %v\n", err)
		}
	}()

	// O programa fica travado aqui esperando um sinal de parada
	<-shutdownChan
	log.Println("Sinal de encerramento recebido. Iniciando Graceful Shutdown...")

	// 7. Contexto com Timeout de 5 segundos para o encerramento limpo
	// Dá tempo para a API terminar de processar requisições ativas de relatórios ou logins
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Erro durante o encerramento forçado do servidor: %v\n", err)
	}

	log.Println("Servidor finalizado com sucesso. Até logo!")
}
