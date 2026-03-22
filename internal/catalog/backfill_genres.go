package catalog

import (
	"context"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/discogs"
	"github.com/gabehf/koito/internal/lastfm"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/mbz"
	"github.com/gabehf/koito/internal/models"
)

// GenreFetcher interface for hybrid genre fetching
type GenreFetcher interface {
	FetchAlbumGenres(ctx context.Context, store db.DB, album *models.Album) ([]string, error)
	FetchArtistGenres(ctx context.Context, store db.DB, artist *models.Artist) ([]string, error)
}

// HybridGenreFetcher implements hybrid genre fetching from multiple sources
type HybridGenreFetcher struct {
	mbz     mbz.MusicBrainzCaller
	discogs DiscogsCaller
	lastfm  LastFmCaller
	spotify SpotifyCaller
}

// DiscogsCaller interface for Discogs API
type DiscogsCaller interface {
	SearchRelease(ctx context.Context, artist, title string) (*discogs.DiscogsSearchResult, error)
	GetReleaseGenres(ctx context.Context, releaseID int) ([]string, error)
}

// LastFmCaller interface for Last.fm API
type LastFmCaller interface {
	GetAlbumTopTags(ctx context.Context, artist, album string) ([]lastfm.LastFmTag, error)
	GetArtistTopTags(ctx context.Context, artist string) ([]lastfm.LastFmTag, error)
}

// SpotifyCaller interface for Spotify API
type SpotifyCaller interface {
	GetArtistGenres(ctx context.Context, artistName string) ([]string, error)
}

// NewHybridGenreFetcher creates a new hybrid genre fetcher
func NewHybridGenreFetcher(mbzc mbz.MusicBrainzCaller, discogsC DiscogsCaller, lastfmC LastFmCaller, spotifyC SpotifyCaller) *HybridGenreFetcher {
	return &HybridGenreFetcher{
		mbz:     mbzc,
		discogs: discogsC,
		lastfm:  lastfmC,
		spotify: spotifyC,
	}
}

// FetchAlbumGenres tries MusicBrainz → Discogs → Last.fm
func (h *HybridGenreFetcher) FetchAlbumGenres(ctx context.Context, store db.DB, album *models.Album) ([]string, error) {
	l := logger.FromContext(ctx)

	// 1. Try MusicBrainz (album-level genres)
	if h.mbz != nil && album.MbzID != nil {
		release, err := h.mbz.GetReleaseWithGenres(ctx, *album.MbzID)
		if err == nil && release.ReleaseGroup != nil {
			genres := mbz.ReleaseGroupToGenres(release.ReleaseGroup)
			if len(genres) > 0 {
				l.Info().Str("source", "musicbrainz").Msgf("FetchAlbumGenres: Found %d genres for album %d (%s)", len(genres), album.ID, album.Title)
				return genres, nil
			}
			l.Debug().Str("source", "musicbrainz").Msgf("FetchAlbumGenres: No genres found in release group for album %d", album.ID)
		} else {
			l.Debug().Str("source", "musicbrainz").Err(err).Msgf("FetchAlbumGenres: Failed to fetch release for album %d", album.ID)
		}
	} else if h.mbz == nil {
		l.Warn().Msg("FetchAlbumGenres: MusicBrainz client not configured, skipping")
	}

	// 2. Try Discogs (release-level genres)
	if h.discogs != nil {
		artistName := h.getAlbumArtistName(ctx, store, album)
		if artistName != "" {
			result, err := h.discogs.SearchRelease(ctx, artistName, album.Title)
			if err == nil && result != nil && len(result.Results) > 0 {
				genres, err := h.discogs.GetReleaseGenres(ctx, result.Results[0].ID)
				if err == nil && len(genres) > 0 {
					l.Info().Str("source", "discogs").Msgf("FetchAlbumGenres: Found %d genres for album %d (%s)", len(genres), album.ID, album.Title)
					return genres, nil
				}
				l.Debug().Str("source", "discogs").Err(err).Msgf("FetchAlbumGenres: No genres found for album %d", album.ID)
			} else {
				l.Debug().Str("source", "discogs").Err(err).Msgf("FetchAlbumGenres: Failed to search release for album %d", album.ID)
			}
		}
	} else {
		l.Warn().Msg("FetchAlbumGenres: Discogs client not configured, skipping fallback")
	}

	// 3. Try Last.fm (album tags)
	if h.lastfm != nil {
		artistName := h.getAlbumArtistName(ctx, store, album)
		if artistName != "" {
			tags, err := h.lastfm.GetAlbumTopTags(ctx, artistName, album.Title)
			if err == nil && len(tags) > 0 {
				genres := h.tagsToGenres(tags)
				if len(genres) > 0 {
					l.Info().Str("source", "lastfm").Msgf("FetchAlbumGenres: Found %d genres for album %d (%s)", len(genres), album.ID, album.Title)
					return genres, nil
				}
				l.Debug().Str("source", "lastfm").Msgf("FetchAlbumGenres: No tags found for album %d", album.ID)
			} else {
				l.Debug().Str("source", "lastfm").Err(err).Msgf("FetchAlbumGenres: Failed to fetch tags for album %d", album.ID)
			}
		}
	} else {
		l.Warn().Msg("FetchAlbumGenres: Last.fm client not configured, skipping fallback")
	}

	return nil, nil
}

