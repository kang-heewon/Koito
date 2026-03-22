package mbz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type MusicBrainzGenre struct {
	Name string `json:"name"`
}

type MusicBrainzTag struct {
	Name string `json:"name"`
}

type MusicBrainzReleaseGroup struct {
	ID           string                    `json:"id"`
	Title        string                    `json:"title"`
	Type         string                    `json:"primary_type"`
	ArtistCredit []MusicBrainzArtistCredit `json:"artist-credit"`
	Releases     []MusicBrainzRelease      `json:"releases"`
	Genres       []MusicBrainzGenre        `json:"genres"`
	Tags         []MusicBrainzTag          `json:"tags"`
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
const releaseGroupGenresFmtStr = "%s/ws/2/release-group/%s?inc=genres"
const releaseFmtStr = "%s/ws/2/release/%s?inc=artists"
const releaseWithGenresFmtStr = "%s/ws/2/release/%s?inc=release-groups+genres+tags"

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

func (c *MusicBrainzClient) GetReleaseGroupGenres(ctx context.Context, id uuid.UUID) ([]string, error) {
	mbzReleaseGroup := new(MusicBrainzReleaseGroup)
	err := c.getEntityCached(ctx, mbzCacheKey("release-group-genres", id), releaseGroupCacheTTL, releaseGroupGenresFmtStr, id, mbzReleaseGroup)
	if err != nil {
		return nil, fmt.Errorf("GetReleaseGroupGenres: %w", err)
	}
	return musicBrainzGenresToNames(mbzReleaseGroup.Genres), nil
}

func (c *MusicBrainzClient) GetReleaseWithGenres(ctx context.Context, id uuid.UUID) (*MusicBrainzRelease, error) {
	mbzRelease := new(MusicBrainzRelease)
	err := c.getEntityCached(ctx, mbzCacheKey("release-with-genres", id), releaseWithGenresCacheTTL, releaseWithGenresFmtStr, id, mbzRelease)
	if err != nil {
		return nil, fmt.Errorf("GetReleaseWithGenres: %w", err)
	}
	if mbzRelease.ReleaseGroup != nil && len(mbzRelease.ReleaseGroup.Genres) == 0 && mbzRelease.ReleaseGroup.ID != "" {
		releaseGroupID, err := uuid.Parse(mbzRelease.ReleaseGroup.ID)
		if err != nil {
			zerolog.Ctx(ctx).Warn().Err(err).Str("release_group_id", mbzRelease.ReleaseGroup.ID).Msg("GetReleaseWithGenres: invalid release group id")
			return mbzRelease, nil
		}
		genres, err := c.GetReleaseGroupGenres(ctx, releaseGroupID)
		if err != nil {
			zerolog.Ctx(ctx).Warn().Err(err).Str("release_group_id", mbzRelease.ReleaseGroup.ID).Msg("GetReleaseWithGenres: failed to fetch release group genres")
			return mbzRelease, nil
		}
		mbzRelease.ReleaseGroup.Genres = musicBrainzGenres(genres)
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
	if rg == nil {
		return nil
	}
	genres := musicBrainzGenresToNames(rg.Genres)
	if len(genres) > 0 {
		return genres
	}
	return musicBrainzTagsToGenres(rg.Tags)
}

func musicBrainzGenresToNames(genres []MusicBrainzGenre) []string {
	ret := make([]string, len(genres))
	for i, genre := range genres {
		ret[i] = genre.Name
	}
	return ret
}

func musicBrainzGenres(genres []string) []MusicBrainzGenre {
	ret := make([]MusicBrainzGenre, len(genres))
	for i, genre := range genres {
		ret[i] = MusicBrainzGenre{Name: genre}
	}
	return ret
}

func musicBrainzTagsToGenres(tags []MusicBrainzTag) []string {
	ret := make([]string, 0, len(tags))
	seen := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		name := strings.ToLower(strings.TrimSpace(tag.Name))
		if name == "" {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		ret = append(ret, name)
	}
	return ret
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
