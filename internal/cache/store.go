package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/gabehf/koito/internal/memkv"
)

type Store interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

type memKV interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, expiration ...time.Duration)
}

type MemKVStore struct {
	store memKV
}

func NewDefaultStore() Store {
	return NewMemKVStore(memkv.Store)
}

func NewMemKVStore(store memKV) *MemKVStore {
	return &MemKVStore{store: store}
}

func (s *MemKVStore) Get(_ context.Context, key string) ([]byte, bool, error) {
	v, ok := s.store.Get(key)
	if !ok {
		return nil, false, nil
	}

	b, ok := v.([]byte)
	if !ok {
		return nil, false, fmt.Errorf("MemKVStore.Get: expected []byte for key %s", key)
	}

	return append([]byte(nil), b...), true, nil
}

func (s *MemKVStore) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	s.store.Set(key, append([]byte(nil), value...), ttl)
	return nil
}
