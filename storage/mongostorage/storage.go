package mongostorage

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"miniurl/generator"
	storage2 "miniurl/storage"
	"time"
)

const dbName = "shortUrls"
const collectionName = "urls"

type storage struct {
	urls *mongo.Collection
}

func NewStorage(mongoUrl string) *storage {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		panic(err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	ensureIndexes(ctx, collection)

	return &storage{
		urls: collection,
	}
}

func ensureIndexes(ctx context.Context, collection *mongo.Collection) {
	indexModels := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{Key: "_id", Value: bsonx.Int32(1)}},
		},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	_, err := collection.Indexes().CreateMany(ctx, indexModels, opts)
	if err != nil {
		panic(fmt.Errorf("failed to ensure indexes - %w", err))
	}
}

func (s *storage) PutURL(ctx context.Context, url storage2.ShortedURL) (storage2.URLKey, error) {
	for attempt := 0; attempt < 5; attempt++ {
		key := storage2.URLKey(generator.GetRandomKey())
		item := urlItem{
			Key: key,
			Url: url,
		}
		_, err := s.urls.InsertOne(ctx, item)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				continue
			}
			return "", fmt.Errorf("something went wrong - %w", storage2.StorageError)
		}

		return key, nil
	}
	return "", fmt.Errorf("too much attempts during inserting - %w", storage2.ErrorCollision)
}

func (s *storage) GetURL(ctx context.Context, key storage2.URLKey) (storage2.ShortedURL, error) {
	var result urlItem
	err := s.urls.FindOne(ctx, bson.M{"_id": key}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", fmt.Errorf("no document with key %v - %w", key, storage2.ErrorNotFound)
		}
		return "", fmt.Errorf("something went wrong - %w", storage2.StorageError)
	}
	return result.Url, nil
}

type urlItem struct {
	Key storage2.URLKey     `bson:"_id"`
	Url storage2.ShortedURL `bson:"url"`
}
