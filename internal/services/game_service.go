package services

import (
	"context"
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

func (s *GameService) FindByID(ctx context.Context, id string) (*models.Game, error) {
	return s.gameRepository.FindByID(ctx, id)
}
