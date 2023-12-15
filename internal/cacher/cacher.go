package cacher

import (
	"context"
	"time"
)

type Cacher interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
	Exist(ctx context.Context, key string) (bool, error)
}
