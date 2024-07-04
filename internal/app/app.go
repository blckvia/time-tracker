package app

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"time-tracker/internal/handler"
	"time-tracker/internal/repository"
	"time-tracker/internal/service"
)

// @title Time Tracker API
// @version 1.0
// @description API Server for Go Service

// @host localhost:8080
// @BasePath /

type App struct {
	Server *http.Server
	Logger *zap.Logger
	db     *pgxpool.Pool
}

func NewApp(ctx context.Context, logger *zap.Logger) *App {
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file", zap.Error(err))
	}

	db, err := repository.NewPostgresDB(ctx, &repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})

	if err != nil {
		logger.Fatal("failed to initialize db", zap.Error(err))
	}

	repos := repository.NewRepository(ctx, db, logger)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        handlers.InitRoutes(),
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return &App{
		Server: srv,
		Logger: logger,
		db:     db,
	}
}

func (a *App) Run() error {
	return a.Server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.db != nil {
		a.db.Close()
	}

	return a.Server.Shutdown(ctx)
}
