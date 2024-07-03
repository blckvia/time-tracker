package service

import (
	"time"

	"time-tracker/internal/entities"
	"time-tracker/internal/repository"
)

type TasksService struct {
	repo repository.Tasks
}

func NewTasksService(repo repository.Tasks) *TasksService {
	return &TasksService{repo: repo}
}

func (s *TasksService) Create(input *entities.Task, userID int) (int, error) {
	return s.repo.Create(input, userID)
}

func (s *TasksService) StartTask(taskID int) error {
	return s.repo.StartTask(taskID)
}

func (s *TasksService) StopTask(taskID int) (time.Duration, error) {
	return s.repo.StopTask(taskID)
}
