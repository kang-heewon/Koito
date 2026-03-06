package catalog

import (
	"context"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/images"
	"github.com/gabehf/koito/internal/logger"
	"github.com/google/uuid"
)

func BackfillImages(ctx context.Context, store db.DB) {
	l := logger.FromContext(ctx)
	l.Info().Msg("BackfillImages: Starting image backfill")

	lastID := int32(0)
	processed := 0

	for {
		select {
		case <-ctx.Done():
			l.Info().Msg("BackfillImages: Context cancelled, stopping")
			return
		default:
		}

		albums, err := store.AlbumsWithoutImages(ctx, lastID)
		if err != nil {
			l.Err(err).Int32("last_id", lastID).Msg("BackfillImages: Failed to get albums without images")
			return
		}

		if len(albums) == 0 {
			break
		}

		for _, album := range albums {
			select {
			case <-ctx.Done():
				l.Info().Msg("BackfillImages: Context cancelled during processing")
				return
			default:
			}

			// Extract artist names
			var artistNames []string
			for _, artist := range album.Artists {
				artistNames = append(artistNames, artist.Name)
			}

			// Skip albums without artists
			if len(artistNames) == 0 {
				l.Debug().Int32("album_id", album.ID).Msg("BackfillImages: Skipping album without artists")
				continue
			}

			// Get album image
			imgURL, err := images.GetAlbumImage(ctx, images.AlbumImageOpts{
				Artists:      artistNames,
				Album:        album.Title,
				ReleaseMbzID: album.MbzID,
			})
			if err != nil {
				l.Warn().Err(err).Int32("album_id", album.ID).Msg("BackfillImages: Failed to get album image")
				continue
			}

			// Skip if no image found
			if imgURL == "" {
				l.Debug().Int32("album_id", album.ID).Msg("BackfillImages: No image found for album")
				continue
			}

			// Download and cache image
			imgID := uuid.New()
			err = DownloadAndCacheImage(ctx, imgID, imgURL, ImageSourceSize())
			if err != nil {
				l.Warn().Err(err).Int32("album_id", album.ID).Str("url", imgURL).Msg("BackfillImages: Failed to download and cache image")
				continue
			}

			// Update album with image
			err = store.UpdateAlbum(ctx, db.UpdateAlbumOpts{
				ID:       album.ID,
				Image:    imgID,
				ImageSrc: imgURL,
			})
			if err != nil {
				l.Warn().Err(err).Int32("album_id", album.ID).Msg("BackfillImages: Failed to update album with image")
				continue
			}

			processed++
			if processed%10 == 0 {
				l.Info().Int("processed", processed).Msg("BackfillImages: Progress")
			}
		}

		lastID = albums[len(albums)-1].ID
	}

	l.Info().Int("processed", processed).Msg("BackfillImages: Completed")
}
