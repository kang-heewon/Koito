// package lastfm provides functions for interacting with the Last.fm API
package lastfm

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
	lastfmBaseUrl     = "https://ws.audioscrobbler.com/2.0/"
	lastfmCachePrefix = "lastfm"
	artistCacheTTL    = 6 * time.Hour
	albumCacheTTL     = 24 * time.Hour
	tagCacheTTL       = 6 * time.Hour
)

// LastFmClient handles communication with Last.fm API
type LastFmClient struct {
	apiKey       string
	userAgent    string
	requestQueue *queue.RequestQueue
	cacheStore   cache.Store
}

// LastFmArtist represents an artist from Last.fm
type LastFmArtist struct {
	Name  string        `json:"name"`
	MBID  string        `json:"mbid"`
	URL   string        `json:"url"`
	Stats LastFmStats   `json:"stats"`
	Tags  LastFmTagList `json:"tags"`
}

// LastFmStats represents listener/play stats
type LastFmStats struct {
	Listeners string `json:"listeners"`
	Playcount string `json:"playcount"`
}

// LastFmTagList represents a list of tags
type LastFmTagList struct {
	Tag []LastFmTag `json:"tag"`
}

// LastFmTag represents a tag from Last.fm
type LastFmTag struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Count int    `json:"count,omitempty"`
}

// LastFmAlbum represents an album from Last.fm
type LastFmAlbum struct {
	Artist string        `json:"artist"`
	Title  string        `json:"title"`
	MBID   string        `json:"mbid"`
	URL    string        `json:"url"`
	Tags   LastFmTagList `json:"tags"`
}

// LastFmTopTags represents top tags response
type LastFmTopTags struct {
	Toptags struct {
		Tag  []LastFmTag `json:"tag"`
		Attr struct {
			Artist string `json:"artist"`
		} `json:"@attr"`
	} `json:"toptags"`
}

// LastFmArtistInfo represents artist info response
type LastFmArtistInfo struct {
	Artist LastFmArtist `json:"artist"`
}

// LastFmAlbumInfo represents album info response
type LastFmAlbumInfo struct {
	Album LastFmAlbum `json:"album"`
}

// LastFmCaller interface for Last.fm API operations
type LastFmCaller interface {
	GetArtistTopTags(ctx context.Context, artist string) ([]LastFmTag, error)
	GetArtistInfo(ctx context.Context, artist string) (*LastFmArtist, error)
	GetAlbumInfo(ctx context.Context, artist, album string) (*LastFmAlbum, error)
	GetAlbumTopTags(ctx context.Context, artist, album string) ([]LastFmTag, error)
	Shutdown()
}

// NewLastFmClient creates a new Last.fm client
func NewLastFmClient() *LastFmClient {
	ret := new(LastFmClient)
	ret.apiKey = cfg.LastFmApiKey()
	ret.userAgent = cfg.UserAgent()
	ret.requestQueue = queue.NewRequestQueue(1, 1) // Last.fm rate limit: ~1 req/sec
	ret.cacheStore = cache.NewDefaultStore()
	return ret
}

// NewLastFmClientWithCache creates a Last.fm client with custom cache (for testing)
func NewLastFmClientWithCache(cacheStore cache.Store) *LastFmClient {
	ret := new(LastFmClient)
	ret.apiKey = "test-api-key"
	ret.userAgent = "koito-test"
	ret.requestQueue = queue.NewRequestQueue(100, 100)
	ret.cacheStore = cacheStore
	return ret
}

func (c *LastFmClient) Shutdown() {
	c.requestQueue.Shutdown()
}

func lastfmCacheKey(entity string, id string) string {
	return fmt.Sprintf("%s:%s:%s", lastfmCachePrefix, entity, url.QueryEscape(id))
}

// GetArtistTopTags fetches top tags for an artist
func (c *LastFmClient) GetArtistTopTags(ctx context.Context, artist string) ([]LastFmTag, error) {
	l := logger.FromContext(ctx)

	cacheKey := lastfmCacheKey("artist_tags", artist)
	if c.cacheStore != nil {
		body, found, err := c.cacheStore.Get(ctx, cacheKey)
		if err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to read Last.fm cache entry")
		} else if found {
			var result LastFmTopTags
			if err := json.Unmarshal(body, &result); err == nil {
				return result.Toptags.Tag, nil
			}
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to unmarshal Last.fm cache entry")
		}
	}

	params := url.Values{}
	params.Set("method", "artist.gettoptags")
	params.Set("artist", artist)
	params.Set("api_key", c.apiKey)
	params.Set("format", "json")

	req, err := http.NewRequest("GET", lastfmBaseUrl+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("GetArtistTopTags: failed to build request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	body, err := c.queue(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetArtistTopTags: request failed: %w", err)
	}

	var result LastFmTopTags
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("GetArtistTopTags: failed to unmarshal response: %w", err)
	}

	// Cache the result
	if c.cacheStore != nil {
		if err := c.cacheStore.Set(ctx, cacheKey, body, tagCacheTTL); err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to store Last.fm cache entry")
		}
	}

	return result.Toptags.Tag, nil
}

