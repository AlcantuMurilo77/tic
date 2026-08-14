package services

import (
	"context"
	"fmt"
	"github.com/AlcantuMurilo77/tic/internal/models"
	"github.com/AlcantuMurilo77/tic/internal/repository"
	"github.com/google/uuid"
	"time"
)

type GameService struct {
	gameRepository *repository.GameRepository
}

func NewGameService(repo *repository.GameRepository) *GameService {
	return &GameService{
		gameRepository: repo,
	}
}

func (s *GameService) Create(
	ctx context.Context,
	userX uuid.UUID,
	userO uuid.UUID,
) (*models.Game, error) {

	game := &models.Game{
		GameUuid:  uuid.New(),
		UserX:     userX,
		UserO:     userO,
		StartedAt: time.Now(),
		Board: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
	}

	if err := s.gameRepository.Create(ctx, game); err != nil {
		return nil, err
	}

	return game, nil
}

func (s *GameService) FindAll(ctx context.Context) ([]models.Game, error) {
	games, err := s.gameRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (s *GameService) Move(
	ctx context.Context,
	gameID uuid.UUID,
	playerID uuid.UUID,
	row int,
	col int,
) (*models.Game, error) {

	game, err := s.gameRepository.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if playerID != game.UserX && playerID != game.UserO {
		return nil, fmt.Errorf("player does not belong to this game")
	}

	return game, nil
}
