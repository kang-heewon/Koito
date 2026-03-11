// package discogs provides functions for interacting with the Discogs API
package discogs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gabehf/koito/internal/cache"
	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/queue"
)

const (
	discogsBaseUrl     = "https://api.discogs.com"
	discogsCachePrefix = "discogs"
	releaseCacheTTL    = 24 * time.Hour
	searchCacheTTL     = 6 * time.Hour
)

// DiscogsClient handles communication with Discogs API
type DiscogsClient struct {
	consumerKey    string
	consumerSecret string
	userAgent      string
	requestQueue   *queue.RequestQueue
	cacheStore     cache.Store
}

// DiscogsSearchResult represents a search result from Discogs
type DiscogsSearchResult struct {
	Results []DiscogsSearchItem `json:"results"`
}

// DiscogsSearchItem represents a single search result item
type DiscogsSearchItem struct {
	ID     int      `json:"id"`
	Type   string   `json:"type"` // release, master, artist
	Title  string   `json:"title"`
	Genres []string `json:"genre"`
	Styles []string `json:"style"`
	Year   int      `json:"year"`
}

// DiscogsRelease represents a release from Discogs
type DiscogsRelease struct {
	ID      int             `json:"id"`
	Title   string          `json:"title"`
	Artists []DiscogsArtist `json:"artists"`
	Genres  []string        `json:"genres"`
	Styles  []string        `json:"styles"`
	Year    int             `json:"year"`
}

// DiscogsArtist represents an artist from Discogs
type DiscogsArtist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Anv  string `json:"anv"` // Artist Name Variation
}

// DiscogsCaller interface for Discogs API operations
type DiscogsCaller interface {
	SearchRelease(ctx context.Context, artist, title string) (*DiscogsSearchResult, error)
	GetRelease(ctx context.Context, releaseID int) (*DiscogsRelease, error)
	GetReleaseGenres(ctx context.Context, releaseID int) ([]string, error)
	Shutdown()
}

// NewDiscogsClient creates a new Discogs client
func NewDiscogsClient() *DiscogsClient {
	ret := new(DiscogsClient)
	ret.consumerKey = cfg.DiscogsConsumerKey()
	ret.consumerSecret = cfg.DiscogsConsumerSecret()
	ret.userAgent = cfg.UserAgent()
	ret.requestQueue = queue.NewRequestQueue(1, 1) // Discogs rate limit: 1 req/sec
	ret.cacheStore = cache.NewDefaultStore()
	return ret
}

// NewDiscogsClientWithCache creates a Discogs client with custom cache (for testing)
func NewDiscogsClientWithCache(cacheStore cache.Store) *DiscogsClient {
	ret := new(DiscogsClient)
	ret.consumerKey = "test-key"
	ret.consumerSecret = "test-secret"
	ret.userAgent = "koito-test"
	ret.requestQueue = queue.NewRequestQueue(100, 100)
	ret.cacheStore = cacheStore
	return ret
}

func (c *DiscogsClient) Shutdown() {
	c.requestQueue.Shutdown()
}

func discogsCacheKey(entity string, id interface{}) string {
	return fmt.Sprintf("%s:%s:%v", discogsCachePrefix, entity, id)
}

// SearchRelease searches for releases by artist and title
func (c *DiscogsClient) SearchRelease(ctx context.Context, artist, title string) (*DiscogsSearchResult, error) {
	l := logger.FromContext(ctx)

	// Check cache first
	cacheKey := discogsCacheKey("search", url.QueryEscape(artist+"-"+title))
	if c.cacheStore != nil {
		body, found, err := c.cacheStore.Get(ctx, cacheKey)
		if err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to read Discogs cache entry")
		} else if found {
			var result DiscogsSearchResult
			if err := json.Unmarshal(body, &result); err == nil {
				return &result, nil
			}
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to unmarshal Discogs cache entry")
		}
	}

	// Build search query
	query := fmt.Sprintf("artist:\"%s\" release:\"%s\"", artist, title)
	searchURL := fmt.Sprintf("%s/database/search?q=%s&type=release", discogsBaseUrl, url.QueryEscape(query))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("SearchRelease: failed to build request: %w", err)
	}

	c.setAuthHeaders(req)

	body, err := c.queue(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("SearchRelease: request failed: %w", err)
	}

	var result DiscogsSearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("SearchRelease: failed to unmarshal response: %w", err)
	}

	// Cache the result
	if c.cacheStore != nil {
		if err := c.cacheStore.Set(ctx, cacheKey, body, searchCacheTTL); err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to store Discogs cache entry")
		}
	}

	return &result, nil
}

// GetRelease fetches a specific release by ID
func (c *DiscogsClient) GetRelease(ctx context.Context, releaseID int) (*DiscogsRelease, error) {
	l := logger.FromContext(ctx)

	cacheKey := discogsCacheKey("release", releaseID)
	if c.cacheStore != nil {
		body, found, err := c.cacheStore.Get(ctx, cacheKey)
		if err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to read Discogs cache entry")
		} else if found {
			var release DiscogsRelease
			if err := json.Unmarshal(body, &release); err == nil {
				return &release, nil
			}
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to unmarshal Discogs cache entry")
		}
	}

	releaseURL := fmt.Sprintf("%s/releases/%d", discogsBaseUrl, releaseID)

	req, err := http.NewRequest("GET", releaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("GetRelease: failed to build request: %w", err)
	}

	c.setAuthHeaders(req)

	body, err := c.queue(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetRelease: request failed: %w", err)
	}

	var release DiscogsRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("GetRelease: failed to unmarshal response: %w", err)
	}

	// Cache the result
	if c.cacheStore != nil {
		if err := c.cacheStore.Set(ctx, cacheKey, body, releaseCacheTTL); err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to store Discogs cache entry")
		}
	}

	return &release, nil
}

// GetReleaseGenres returns genres and styles for a release
func (c *DiscogsClient) GetReleaseGenres(ctx context.Context, releaseID int) ([]string, error) {
	release, err := c.GetRelease(ctx, releaseID)
	if err != nil {
		return nil, err
	}

	// Combine genres and styles
	genres := make([]string, 0, len(release.Genres)+len(release.Styles))
	genres = append(genres, release.Genres...)
	genres = append(genres, release.Styles...)

	return genres, nil
}

func (c *DiscogsClient) setAuthHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", fmt.Sprintf("Discogs key=%s, secret=%s", c.consumerKey, c.consumerSecret))
}

func (c *DiscogsClient) queue(ctx context.Context, req *http.Request) ([]byte, error) {
	l := logger.FromContext(ctx)
	req.Header.Set("Accept", "application/json")

	resultChan := c.requestQueue.Enqueue(func(client *http.Client, done chan<- queue.RequestResult) {
		resp, err := client.Do(req)
		if err != nil {
			l.Err(err).Str("url", req.URL.String()).Msg("Failed to contact Discogs")
			done <- queue.RequestResult{Err: err}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 || resp.StatusCode < 200 {
			err = fmt.Errorf("received non-ok status from Discogs: %s", resp.Status)
			done <- queue.RequestResult{Body: nil, Err: err}
			return
		}

		body, err := io.ReadAll(resp.Body)
		done <- queue.RequestResult{Body: body, Err: err}
	})

	result := <-resultChan
	return result.Body, result.Err
}
