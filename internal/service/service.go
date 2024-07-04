package service

import (
	"time"

	"time-tracker/internal/entities"
	"time-tracker/internal/repository"
)

type Users interface {
	Create(input *entities.Users) (int, error)
	Update(userID int, input *entities.Users) error
	Delete(userID int) error
	GetAll(filters map[string]string, limit, offset int) (entities.GetAllUsers, error)
	GetByID(userID int) (entities.Users, error)
	Stats(userID int) (entities.UserStats, error)
}

type Tasks interface {
	Create(input *entities.Task, userID int) (int, error)
	StartTask(taskID int) error
	StopTask(taskID int) (time.Duration, error)
}

type Service struct {
	Users
	Tasks
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Users: NewUsersService(repos.Users),
		Tasks: NewTasksService(repos.Tasks),
	}
}
