package handlers

import "github.com/AlcantuMurilo77/tic/internal/repository"

type GameHandler struct {
	gameRepository *repository.GameRepository
}

func NewGameHandler(repo *repository.GameRepository) *GameHandler {
	return &GameHandler{
		gameRepository: repo,
	}
}
