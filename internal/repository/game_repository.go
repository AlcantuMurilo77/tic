package repository

import (
	"context"
	"fmt"
	"github.com/AlcantuMurilo77/tic/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	if game.GameUuid == uuid.Nil {
		game.GameUuid = uuid.New()
	}
	_, err := r.collection.InsertOne(ctx, game)
	if err != nil {
		return err
	}
	return nil
}

func (r *GameRepository) FindAll(ctx context.Context) ([]models.Game, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	games := make([]models.Game, 0)

	if err := cursor.All(ctx, &games); err != nil {
		return nil, err
	}

	return games, nil

}

func (r *GameRepository) FindByID(ctx context.Context, id string) (*models.Game, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid game id:  %w", err)
	}

	var game models.Game

	err = r.collection.FindOne(ctx, bson.M{
		"_id": parsedID,
	}).Decode(&game)

	if err != nil {
		return nil, fmt.Errorf("find game: %w", err)
	}
	return &game, nil
}
