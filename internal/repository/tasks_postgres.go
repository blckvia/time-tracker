package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"time-tracker/internal/entities"
)

type TasksPostgres struct {
	ctx context.Context
	db  *pgxpool.Pool
	zap.Logger
}

func NewTasksPostgres(ctx context.Context, db *pgxpool.Pool, logger *zap.Logger) *TasksPostgres {
	return &TasksPostgres{
		ctx:    ctx,
		db:     db,
		Logger: *logger,
	}
}

func (r *TasksPostgres) Create(input *entities.Task) (int, error) {
	return 0, nil
}

func (r *TasksPostgres) StartTask(taskID int) error {
	return nil
}

func (r *TasksPostgres) StopTask(taskID int) error {
	return nil
}
