package rediscachedstorage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"miniurl/storage"
	"time"
)

const cacheTTL = 10 * time.Second

func NewStorage(persistentStorage storage.Storage, client *redis.Client) *Storage {
	return &Storage{
		client:            client,
		persistentStorage: persistentStorage,
	}
}

type Storage struct {
	client            *redis.Client
	persistentStorage storage.Storage
}

//TODO: planning to add mongo sharding
var _ storage.Storage = (*Storage)(nil)

func (s *Storage) PutURL(ctx context.Context, url storage.ShortedURL) (storage.URLKey, error) {
	key, err := s.persistentStorage.PutURL(ctx, url)
	if err != nil {
		return "", err
	}
	fullKey := s.fullKey(key)
	resp := s.client.Set(ctx, fullKey, string(url), cacheTTL)
	if err := resp.Err(); err != nil {
		log.Printf("Failed to save key %s to redis", fullKey)
		return "", err
	}

	return key, nil
}

func (s *Storage) GetURL(ctx context.Context, key storage.URLKey) (storage.ShortedURL, error) {
	fullKey := s.fullKey(key)
	rawUrl, err := s.client.Get(ctx, fullKey).Result()
	switch {
	case err == redis.Nil:
	// go to persistence
	case err != nil:
		return "", err
	default:
		log.Println("Successfully loaded key from cache")
		return storage.ShortedURL(rawUrl), nil
	}

	url, err := s.persistentStorage.GetURL(ctx, key)
	if err != nil {
		return "", err
	}

	resp := s.client.Set(ctx, fullKey, string(url), cacheTTL)
	if err := resp.Err(); err != nil {
		log.Printf("Failed to save key %s to redis", fullKey)
		return "", err
	}

	log.Println("Successfully loaded key from persistence")
	return url, nil
}

func (s *Storage) fullKey(key storage.URLKey) string {
	return "su:" + string(key)
}
