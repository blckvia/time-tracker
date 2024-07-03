package repository

import (
	"context"
	"fmt"

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
	tx, err := r.db.Begin(r.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(r.ctx)

	var exists bool

	err = tx.QueryRow(r.ctx, fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", tasksTable), taskID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return ErrNotFound
	}

	_, err = tx.Exec(r.ctx, fmt.Sprintf("UPDATE %s SET started_at = NOW() WHERE id = $1", tasksTable), taskID)
	if err != nil {
		return err
	}

	err = tx.Commit(r.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *TasksPostgres) StopTask(taskID int) error {
	return nil
}
