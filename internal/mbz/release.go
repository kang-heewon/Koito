package mbz

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
)

type MusicBrainzGenre struct {
	Name string `json:"name"`
}

type MusicBrainzReleaseGroup struct {
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
