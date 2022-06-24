package db

import (
	"context"
	"strings"
	"time"

	"github.com/sushilman/userservice/models"
	"github.com/sushilman/userservice/usererrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	QUERY_TIMEOUT = 5 * time.Second
	COLLECTION   = "users"
)

type IUserStorage interface {
	Insert(models.User) error
	GetAll(models.GetUserQueryParams) ([]models.User, error)
	GetById(string) (*models.User, error)
	Update(models.User) error
	DeleteById(string) error
}

// implements the IStorage interface
type userstorage struct {
	database *mongo.Database
}

func NewStorage(database *mongo.Database) IUserStorage {
	return &userstorage{
		database,
	}
}

func (s *userstorage) Insert(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), QUERY_TIMEOUT)
	defer cancel()

	_, err := s.database.Collection(COLLECTION).InsertOne(ctx, user)

	return err // 'nil' if successful
}

func (s *userstorage) GetAll(queryParams models.GetUserQueryParams) (users []models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), QUERY_TIMEOUT)
	defer cancel()

	filter := bson.M{}
	if queryParams.Country != "" {
		filter["country"] = strings.ToUpper(queryParams.Country)
	}
	if queryParams.FirstName != "" {
		filter["first_name"] = queryParams.FirstName
	}
	if queryParams.LastName != "" {
		filter["last_name"] = queryParams.LastName
	}
	if queryParams.NickName != "" {
		filter["nickname"] = queryParams.NickName
	}
	if queryParams.Email != "" {
		filter["email"] = queryParams.Email
	}

	opts := options.Find()
	opts.SetSkip(int64(queryParams.Offset))
	opts.SetLimit(int64(queryParams.Limit))

	cursor, err := s.database.Collection(COLLECTION).Find(ctx, filter, opts)
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

func (s *userstorage) GetById(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), QUERY_TIMEOUT)
	defer cancel()

	user := models.User{}
	query := bson.M{"id": id}
	err := s.database.Collection(COLLECTION).FindOne(ctx, query).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return nil, &usererrors.NotFoundError{}
	}

	return &user, err
}

func (s *userstorage) Update(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), QUERY_TIMEOUT)
	defer cancel()

	filter := bson.M{"id": user.Id}
	userUpdate := bson.M{"$set": user}

	updateResult, err := s.database.Collection(COLLECTION).UpdateOne(ctx, filter, userUpdate)

	if err != nil {
		return err
	}

	// No match was found for the given ID
	if updateResult.MatchedCount == 0 {
		return &usererrors.NotFoundError{}
	}

	return nil
}

func (s *userstorage) DeleteById(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), QUERY_TIMEOUT)
	defer cancel()

	query := bson.M{"id": id}
	_, err := s.database.Collection(COLLECTION).DeleteOne(ctx, query)

	if err != nil {
		return err
	}

	return nil
}
