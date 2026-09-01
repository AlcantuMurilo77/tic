package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlcantuMurilo77/tic/internal/models"
	"github.com/AlcantuMurilo77/tic/internal/repository"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RematchService struct {
	rematchRepository *repository.RematchRepository
	gameRepository    *repository.GameRepository
}

func NewRematchService(
	rematchRepository *repository.RematchRepository,
	gameRepository *repository.GameRepository,
) *RematchService {
	return &RematchService{
		rematchRepository: rematchRepository,
		gameRepository:    gameRepository,
	}
}

func (s *RematchService) Request(
	ctx context.Context,
	gameID uuid.UUID,
	playerID uuid.UUID,
) (*models.RematchRequest, error) {
	game, err := s.gameRepository.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if game.Status != models.GameFinished {
		return nil, fmt.Errorf("game is not finished")
	}

	if playerID != game.UserX && playerID != game.UserO {
		return nil, fmt.Errorf("player does not belong to this game")
	}

	existingRequest, err := s.rematchRepository.FindByOriginalGameID(
		ctx, gameID,
	)

	if err == nil {
		return existingRequest, nil
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	request := &models.RematchRequest{
		RematchUuid:         uuid.New(),
		OriginalGameId:      gameID,
		RequestedByPlayerId: playerID,
		Status:              models.RematchPending,
		CreatedAt:           time.Now(),
	}

	if err := s.rematchRepository.Create(
		ctx,
		request,
	); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return s.rematchRepository.FindByOriginalGameID(ctx, gameID)
		}
		return nil, err
	}

	return request, nil
}

func (s *RematchService) Accept(
	ctx context.Context,
	rematchID uuid.UUID,
	playerID uuid.UUID,
) (*models.Game, error) {
	request, err := s.rematchRepository.FindByID(ctx, rematchID)
	if err != nil {
		return nil, err
	}

	if request.Status == models.RematchAccepted {
		if request.NewGameId == nil {
			return nil, fmt.Errorf("accepted rematch has no game")
		}

		return s.gameRepository.FindByID(
			ctx,
			*request.NewGameId,
		)

	}

	originalGame, err := s.gameRepository.FindByID(
		ctx,
		request.OriginalGameId,
	)

	if err != nil {
		return nil, err
	}

	if playerID != originalGame.UserX && playerID != originalGame.UserO {
		return nil, fmt.Errorf("player does not belong to this game")
	}

	if playerID == request.RequestedByPlayerId {
		return nil, fmt.Errorf("player cannot accept their own rematch")
	}

	newGameID := uuid.New()
	now := time.Now()

	if err := s.rematchRepository.ClaimAcceptance(
		ctx,
		request.RematchUuid,
		playerID,
		newGameID,
		now,
	); err != nil {
		current, findErr := s.rematchRepository.FindByID(ctx, rematchID)
		if findErr != nil {
			return nil, err
		}
		if current.Status == models.RematchAccepted && current.NewGameId != nil {
			return s.gameRepository.FindByID(ctx, *current.NewGameId)
		}
		return nil, err
	}

	originalGameID := originalGame.GameUuid
	newGame := &models.Game{
		GameUuid:    newGameID,
		UserX:       originalGame.UserO,
		UserO:       originalGame.UserX,
		CurrentTurn: originalGame.UserO,
		Status:      models.GameReady,
		StartedAt:   time.Now(),
		Board: [][]int{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		RematchOfGameID: &originalGameID,
	}

	if err := s.gameRepository.Create(
		ctx,
		newGame,
	); err != nil {
		_ = s.rematchRepository.ReleaseAcceptance(
			ctx,
			request.RematchUuid,
			newGameID,
		)
		return nil, err
	}

	return newGame, nil

}
