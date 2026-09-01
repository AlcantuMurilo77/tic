package repository

import (
	"context"
	"fmt"
	"github.com/AlcantuMurilo77/tic/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type RematchRepository struct {
	collection *mongo.Collection
}

func (r *RematchRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "original_game_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func NewRematchRepository(
	db *mongo.Database,
) *RematchRepository {
	return &RematchRepository{
		collection: db.Collection("rematch_requests"),
	}
}

func (r *RematchRepository) Create(
	ctx context.Context,
	request *models.RematchRequest,
) error {
	_, err := r.collection.InsertOne(ctx, request)
	return err
}

func (r *RematchRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.RematchRequest, error) {
	var request models.RematchRequest

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&request)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *RematchRepository) FindByOriginalGameID(
	ctx context.Context,
	gameID uuid.UUID,
) (*models.RematchRequest, error) {
	var request models.RematchRequest

	err := r.collection.FindOne(
		ctx,
		bson.M{"original_game_id": gameID},
	).Decode(&request)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *RematchRepository) ClaimAcceptance(
	ctx context.Context,
	requestID uuid.UUID,
	acceptedBy uuid.UUID,
	newGameID uuid.UUID,
	acceptedAt time.Time,
) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":    requestID,
			"status": models.RematchPending,
		},
		bson.M{
			"$set": bson.M{
				"status":                models.RematchAccepted,
				"accepted_by_player_id": acceptedBy,
				"new_game_id":           newGameID,
				"accepted_at":           acceptedAt,
			},
		},
	)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("rematch is no longer pending")
	}

	return nil
}

func (r *RematchRepository) ReleaseAcceptance(
	ctx context.Context,
	requestID uuid.UUID,
	newGameID uuid.UUID,
) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":         requestID,
			"status":      models.RematchAccepted,
			"new_game_id": newGameID,
		},
		bson.M{
			"$set": bson.M{"status": models.RematchPending},
			"$unset": bson.M{
				"accepted_by_player_id": "",
				"new_game_id":           "",
				"accepted_at":           "",
			},
		},
	)
	return err
}
