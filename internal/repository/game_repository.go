package repository

import (
	"context"
	"github.com/AlcantuMurilo77/tic/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
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

func (r *GameRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Game, error) {

	var game models.Game

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&game)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *GameRepository) UpdateBoard(
	ctx context.Context,
	gameID uuid.UUID,
	board [][]int,
	currentTurn uuid.UUID,
) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": gameID},
		bson.M{
			"$set": bson.M{
				"board":        board,
				"current_turn": currentTurn,
			},
		},
	)

	return err
}

func (r *GameRepository) Join(
	ctx context.Context,
	gameID uuid.UUID,
	userO uuid.UUID,
	currentTurn uuid.UUID,
	startedAt time.Time,
) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": gameID},
		bson.M{
			"$set": bson.M{
				"user_o":       userO,
				"status":       models.GameReady,
				"current_turn": currentTurn,
				"started_at":   startedAt,
			},
		},
	)
	return err
}

func (r *GameRepository) UpdateState(
	ctx context.Context,
	game *models.Game,
) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": game.GameUuid},
		bson.M{
			"$set": bson.M{
				"board":        game.Board,
				"current_turn": game.CurrentTurn,
				"winner_id":    game.WinnerID,
				"status":       game.Status,
				"ended_at":     game.EndedAt,
			},
		},
	)
	return err
}
