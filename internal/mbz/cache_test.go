package mbz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gabehf/koito/internal/cache"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetReleaseGroup_CacheHitSkipsSecondRequest(t *testing.T) {
	id := uuid.New()

	var mu sync.Mutex
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"title":"cached-rg","releases":[]}`))
	}))
	defer server.Close()

	client := newMusicBrainzClientWithCache(server.URL, cache.NewDefaultStore())
	defer client.Shutdown()

	ctx := context.Background()

	first, err := client.GetReleaseGroup(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, first)
	assert.Equal(t, "cached-rg", first.Title)

	second, err := client.GetReleaseGroup(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, second)
	assert.Equal(t, "cached-rg", second.Title)

	mu.Lock()
	count := requestCount
	mu.Unlock()
	assert.Equal(t, 1, count)
}

type testStore struct {
	mu   sync.RWMutex
	data map[string]testItem
}

type testItem struct {
	value     []byte
	expiresAt time.Time
}

func newTestStore() *testStore {
	return &testStore{data: make(map[string]testItem)}
}

func (s *testStore) Get(_ context.Context, key string) ([]byte, bool, error) {
	s.mu.RLock()
	it, ok := s.data[key]
	s.mu.RUnlock()
	if !ok {
		return nil, false, nil
	}
	if !it.expiresAt.IsZero() && time.Now().After(it.expiresAt) {
		s.mu.Lock()
		delete(s.data, key)
		s.mu.Unlock()
		return nil, false, nil
	}
	return append([]byte(nil), it.value...), true, nil
}

func (s *testStore) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = testItem{
		value:     append([]byte(nil), value...),
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

func (s *testStore) ttl(key string) (time.Duration, bool) {
	s.mu.RLock()
	it, ok := s.data[key]
	s.mu.RUnlock()
	if !ok {
		return 0, false
	}
	return time.Until(it.expiresAt), true
}

func TestGetRelease_UsesExpectedTTL(t *testing.T) {
	id := uuid.New()
	store := newTestStore()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"title":"x","id":"` + id.String() + `"}`))
	}))
	defer server.Close()

	client := newMusicBrainzClientWithCache(server.URL, store)
	defer client.Shutdown()

	_, err := client.GetRelease(context.Background(), id)
	require.NoError(t, err)

	ttl, ok := store.ttl(mbzCacheKey("release", id))
	require.True(t, ok)
	assert.GreaterOrEqual(t, ttl, releaseCacheTTL-time.Minute)
	assert.LessOrEqual(t, ttl, releaseCacheTTL+time.Minute)
}

func TestGetArtistPrimaryAliases_CacheHitSkipsSecondRequest(t *testing.T) {
	id := uuid.New()

	var mu sync.Mutex
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"artist-a","aliases":[{"name":"artist-a-primary","primary":true}],"genres":[]}`))
	}))
	defer server.Close()

	client := newMusicBrainzClientWithCache(server.URL, cache.NewDefaultStore())
	defer client.Shutdown()

	ctx := context.Background()

	first, err := client.GetArtistPrimaryAliases(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, []string{"artist-a", "artist-a-primary"}, first)

	second, err := client.GetArtistPrimaryAliases(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, []string{"artist-a", "artist-a-primary"}, second)

	mu.Lock()
	count := requestCount
	mu.Unlock()
	assert.Equal(t, 1, count)
}
