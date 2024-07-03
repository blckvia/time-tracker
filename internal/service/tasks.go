package service

import (
	"time-tracker/internal/entities"
	"time-tracker/internal/repository"
)

type TasksService struct {
	repo repository.Tasks
}

func NewTasksService(repo repository.Tasks) *TasksService {
	return &TasksService{repo: repo}
}

func (s *TasksService) Create(input *entities.Task) (int, error) {
	return s.repo.Create(input)
}

func (s *TasksService) StartTask(taskID int) error {
	return s.repo.StartTask(taskID)
}

func (s *TasksService) StopTask(taskID int) error {
	return s.repo.StopTask(taskID)
}
