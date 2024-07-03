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

func (r *TasksPostgres) Create(input *entities.Task, userID int) (int, error) {
	var id int

	var exists bool

	err := r.db.QueryRow(r.ctx, fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", tasksTable), userID).Scan(&exists)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, ErrNotFound
	}

	query := fmt.Sprintf(`INSERT INTO %s (task, user_id, description) VALUES ($1, $2, $3) RETURNING id`, tasksTable)

	conn, err := r.db.Acquire(r.ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	pgxConn := conn.Conn()

	_, err = pgxConn.Prepare(r.ctx, "createTask", query)
	if err != nil {
		return 0, err
	}

	row := pgxConn.QueryRow(r.ctx, "createTask", input.Task, userID, input.Description)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
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

	var timer bool

	err = r.db.QueryRow(r.ctx, fmt.Sprintf("SELECT tp.timer FROM %s tp WHERE tp.id = $1", tasksTable), taskID).Scan(&timer)
	if err != nil {
		return err
	}
	fmt.Println(timer)

	if timer == true {
		return fmt.Errorf("timer is already active for task %d", taskID)
	}

	_, err = tx.Exec(r.ctx, fmt.Sprintf(`
	UPDATE %s 
	SET start_time = CASE WHEN timer = false THEN NOW() ELSE start_time END,
	    timer = CASE WHEN timer = false THEN true ELSE timer END
	WHERE id = $1`, tasksTable), taskID)
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
