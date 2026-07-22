package repository

import (
	"context"
	"github.com/AlcantuMurilo77/tic/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GameRepository struct {
	collection *mongo.Collection
}

func NewGameRepository(db *mongo.Database) *GameRepository {
	return &GameRepository{
		collection: db.Collection("game"),
	}
}

func (r GameRepository) Create(ctx context.Context, game *models.Game) error {
	_, err := r.collection.InsertOne(ctx, game)
	if err != nil {
		return err
	}
	return nil
}
