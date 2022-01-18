package storage

import (
	"context"
	"errors"
	"fmt"
)

var (
	StorageError   = errors.New("storage")
	ErrorCollision = fmt.Errorf("%w.collision", StorageError)
	ErrorNotFound  = fmt.Errorf("%w.not_found", StorageError)
)

type ShortedURL string
type URLKey string

type Storage interface {
	PutURL(ctx context.Context, url ShortedURL) (URLKey, error)
	GetURL(ctx context.Context, key URLKey) (ShortedURL, error)
}