// FetchArtistGenres tries MusicBrainz → Last.fm → Spotify
func (h *HybridGenreFetcher) FetchArtistGenres(ctx context.Context, store db.DB, artist *models.Artist) ([]string, error) {
	l := logger.FromContext(ctx)

	// 1. Try MusicBrainz
	if h.mbz != nil && artist.MbzID != nil {
		genres, err := h.mbz.GetArtistGenres(ctx, *artist.MbzID)
		if err == nil && len(genres) > 0 {
			l.Info().Str("source", "musicbrainz").Msgf("FetchArtistGenres: Found %d genres for artist %d (%s)", len(genres), artist.ID, artist.Name)
			return genres, nil
		}
		l.Debug().Str("source", "musicbrainz").Err(err).Msgf("FetchArtistGenres: Failed to fetch genres for artist %d", artist.ID)
	} else if h.mbz == nil {
		l.Warn().Msg("FetchArtistGenres: MusicBrainz client not configured, skipping")
	}

	// 2. Try Last.fm
	if h.lastfm != nil {
		artistName := artist.Name
		if artistName != "" {
			tags, err := h.lastfm.GetArtistTopTags(ctx, artistName)
			if err == nil && len(tags) > 0 {
				genres := h.tagsToGenres(tags)
				if len(genres) > 0 {
					l.Info().Str("source", "lastfm").Msgf("FetchArtistGenres: Found %d genres for artist %d (%s)", len(genres), artist.ID, artist.Name)
					return genres, nil
				}
				l.Debug().Str("source", "lastfm").Msgf("FetchArtistGenres: No tags found for artist %d", artist.ID)
			} else {
				l.Debug().Str("source", "lastfm").Err(err).Msgf("FetchArtistGenres: Failed to fetch tags for artist %d", artist.ID)
			}
		}
	} else {
		l.Warn().Msg("FetchArtistGenres: Last.fm client not configured, skipping fallback")
	}

	// 3. Try Spotify (artist-level genres only)
	if h.spotify != nil {
		artistName := artist.Name
		if artistName != "" {
			genres, err := h.spotify.GetArtistGenres(ctx, artistName)
			if err == nil && len(genres) > 0 {
				l.Info().Str("source", "spotify").Msgf("FetchArtistGenres: Found %d genres for artist %d (%s)", len(genres), artist.ID, artist.Name)
				return genres, nil
			}
			l.Debug().Str("source", "spotify").Err(err).Msgf("FetchArtistGenres: No genres found for artist %d", artist.ID)
		}
	} else {
		l.Warn().Msg("FetchArtistGenres: Spotify client not configured, skipping fallback")
	}

	return nil, nil
}

func (h *HybridGenreFetcher) getAlbumArtistName(ctx context.Context, store db.DB, album *models.Album) string {
	// Get from album's Artists field if available
	if len(album.Artists) > 0 && album.Artists[0].Name != "" {
		return album.Artists[0].Name
	}
	// Fallback: get from album's artist relation
	artists, err := store.GetArtistsForAlbum(ctx, album.ID)
	if err != nil || len(artists) == 0 {
		return ""
	}
	return artists[0].Name
}

func (h *HybridGenreFetcher) tagsToGenres(tags []lastfm.LastFmTag) []string {
	genres := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag.Name != "" {
			genres = append(genres, tag.Name)
		}
	}
	return genres
}

