package db

import (
	"context"
	"errors"
	"time"

	"github.com/sushilman/userservice/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	queryTimeout = 5 * time.Second
	COLLECTION   = "users"
)

type IStorage interface {
	Insert(context.Context, models.User) error
	GetAll(context.Context) ([]models.User, error)
	GetById(context.Context, string) (*models.User, error)
	Update(context.Context, models.User) (bool, error)
}

// implements the IStorage interface
type storage struct {
	database *mongo.Database
}

func NewStorage(database *mongo.Database) IStorage {
	return &storage{
		database,
	}
}

func (s *storage) Insert(ctx context.Context, user models.User) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	_, err := s.database.Collection(COLLECTION).InsertOne(ctx, user)

	return err // 'nil' if successful
}

//TODO: implement filtering by query params
func (s *storage) GetAll(ctx context.Context) (users []models.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	filter := bson.M{}
	cursor, err := s.database.Collection(COLLECTION).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	errFetchAll := cursor.All(ctx, &users)
	if errFetchAll != nil {
		return nil, errFetchAll
	}

	return users, nil
}

func (s *storage) GetById(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	user := models.User{}
	query := bson.M{"id": id}
	err := s.database.Collection(COLLECTION).FindOne(ctx, query).Decode(&user)

	// return empty document if not found
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("not_found")
	}

	return &user, err
}

func (s *storage) Update(ctx context.Context, user models.User) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	filter := bson.M{"id": user.Id}
	userUpdate := bson.M{"$set": user}

	updateResult, err := s.database.Collection(COLLECTION).UpdateOne(ctx, filter, userUpdate)

	if err != nil {
		return false, err
	}

	// No match was found for the given ID
	if updateResult.MatchedCount == 0 {
		return false, nil
	}

	return true, nil
}