package handlers

import "github.com/AlcantuMurilo77/tic/internal/repository"

type UserHandler struct {
	userRepository *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepository: repo,
	}
}
