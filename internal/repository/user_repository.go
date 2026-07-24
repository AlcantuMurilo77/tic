package repository

import (
	"context"
	"fmt"
	"github.com/AlcantuMurilo77/tic/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("user"),
	}
}

func (r UserRepository) Create(ctx context.Context, user *models.User) error {

	if user.UserUuid == uuid.Nil {
		user.UserUuid = uuid.New()
	}
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	users := make([]models.User, 0)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil

}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid game id:  %w", err)
	}

	var user models.User

	err = r.collection.FindOne(ctx, bson.M{
		"_id": parsedID,
	}).Decode(&user)

	if err != nil {
		return nil, fmt.Errorf("find game: %w", err)
	}
	return &user, nil
}
