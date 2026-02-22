// package mbz provides functions for interacting with the musicbrainz api
package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gabehf/koito/internal/cache"
	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/queue"
	"github.com/google/uuid"
)

type MusicBrainzArea struct {
	Name           string   `json:"name"`
	Iso3166_1Codes []string `json:"iso-3166-1-codes"`
}

type MusicBrainzClient struct {
	url          string
	userAgent    string
	requestQueue *queue.RequestQueue
	cacheStore   cache.Store
}

const (
	mbzCachePrefix            = "mbz"
	artistCacheTTL            = 6 * time.Hour
	releaseGroupCacheTTL      = 6 * time.Hour
	releaseCacheTTL           = 24 * time.Hour
	releaseWithGenresCacheTTL = 6 * time.Hour
	recordingCacheTTL         = 24 * time.Hour
)

type MusicBrainzCaller interface {
	GetArtistPrimaryAliases(ctx context.Context, id uuid.UUID) ([]string, error)
	GetArtistGenres(ctx context.Context, id uuid.UUID) ([]string, error)
	GetReleaseTitles(ctx context.Context, RGID uuid.UUID) ([]string, error)
	GetTrack(ctx context.Context, id uuid.UUID) (*MusicBrainzTrack, error)
	GetReleaseGroup(ctx context.Context, id uuid.UUID) (*MusicBrainzReleaseGroup, error)
	GetRelease(ctx context.Context, id uuid.UUID) (*MusicBrainzRelease, error)
	GetReleaseWithGenres(ctx context.Context, id uuid.UUID) (*MusicBrainzRelease, error)
	Shutdown()
}

func NewMusicBrainzClient() *MusicBrainzClient {
	ret := new(MusicBrainzClient)
	ret.url = cfg.MusicBrainzUrl()
	ret.userAgent = cfg.UserAgent()
	ret.requestQueue = queue.NewRequestQueue(cfg.MusicBrainzRateLimit(), cfg.MusicBrainzRateLimit())
	ret.cacheStore = cache.NewDefaultStore()
	return ret
}

func newMusicBrainzClientWithCache(url string, cacheStore cache.Store) *MusicBrainzClient {
	ret := new(MusicBrainzClient)
	ret.url = url
	ret.userAgent = "koito-test"
	ret.requestQueue = queue.NewRequestQueue(100, 100)
	ret.cacheStore = cacheStore
	return ret
}

func (c *MusicBrainzClient) Shutdown() {
	c.requestQueue.Shutdown()
}

func mbzCacheKey(entity string, id uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", mbzCachePrefix, entity, id.String())
}

func (c *MusicBrainzClient) getEntity(ctx context.Context, fmtStr string, id uuid.UUID, result any) error {
	return c.getEntityCached(ctx, "", 0, fmtStr, id, result)
}

func (c *MusicBrainzClient) getEntityCached(ctx context.Context, cacheKey string, ttl time.Duration, fmtStr string, id uuid.UUID, result any) error {
	l := logger.FromContext(ctx)

	if c.cacheStore != nil && cacheKey != "" && ttl > 0 {
		body, found, err := c.cacheStore.Get(ctx, cacheKey)
		if err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to read MusicBrainz cache entry")
		} else if found {
			err = json.Unmarshal(body, result)
			if err == nil {
				return nil
			}
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to unmarshal MusicBrainz cache entry")
		}
	}

	body, err := c.getEntityBody(ctx, fmtStr, id)
	if err != nil {
		return fmt.Errorf("getEntityCached: %w", err)
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		l.Err(err).Str("body", string(body)).Msg("Failed to unmarshal MusicBrainz response body")
		return fmt.Errorf("getEntityCached: %w", err)
	}

	if c.cacheStore != nil && cacheKey != "" && ttl > 0 {
		err := c.cacheStore.Set(ctx, cacheKey, body, ttl)
		if err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to store MusicBrainz cache entry")
		}
	}

	return nil
}

func (c *MusicBrainzClient) getEntityBody(ctx context.Context, fmtStr string, id uuid.UUID) ([]byte, error) {
	l := logger.FromContext(ctx)
	url := fmt.Sprintf(fmtStr, c.url, id.String())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		l.Err(err).Msg("Failed to build MusicBrainz request")
		return nil, fmt.Errorf("getEntityBody: %w", err)
	}
	l.Debug().Msg("Adding MusicBrainz request to queue")
	body, err := c.queue(ctx, req)
	if err != nil {
		l.Err(err).Msg("MusicBrainz request failed")
		return nil, fmt.Errorf("getEntityBody: %w", err)
	}
	return body, nil
}

func (c *MusicBrainzClient) queue(ctx context.Context, req *http.Request) ([]byte, error) {
	l := logger.FromContext(ctx)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resultChan := c.requestQueue.Enqueue(func(client *http.Client, done chan<- queue.RequestResult) {
		resp, err := client.Do(req)
		if err != nil {
			l.Err(err).Str("url", req.RequestURI).Msg("Failed to contact MusicBrainz")
			done <- queue.RequestResult{Err: err}
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 || resp.StatusCode < 200 {
			err = fmt.Errorf("recieved non-ok status from MusicBrainz: %s", resp.Status)
			done <- queue.RequestResult{Body: nil, Err: err}
			return
		}

		body, err := io.ReadAll(resp.Body)
		done <- queue.RequestResult{Body: body, Err: err}
	})

	result := <-resultChan
	return result.Body, result.Err
}
