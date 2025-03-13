package repository

import (
	"context"
	"errors"
	"github.com/vkhrushchev/GophKeeper/internal/server/entity"
	"github.com/vkhrushchev/GophKeeper/internal/server/usecase"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log/slog"
)

type MongoUserRepository struct {
	log         *slog.Logger
	mongoClient *mongo.Client
}

func (r *MongoUserRepository) Save(ctx context.Context, userEntity entity.User) (entity.User, error) {
	collection := r.mongoClient.Database("GophKeeper").Collection("users")

	count, err := collection.CountDocuments(ctx, bson.M{
		"username": userEntity.Username,
		"email":    userEntity.Email,
	})

	if err != nil {
		r.log.Error("repository: unexpected error", "error", err)
		return entity.User{}, usecase.ErrUnexpected
	}

	if count > 0 {
		return userEntity, usecase.ErrUserAlreadyExists
	}

	result, err := collection.InsertOne(ctx, userEntity)
	if err != nil {
		r.log.Error("repository: unexpected error", "error", err)
		return entity.User{}, usecase.ErrUnexpected
	}

	userEntity.ID = result.InsertedID.(string)

	return userEntity, nil
}

func (r *MongoUserRepository) Get(ctx context.Context, username string) (entity.User, error) {
	collection := r.mongoClient.Database("GophKeeper").Collection("users")

	var result entity.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&result)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		r.log.Debug("repository: user not found", "username", username)
		return entity.User{}, usecase.ErrUserNotFound
	} else if err != nil {
		r.log.Error("repository: unexpected error", "error", err)
		return entity.User{}, usecase.ErrUnexpected
	}

	return result, nil
}
