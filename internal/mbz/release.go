package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type MusicBrainzGenre struct {
	Name string `json:"name"`
}

type MusicBrainzReleaseGroup struct {
	ID           string                    `json:"id"`
	Title        string                    `json:"title"`
	Type         string                    `json:"primary_type"`
	ArtistCredit []MusicBrainzArtistCredit `json:"artist-credit"`
	Releases     []MusicBrainzRelease      `json:"releases"`
	Genres       []MusicBrainzGenre        `json:"genres"`
}

type MusicBrainzRelease struct {
	Title              string                    `json:"title"`
	ID                 string                    `json:"id"`
	ArtistCredit       []MusicBrainzArtistCredit `json:"artist-credit"`
	Status             string                    `json:"status"`
	TextRepresentation TextRepresentation        `json:"text-representation"`
	ReleaseGroup       *MusicBrainzReleaseGroup  `json:"release-group"`
}
type MusicBrainzArtistCredit struct {
	Artist MusicBrainzArtist `json:"artist"`
	Name   string            `json:"name"`
}
type TextRepresentation struct {
	Language string `json:"language"`
	Script   string `json:"script"`
}

const releaseGroupFmtStr = "%s/ws/2/release-group/%s?inc=releases+artists+genres"
const releaseFmtStr = "%s/ws/2/release/%s?inc=artists"
const releaseWithGenresFmtStr = "%s/ws/2/release/%s?inc=release-groups+genres"

func (c *MusicBrainzClient) GetReleaseGroup(ctx context.Context, id uuid.UUID) (*MusicBrainzReleaseGroup, error) {
	mbzRG := new(MusicBrainzReleaseGroup)
	err := c.getEntityCached(ctx, mbzCacheKey("release-group", id), releaseGroupCacheTTL, releaseGroupFmtStr, id, mbzRG)
	if err != nil {
		return nil, fmt.Errorf("GetReleaseGroup: %w", err)
	}
	return mbzRG, nil
}

func (c *MusicBrainzClient) GetRelease(ctx context.Context, id uuid.UUID) (*MusicBrainzRelease, error) {
	mbzRelease := new(MusicBrainzRelease)
	err := c.getEntityCached(ctx, mbzCacheKey("release", id), releaseCacheTTL, releaseFmtStr, id, mbzRelease)
	if err != nil {
		return nil, fmt.Errorf("GetRelease: %w", err)
	}
	return mbzRelease, nil
}

func (c *MusicBrainzClient) GetReleaseWithGenres(ctx context.Context, id uuid.UUID) (*MusicBrainzRelease, error) {
	mbzRelease := new(MusicBrainzRelease)
	err := c.getEntityCached(ctx, mbzCacheKey("release-with-genres", id), releaseWithGenresCacheTTL, releaseWithGenresFmtStr, id, mbzRelease)
	if err != nil {
		return nil, fmt.Errorf("GetReleaseWithGenres: %w", err)
	}
	return mbzRelease, nil
}

func (c *MusicBrainzClient) GetReleaseTitles(ctx context.Context, RGID uuid.UUID) ([]string, error) {
	releaseGroup, err := c.GetReleaseGroup(ctx, RGID)
	if err != nil {
		return nil, fmt.Errorf("GetReleaseTitles: %w", err)
	}

	var titles []string
	for _, release := range releaseGroup.Releases {
		if !slices.Contains(titles, release.Title) {
			titles = append(titles, release.Title)
		}
	}

	return titles, nil
}

func ReleaseGroupToTitles(rg *MusicBrainzReleaseGroup) []string {
	var titles []string
	for _, release := range rg.Releases {
		if !slices.Contains(titles, release.Title) {
			titles = append(titles, release.Title)
		}
	}
	return titles
}

func ReleaseGroupToGenres(rg *MusicBrainzReleaseGroup) []string {
	genres := make([]string, len(rg.Genres))
	for i, g := range rg.Genres {
		genres[i] = g.Name
	}
	return genres
}

// Searches for Pseudo-Releases of release groups with Latin script, and returns them as an array
func (c *MusicBrainzClient) GetLatinTitles(ctx context.Context, id uuid.UUID) ([]string, error) {
	rg, err := c.GetReleaseGroup(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetLatinTitles: %w", err)
	}
	titles := make([]string, 0)
	for _, r := range rg.Releases {
		if r.Status == "Pseudo-Release" && r.TextRepresentation.Script == "Latn" { // not a typo
			titles = append(titles, r.Title)
		}
	}
	return titles, nil
}

type MusicBrainzSearchResult struct {
	Releases []MusicBrainzSearchRelease `json:"releases"`
}

type MusicBrainzSearchRelease struct {
	ID           string                    `json:"id"`
	Score        int                       `json:"score"`
	Title        string                    `json:"title"`
	ArtistCredit []MusicBrainzArtistCredit `json:"artist-credit"`
	ReleaseGroup *MusicBrainzReleaseGroup  `json:"release-group"`
}

const searchReleaseFmtStr = "%s/ws/2/release/?query=release:\"%s\" AND artist:\"%s\"&limit=5&fmt=json"
const searchReleaseCacheTTL = 24 * time.Hour

func (c *MusicBrainzClient) SearchRelease(ctx context.Context, artist, title string) (*MusicBrainzSearchResult, error) {
	l := zerolog.Ctx(ctx)

	titleEscaped := url.QueryEscape(title)
	artistEscaped := url.QueryEscape(artist)

	cacheKey := fmt.Sprintf("mbz:search:release-artist:%s:%s", titleEscaped, artistEscaped)

	// Check cache first
	if c.cacheStore != nil {
		body, found, err := c.cacheStore.Get(ctx, cacheKey)
		if err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to read MusicBrainz search cache entry")
		} else if found {
			mbzResult := new(MusicBrainzSearchResult)
			err = json.Unmarshal(body, mbzResult)
			if err == nil {
				return mbzResult, nil
			}
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to unmarshal MusicBrainz search cache entry")
		}
	}

	// Build search URL
	searchURL := fmt.Sprintf(searchReleaseFmtStr, c.url, titleEscaped, artistEscaped)
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		l.Err(err).Msg("Failed to build MusicBrainz search request")
		return nil, nil
	}

	body, err := c.queue(ctx, req)
	if err != nil {
		l.Err(err).Str("artist", artist).Str("title", title).Msg("MusicBrainz search request failed")
		return nil, nil
	}

	mbzResult := new(MusicBrainzSearchResult)
	err = json.Unmarshal(body, mbzResult)
	if err != nil {
		l.Err(err).Str("body", string(body)).Msg("Failed to unmarshal MusicBrainz search response")
		return nil, nil
	}

	// Cache the result
	if c.cacheStore != nil {
		err := c.cacheStore.Set(ctx, cacheKey, body, searchReleaseCacheTTL)
		if err != nil {
			l.Warn().Err(err).Str("cache_key", cacheKey).Msg("Failed to store MusicBrainz search cache entry")
		}
	}

	return mbzResult, nil
}

