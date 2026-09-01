package services

import (
	"context"
	"fmt"

	"github.com/AlcantuMurilo77/tic/game"
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
) (*models.Game, error) {

	game := &models.Game{
		GameUuid:    uuid.New(),
		UserX:       userX,
		UserO:       uuid.Nil,
		CurrentTurn: userX,
		Status:      models.GameWaiting,
		StartedAt:   time.Now(),
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

func (s *GameService) FindByID(ctx context.Context, id uuid.UUID) (*models.Game, error) {
	game, err := s.gameRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func boardIsFull(board [][]int) bool {
	for _, row := range board {
		for _, cell := range row {
			if cell == 0 {
				return false
			}
		}
	}
	return true
}

func (s *GameService) Move(
	ctx context.Context,
	gameID uuid.UUID,
	playerID uuid.UUID,
	row int,
	col int,
) (*models.Game, error) {

	gameModel, err := s.gameRepository.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if playerID != gameModel.UserX && playerID != gameModel.UserO {
		return nil, fmt.Errorf("player does not belong to this game")
	}

	if playerID != gameModel.CurrentTurn {
		return nil, fmt.Errorf("not your turn")
	}

	if row < 0 || row >= 3 || col < 0 || col >= 3 {
		return nil, fmt.Errorf("invalid position")
	}

	player := 1
	if playerID == gameModel.UserO {
		player = 2
	}

	tic := game.NewTicTacToeFromBoard(gameModel.Board)

	if !tic.CheckIfLegalMove(row, col) {
		return nil, fmt.Errorf("position already occupied")
	}

	won := tic.Move(row, col, player)

	gameModel.Board = tic.Board

	if won {
		gameModel.WinnerID = playerID
		gameModel.Status = models.GameFinished
		gameModel.EndedAt = time.Now()
	} else if boardIsFull(gameModel.Board) {
		gameModel.WinnerID = uuid.Nil
		gameModel.Status = models.GameFinished
		gameModel.EndedAt = time.Now()
	}

	if gameModel.Status != models.GameFinished {
		if playerID == gameModel.UserX {
			gameModel.CurrentTurn = gameModel.UserO
		} else {
			gameModel.CurrentTurn = gameModel.UserX
		}
	}

	if err := s.gameRepository.UpdateState(
		ctx,
		gameModel,
	); err != nil {
		return nil, err
	}

	return gameModel, nil
}

func (s *GameService) Join(
	ctx context.Context,
	gameID uuid.UUID,
	playerID uuid.UUID,
) (*models.Game, error) {

	game, err := s.gameRepository.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if game.Status != models.GameWaiting {
		return nil, fmt.Errorf("game is not waiting for a player")
	}

	if game.UserX == playerID {
		return nil, fmt.Errorf("creator cannot join their own game")
	}

	if game.UserO != uuid.Nil {
		return nil, fmt.Errorf("game already has a second player")
	}

	game.UserO = playerID
	game.CurrentTurn = game.UserX
	game.StartedAt = time.Now()
	game.Status = models.GameReady

	if err := s.gameRepository.Join(
		ctx,
		gameID,
		playerID,
		game.UserX,
		game.StartedAt,
	); err != nil {
		return nil, err
	}
	return game, nil
}
