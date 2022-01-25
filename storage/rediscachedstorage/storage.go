package rediscachedstorage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"miniurl/storage"
	"time"
)

const cacheTTL = 10 * time.Second

func NewStorage(redisUrl string, persistenceStorage storage.Storage) *Storage {
	return &Storage{
		persistence: persistenceStorage,
		client:      redis.NewClient(&redis.Options{Addr: redisUrl}),
	}
}

type Storage struct {
	persistence storage.Storage
	client      *redis.Client
}

var _ storage.Storage = (*Storage)(nil)

func (s *Storage) PutURL(ctx context.Context, url storage.ShortedURL) (storage.URLKey, error) {
	key, err := s.persistence.PutURL(ctx, url)
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

	url, err := s.persistence.GetURL(ctx, key)
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
