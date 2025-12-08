package catalog

import (
	"context"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/mbz"
)

func BackfillAlbumGenres(ctx context.Context, store db.DB, mbzc mbz.MusicBrainzCaller) error {
	l := logger.FromContext(ctx)
	l.Info().Msg("BackfillAlbumGenres: Starting album genre backfill")

	var lastID int32 = 0
	totalProcessed := 0

	for {
		albums, err := store.AlbumsWithoutGenres(ctx, lastID)
		if err != nil {
			l.Err(err).Msg("BackfillAlbumGenres: Failed to get albums without genres")
			return err
		}

		if len(albums) == 0 {
			break
		}

		for _, album := range albums {
			lastID = album.ID

			release, err := mbzc.GetReleaseWithGenres(ctx, album.MbzID)
			if err != nil {
				l.Debug().Err(err).Msgf("BackfillAlbumGenres: Failed to get release for album %d", album.ID)
				continue
			}

			if release.ReleaseGroup == nil {
				l.Debug().Msgf("BackfillAlbumGenres: No release group found for album %d", album.ID)
				continue
			}

			genres := mbz.ReleaseGroupToGenres(release.ReleaseGroup)
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
	return nil
}

func BackfillArtistGenres(ctx context.Context, store db.DB, mbzc mbz.MusicBrainzCaller) error {
	l := logger.FromContext(ctx)
	l.Info().Msg("BackfillArtistGenres: Starting artist genre backfill")

	var lastID int32 = 0
	totalProcessed := 0

	for {
		artists, err := store.ArtistsWithoutGenres(ctx, lastID)
		if err != nil {
			l.Err(err).Msg("BackfillArtistGenres: Failed to get artists without genres")
			return err
		}

		if len(artists) == 0 {
			break
		}

		for _, artist := range artists {
			lastID = artist.ID

			genres, err := mbzc.GetArtistGenres(ctx, artist.MbzID)
			if err != nil {
				l.Debug().Err(err).Msgf("BackfillArtistGenres: Failed to get genres for artist %d", artist.ID)
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
	return nil
}

func BackfillGenres(ctx context.Context, store db.DB, mbzc mbz.MusicBrainzCaller) {
	l := logger.FromContext(ctx)

	if err := BackfillAlbumGenres(ctx, store, mbzc); err != nil {
		l.Err(err).Msg("BackfillGenres: Album genre backfill failed")
	}

	if err := BackfillArtistGenres(ctx, store, mbzc); err != nil {
		l.Err(err).Msg("BackfillGenres: Artist genre backfill failed")
	}
}