// GetArtistInfo fetches artist information
func (c *LastFmClient) GetArtistInfo(ctx context.Context, artist string) (*LastFmArtist, error) {
	l := logger.FromContext(ctx)

	cacheKey := lastfmCacheKey("artist", artist)
	if c.cacheStore != nil {
		body, found, err := c.cacheStore.Get(ctx, cacheKey)
		if err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to read Last.fm cache entry")
		} else if found {
			var result LastFmArtistInfo
			if err := json.Unmarshal(body, &result); err == nil {
				return &result.Artist, nil
			}
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to unmarshal Last.fm cache entry")
		}
	}

	params := url.Values{}
	params.Set("method", "artist.getinfo")
	params.Set("artist", artist)
	params.Set("api_key", c.apiKey)
	params.Set("format", "json")

	req, err := http.NewRequest("GET", lastfmBaseUrl+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("GetArtistInfo: failed to build request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	body, err := c.queue(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetArtistInfo: request failed: %w", err)
	}

	var result LastFmArtistInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("GetArtistInfo: failed to unmarshal response: %w", err)
	}

	// Cache the result
	if c.cacheStore != nil {
		if err := c.cacheStore.Set(ctx, cacheKey, body, artistCacheTTL); err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to store Last.fm cache entry")
		}
	}

	return &result.Artist, nil
}

// GetAlbumInfo fetches album information
func (c *LastFmClient) GetAlbumInfo(ctx context.Context, artist, album string) (*LastFmAlbum, error) {
	l := logger.FromContext(ctx)

	cacheKey := lastfmCacheKey("album", artist+"-"+album)
	if c.cacheStore != nil {
		body, found, err := c.cacheStore.Get(ctx, cacheKey)
		if err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to read Last.fm cache entry")
		} else if found {
			var result LastFmAlbumInfo
			if err := json.Unmarshal(body, &result); err == nil {
				return &result.Album, nil
			}
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to unmarshal Last.fm cache entry")
		}
	}

	params := url.Values{}
	params.Set("method", "album.getinfo")
	params.Set("artist", artist)
	params.Set("album", album)
	params.Set("api_key", c.apiKey)
	params.Set("format", "json")

	req, err := http.NewRequest("GET", lastfmBaseUrl+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("GetAlbumInfo: failed to build request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	body, err := c.queue(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetAlbumInfo: request failed: %w", err)
	}

	var result LastFmAlbumInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("GetAlbumInfo: failed to unmarshal response: %w", err)
	}

	// Cache the result
	if c.cacheStore != nil {
		if err := c.cacheStore.Set(ctx, cacheKey, body, albumCacheTTL); err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to store Last.fm cache entry")
		}
	}

	return &result.Album, nil
}

// GetAlbumTopTags fetches top tags for an album
func (c *LastFmClient) GetAlbumTopTags(ctx context.Context, artist, album string) ([]LastFmTag, error) {
	l := logger.FromContext(ctx)

	cacheKey := lastfmCacheKey("album_tags", artist+"-"+album)
	if c.cacheStore != nil {
		body, found, err := c.cacheStore.Get(ctx, cacheKey)
		if err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to read Last.fm cache entry")
		} else if found {
			var result struct {
				Toptags struct {
					Tag []LastFmTag `json:"tag"`
				} `json:"toptags"`
			}
			if err := json.Unmarshal(body, &result); err == nil {
				return result.Toptags.Tag, nil
			}
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to unmarshal Last.fm cache entry")
		}
	}

	params := url.Values{}
	params.Set("method", "album.gettoptags")
	params.Set("artist", artist)
	params.Set("album", album)
	params.Set("api_key", c.apiKey)
	params.Set("format", "json")

	req, err := http.NewRequest("GET", lastfmBaseUrl+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("GetAlbumTopTags: failed to build request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	body, err := c.queue(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetAlbumTopTags: request failed: %w", err)
	}

	var result struct {
		Toptags struct {
			Tag []LastFmTag `json:"tag"`
		} `json:"toptags"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("GetAlbumTopTags: failed to unmarshal response: %w", err)
	}

	// Cache the result
	if c.cacheStore != nil {
		if err := c.cacheStore.Set(ctx, cacheKey, body, tagCacheTTL); err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to store Last.fm cache entry")
		}
	}

	return result.Toptags.Tag, nil
}

func (c *LastFmClient) queue(ctx context.Context, req *http.Request) ([]byte, error) {
	l := logger.FromContext(ctx)

	resultChan := c.requestQueue.Enqueue(func(client *http.Client, done chan<- queue.RequestResult) {
		resp, err := client.Do(req)
		if err != nil {
			l.Err(err).Str("url", req.URL.String()).Msg("Failed to contact Last.fm")
			done <- queue.RequestResult{Err: err}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 || resp.StatusCode < 200 {
			err = fmt.Errorf("received non-ok status from Last.fm: %s", resp.Status)
			done <- queue.RequestResult{Body: nil, Err: err}
			return
		}

		body, err := io.ReadAll(resp.Body)
		done <- queue.RequestResult{Body: body, Err: err}
	})

	result := <-resultChan
	return result.Body, result.Err
}
