package service

import (
	"time-tracker/internal/entities"
	"time-tracker/internal/repository"
)

type UsersService struct {
	repo repository.Users
}

func NewUsersService(repo repository.Users) *UsersService {
	return &UsersService{repo: repo}
}

func (s *UsersService) Create(input *entities.Users) (int, error) {
	return s.repo.Create(input)
}

func (s *UsersService) Update(userID int, input *entities.Users) error {
	return s.repo.Update(userID, input)
}

func (s *UsersService) Delete(userID int) error {
	return s.repo.Delete(userID)
}

func (s *UsersService) GetAll(filters map[string]any, limit, offset int) (entities.GetAllUsers, error) {
	return s.repo.GetAll(filters, limit, offset)
}

func (s *UsersService) GetByID(userID int) (entities.Users, error) {
	return s.repo.GetByID(userID)
}