// BackfillAlbumGenres backfills genres for albums without genres
func BackfillAlbumGenres(ctx context.Context, store db.DB, fetcher *HybridGenreFetcher) (int, error) {
	l := logger.FromContext(ctx)
	l.Info().Msg("BackfillAlbumGenres: Starting album genre backfill")

	var lastID int32 = 0
	totalProcessed := 0

	for {
		select {
		case <-ctx.Done():
			return totalProcessed, ctx.Err()
		default:
		}

		albumIDs, err := store.AlbumsWithoutGenres(ctx, lastID)
		if err != nil {
			l.Err(err).Msg("BackfillAlbumGenres: Failed to get albums without genres")
			return totalProcessed, err
		}

		if len(albumIDs) == 0 {
			break
		}

		for _, albumItem := range albumIDs {
			lastID = albumItem.ID

			// Get full album details from DB
			album, err := store.GetAlbum(ctx, db.GetAlbumOpts{ID: albumItem.ID})
			if err != nil {
				l.Debug().Err(err).Msgf("BackfillAlbumGenres: Failed to get album %d", albumItem.ID)
				continue
			}

			genres, err := fetcher.FetchAlbumGenres(ctx, store, album)
			if err != nil {
				l.Debug().Err(err).Msgf("BackfillAlbumGenres: Failed to fetch genres for album %d", album.ID)
				continue
			}

			if len(genres) == 0 {
				continue
			}

			err = store.SaveAlbumGenres(ctx, album.ID, genres)
			if err != nil {
				l.Warn().Err(err).Msgf("BackfillAlbumGenres: Failed to save genres for album %d", album.ID)
				continue
			}

			l.Debug().Msgf("BackfillAlbumGenres: Saved %d genres for album %d", len(genres), album.ID)
			totalProcessed++
		}
	}

	l.Info().Msgf("BackfillAlbumGenres: Completed. Updated %d albums with genres", totalProcessed)
	return totalProcessed, nil
}

// BackfillArtistGenres backfills genres for artists without genres
func BackfillArtistGenres(ctx context.Context, store db.DB, fetcher *HybridGenreFetcher) (int, error) {
	l := logger.FromContext(ctx)
	l.Info().Msg("BackfillArtistGenres: Starting artist genre backfill")

	var lastID int32 = 0
	totalProcessed := 0

	for {
		select {
		case <-ctx.Done():
			return totalProcessed, ctx.Err()
		default:
		}

		artistIDs, err := store.ArtistsWithoutGenres(ctx, lastID)
		if err != nil {
			l.Err(err).Msg("BackfillArtistGenres: Failed to get artists without genres")
			return totalProcessed, err
		}

		if len(artistIDs) == 0 {
			break
		}

		for _, artistItem := range artistIDs {
			lastID = artistItem.ID

			// Get full artist details from DB
			artist, err := store.GetArtist(ctx, db.GetArtistOpts{ID: artistItem.ID})
			if err != nil {
				l.Debug().Err(err).Msgf("BackfillArtistGenres: Failed to get artist %d", artistItem.ID)
				continue
			}

			genres, err := fetcher.FetchArtistGenres(ctx, store, artist)
			if err != nil {
				l.Debug().Err(err).Msgf("BackfillArtistGenres: Failed to fetch genres for artist %d", artist.ID)
				continue
			}

			if len(genres) == 0 {
				continue
			}

			err = store.SaveArtistGenres(ctx, artist.ID, genres)
			if err != nil {
				l.Warn().Err(err).Msgf("BackfillArtistGenres: Failed to save genres for artist %d", artist.ID)
				continue
			}

			l.Debug().Msgf("BackfillArtistGenres: Saved %d genres for artist %d", len(genres), artist.ID)
			totalProcessed++
		}
	}

	l.Info().Msgf("BackfillArtistGenres: Completed. Updated %d artists with genres", totalProcessed)
	return totalProcessed, nil
}

// BackfillGenres runs both album and artist genre backfill
func BackfillGenres(ctx context.Context, store db.DB, fetcher *HybridGenreFetcher) {
	l := logger.FromContext(ctx)
	l.Info().Msg("BackfillGenres: Starting genre backfill")

	var albumsUpdated, artistsUpdated int

	albumCount, err := BackfillAlbumGenres(ctx, store, fetcher)
	if err != nil {
		l.Err(err).Msg("BackfillGenres: Album genre backfill failed")
	} else {
		albumsUpdated = albumCount
	}

	artistCount, err := BackfillArtistGenres(ctx, store, fetcher)
	if err != nil {
		l.Err(err).Msg("BackfillGenres: Artist genre backfill failed")
	} else {
		artistsUpdated = artistCount
	}

	totalGenres := countTotalGenres(ctx, store)
	l.Info().Msgf("BackfillGenres: Completed. Backfilled %d albums, %d artists, total %d genres", albumsUpdated, artistsUpdated, totalGenres)
}

func countTotalGenres(ctx context.Context, store db.DB) int {
	stats, err := store.GetGenreStatsByListenCount(ctx, db.PeriodToTimeframe(db.PeriodAllTime))
	if err != nil {
		return 0
	}
	return len(stats)
}
