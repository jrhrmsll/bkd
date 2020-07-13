package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"bookmarks/internal"
)

const (
	bookmarksDatabase   = "bkd"
	bookmarksCollection = "bookmarks"
)

type BookmarksRepository struct {
	ctx    context.Context
	client *mongo.Client
	logger *log.Logger
}

func NewBookmarksRepository(
	ctx context.Context,
	client *mongo.Client,
	logger *log.Logger,
) *BookmarksRepository {
	return &BookmarksRepository{
		ctx:    ctx,
		client: client,
		logger: logger,
	}
}

func (repository *BookmarksRepository) collection() *mongo.Collection {
	return repository.client.Database(bookmarksDatabase).Collection(bookmarksCollection)
}

func (repository *BookmarksRepository) Store(bookmark *internal.Bookmark) (string, error) {
	result, err := repository.collection().InsertOne(repository.ctx, bookmark)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", result.InsertedID), nil
}

func (repository *BookmarksRepository) Find() ([]*internal.Bookmark, error) {
	cursor, err := repository.collection().Find(repository.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(repository.ctx)

	bookmarks := make([]*internal.Bookmark, 0)

	err = cursor.All(repository.ctx, &bookmarks)
	if err != nil {
		return nil, err
	}

	return bookmarks, nil
}

func (repository *BookmarksRepository) Delete(version string) (bool, error) {
	filter := bson.D{
		{
			Key: "version",
			Value: bson.D{{
				Key:   "$ne",
				Value: version,
			}},
		},
		{
			Key:   "mode",
			Value: "auto",
		},
	}

	result, err := repository.collection().DeleteMany(repository.ctx, filter)

	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}
