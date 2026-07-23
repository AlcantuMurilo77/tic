package services

import (
	"context"
	"github.com/AlcantuMurilo77/tic/internal/models"
	"github.com/AlcantuMurilo77/tic/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: repo,
	}
}

func (s *UserService) Create(
	ctx context.Context,
	name string,
	country string,
	xman bool,
) (*models.User, error) {

	user := &models.User{
		UserUuid: uuid.New(),
		Name:     name,
		Country:  country,
		Xman:     xman,
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindAll(ctx context.Context) ([]models.User, error) {
	users, err := s.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
