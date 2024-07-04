package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"time-tracker/internal/entities"
)

var ErrNotFound = errors.New("record not found")

type UsersPostgres struct {
	ctx context.Context
	db  *pgxpool.Pool
	zap.Logger
}

func NewUsersPostgres(ctx context.Context, db *pgxpool.Pool, logger *zap.Logger) *UsersPostgres {
	return &UsersPostgres{
		ctx:    ctx,
		db:     db,
		Logger: *logger,
	}
}

func (r *UsersPostgres) Create(input *entities.Users) (int, error) {
	var id int

	query := fmt.Sprintf(`INSERT INTO %s (name, surname, patronymic, passport_series, passport_number, address) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, usersTable)

	conn, err := r.db.Acquire(r.ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	pgxConn := conn.Conn()

	_, err = pgxConn.Prepare(r.ctx, "createUser", query)
	if err != nil {
		return 0, err
	}

	row := pgxConn.QueryRow(r.ctx, "createUser", input.Name, input.Surname, input.Patronymic, input.PassportSeries, input.PassportNumber, input.Address)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UsersPostgres) Update(userID int, input *entities.Users) error {
	tx, err := r.db.Begin(r.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(r.ctx)

	var exists bool

	err = tx.QueryRow(r.ctx, fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", usersTable), userID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}

	_, err = tx.Exec(r.ctx, fmt.Sprintf("SELECT up.id, up.passport_series, up.passport_number, up.name, up.surname, up.patronymic, up.address FROM %s up WHERE id = $1 FOR UPDATE", usersTable), userID)
	if err != nil {
		return err
	}
	setValues := make([]string, 0)
	args := make([]any, 0)
	argID := 1

	if input.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name = $%d", argID))
		args = append(args, input.Name)
		argID++
	}
	if input.Surname != "" {
		setValues = append(setValues, fmt.Sprintf("surname = $%d", argID))
		args = append(args, input.Surname)
		argID++
	}
	if input.Patronymic != "" {
		setValues = append(setValues, fmt.Sprintf("patronymic = $%d", argID))
		args = append(args, input.Patronymic)
		argID++
	}
	if input.PassportSeries != "" {
		setValues = append(setValues, fmt.Sprintf("passport_series = $%d", argID))
		args = append(args, input.PassportSeries)
		argID++
	}
	if input.PassportNumber != "" {
		setValues = append(setValues, fmt.Sprintf("passport_number = $%d", argID))
		args = append(args, input.PassportNumber)
		argID++
	}
	if input.Address != "" {
		setValues = append(setValues, fmt.Sprintf("address = $%d", argID))
		args = append(args, input.Address)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", usersTable, setQuery, argID)
	args = append(args, userID)

	_, err = tx.Exec(r.ctx, query, args...)
	if err != nil {
		return err
	}

	err = tx.Commit(r.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersPostgres) Delete(userID int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, usersTable)

	commandTag, err := r.db.Exec(r.ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *UsersPostgres) GetAll(filters map[string]string, limit, offset int) (entities.GetAllUsers, error) {
	var users []entities.Users

	conditions := []string{}
	args := []any{}
	argID := 1

	for field, value := range filters {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", field, argID))
		args = append(args, value)
		argID++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`SELECT up.id, up.passport_series, up.passport_number, up.name, up.surname, up.patronymic, up.address FROM %s up%s LIMIT $%d OFFSET $%d`, usersTable, whereClause, argID, argID+1)
	args = append(args, limit, offset)

	conn, err := r.db.Acquire(r.ctx)
	if err != nil {
		return entities.GetAllUsers{}, err
	}
	defer conn.Release()

	rows, err := conn.Conn().Query(r.ctx, query, args...)
	if err != nil {
		return entities.GetAllUsers{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entities.Users
		if err := rows.Scan(&user.ID, &user.PassportSeries, &user.PassportNumber, &user.Name, &user.Surname, &user.Patronymic, &user.Address); err != nil {
			return entities.GetAllUsers{}, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return entities.GetAllUsers{}, err
	}

	// Count total users
	countQuery := fmt.Sprintf(`SELECT COUNT(up.id) FROM %s up%s`, usersTable, whereClause)
	var total int
	err = conn.Conn().QueryRow(r.ctx, countQuery, args[:len(args)-2]...).Scan(&total) // Skip limit and offset
	if err != nil {
		return entities.GetAllUsers{}, err
	}

	meta := entities.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	response := entities.GetAllUsers{
		Meta:  meta,
		Users: users,
	}

	return response, nil
}

func (r *UsersPostgres) GetByID(userID int) (entities.Users, error) {
	var user entities.Users
	query := fmt.Sprintf(`SELECT up.id, up.passport_series, up.passport_number, up.name, up.surname, up.patronymic, up.address FROM %s up WHERE up.id = $1`, usersTable)

	conn, err := r.db.Acquire(r.ctx)
	if err != nil {
		return entities.Users{}, err
	}
	defer conn.Release()

	pgxConn := conn.Conn()

	_, err = pgxConn.Prepare(r.ctx, "getUser", query)
	if err != nil {
		return user, err
	}

	err = pgxConn.QueryRow(r.ctx, "getUser", userID).Scan(&user.ID, &user.PassportSeries, &user.PassportNumber, &user.Name, &user.Surname, &user.Patronymic, &user.Address)
	if err != nil {
		return user, err
	}

	row := pgxConn.QueryRow(r.ctx, "getUser", userID)

	if err := row.Scan(&user.ID, &user.PassportSeries, &user.PassportNumber, &user.Name, &user.Surname, &user.Patronymic, &user.Address); err != nil {
		return user, err
	}

	return user, nil
}

//	func (r *UsersPostgres) Stats(userID int) (entities.UserStats, error) {
//		var userStats entities.UserStats
//
//		query := fmt.Sprintf(`SELECT ut.name, ut.surname, ut.patronymic, ut.passport_series, ut.passport_number, ut.address, t.task, t.description, t.overall_time
//			FROM %s ut LEFT JOIN %s t ON t.user_id = ut.id
//			WHERE ut.id = $1 AND t.overall_time IS NOT NULL ORDER BY t.overall_time DESC
//		`, usersTable, tasksTable)
//
//		conn, err := r.db.Acquire(r.ctx)
//		if err != nil {
//			return userStats, err
//		}
//		defer conn.Release()
//
//		pgxConn := conn.Conn()
//
//		_, err = pgxConn.Prepare(r.ctx, "getUserStats", query)
//		if err != nil {
//			return userStats, err
//		}
//
//		err = pgxConn.QueryRow(r.ctx, "getUserStats", userID).Scan(&userStats.Name, &userStats.Surname, &userStats.Patronymic, &userStats.PassportSeries, &userStats.PassportNumber, &userStats.Address, &userStats.Tasks.Task, &userStats.Tasks.Description, &userStats.OverallTime)
//
//		row := pgxConn.QueryRow(r.ctx, "getUserStats", userID)
//
//		if err := row.Scan(&userStats.Name, &userStats.Surname, &userStats.Patronymic, &userStats.PassportSeries, &userStats.PassportNumber, &userStats.Address, &userStats.Tasks.Task, &userStats.Tasks.Description, &userStats.OverallTime); err != nil {
//			return userStats, err
//		}
//		return userStats, nil
//	}
func (r *UsersPostgres) Stats(userID int) (entities.UserStats, error) {
	var userStats entities.UserStats
	var tasks []entities.Task

	query := fmt.Sprintf(`SELECT ut.name, ut.surname, ut.patronymic, ut.passport_series, ut.passport_number, ut.address, t.task, t.description, t.overall_time 
		FROM %s ut LEFT JOIN %s t ON t.user_id = ut.id 
		WHERE ut.id = $1 AND t.overall_time IS NOT NULL ORDER BY t.overall_time DESC `, usersTable, tasksTable)

	conn, err := r.db.Acquire(r.ctx)
	if err != nil {
		return userStats, err
	}
	defer conn.Release()

	rows, err := conn.Query(r.ctx, query, userID)
	if err != nil {
		return userStats, err
	}
	defer rows.Close()

	for rows.Next() {
		var task entities.Task
		err = rows.Scan(&userStats.Name, &userStats.Surname, &userStats.Patronymic, &userStats.PassportSeries, &userStats.PassportNumber, &userStats.Address, &task.Task, &task.Description, &task.OverallTime)
		if err != nil {
			return userStats, err
		}
		tasks = append(tasks, task)
	}

	userStats.Tasks = tasks

	// Calculate overall time
	var totalDuration time.Duration
	for _, task := range tasks {
		totalDuration += task.OverallTime
	}
	userStats.OverallTime = entities.FormatDuration(totalDuration)

	return userStats, nil
}
