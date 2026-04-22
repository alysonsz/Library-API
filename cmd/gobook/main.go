package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"project-go/internal/cli"
	"project-go/internal/config"
	"project-go/internal/database"
	"project-go/internal/services"
	"project-go/internal/web"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "project-go/docs"
)

// @title           Library API
// @version         1.0
// @description     API RESTful para gerenciamento de livraria com suporte a SQLite e PostgreSQL.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Support
// @contact.url    https://github.com/alysonsz/library-api

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @schemes   http https
func main() {
	cfg := config.Load()

	logLevel := slog.LevelInfo
	if cfg.LogLevel == "debug" {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	var rawDB *sql.DB
	var err error

	if cfg.DBDriver == "postgres" {
		if cfg.DBURL == "" {
			slog.Error("DB_URL is required for postgres driver")
			os.Exit(1)
		}
		rawDB, err = sql.Open("pgx", cfg.DBURL)
	} else {
		rawDB, err = sql.Open("sqlite3", cfg.DBPath)
	}

	if err != nil {
		slog.Error("failed to open database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer rawDB.Close()

	db := database.NewDB(rawDB, cfg.DBDriver)

	if err := database.Migrate(db); err != nil {
		slog.Error("failed to migrate database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	bookService := services.NewBookService(db)
	authorService := services.NewAuthorService(db)
	loanService := services.NewLoanService(db)

	bookHandlers := web.NewBookHandlers(bookService)
	authorHandlers := web.NewAuthorHandlers(authorService)
	loanHandlers := web.NewLoanHandlers(loanService)
	healthHandler := web.NewHealthHandler(db)

	if len(os.Args) > 1 && (os.Args[1] == "simulate" || os.Args[1] == "search" || os.Args[1] == "authors" || os.Args[1] == "loans") {
		bookCLI := cli.NewBookCLI(bookService, authorService, loanService)
		bookCLI.Run()
		return
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /books", bookHandlers.GetBooks)
	router.HandleFunc("POST /books", bookHandlers.CreateBook)
	router.HandleFunc("GET /books/{id}", bookHandlers.GetBookByID)
	router.HandleFunc("PUT /books/{id}", bookHandlers.UpdateBook)
	router.HandleFunc("DELETE /books/{id}", bookHandlers.DeleteBook)

	router.HandleFunc("GET /authors", authorHandlers.GetAuthors)
	router.HandleFunc("POST /authors", authorHandlers.CreateAuthor)
	router.HandleFunc("GET /authors/{id}", authorHandlers.GetAuthorByID)
	router.HandleFunc("PUT /authors/{id}", authorHandlers.UpdateAuthor)
	router.HandleFunc("DELETE /authors/{id}", authorHandlers.DeleteAuthor)
	router.HandleFunc("GET /authors/{id}/books", authorHandlers.GetBooksByAuthor)

	router.HandleFunc("GET /loans", loanHandlers.GetLoans)
	router.HandleFunc("POST /loans", loanHandlers.CreateLoan)
	router.HandleFunc("GET /loans/{id}", loanHandlers.GetLoanByID)
	router.HandleFunc("POST /loans/{id}/return", loanHandlers.ReturnLoan)
	router.HandleFunc("GET /loans/overdue", loanHandlers.GetOverdueLoans)
	router.HandleFunc("GET /health", healthHandler.Check)
	router.HandleFunc("GET /swagger/*", httpSwagger.WrapHandler)

	rateLimiter := web.NewRateLimiter(100, time.Minute)
	handler := web.LoggingMiddleware(web.CORSMiddleware(rateLimiter.Middleware(router)))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	go func() {
		slog.Info("server starting", slog.String("addr", server.Addr))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("server exited")
}
