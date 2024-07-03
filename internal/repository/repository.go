package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"time-tracker/internal/entities"
)

type Users interface {
	Create(input *entities.Users) (int, error)
	Update(userID int, input *entities.Users) error
	Delete(userID int) error
	GetAll(filters map[string]any, limit, offset int) (entities.GetAllUsers, error)
	GetByID(userID int) (entities.Users, error)
}

type Tasks interface {
	Create(input *entities.Task) (int, error)
	StartTask(taskID int) error
	StopTask(taskID int) error
}

type Repository struct {
	Users
	Tasks
}

func NewRepository(ctx context.Context, db *pgxpool.Pool, logger *zap.Logger) *Repository {
	return &Repository{
		Users: NewUsersPostgres(ctx, db, logger),
		Tasks: NewTasksPostgres(ctx, db, logger),
	}
}