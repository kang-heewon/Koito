package catalog

import (
	"context"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/mbz"
	"github.com/google/uuid"
)

func BackfillMbzMatching(ctx context.Context, store db.DB, mbzc mbz.MusicBrainzCaller) error {
	l := logger.FromContext(ctx)
	l.Info().Msg("BackfillMbzMatching: Starting MBZ ID matching for albums")

	var lastID int32 = 0
	totalProcessed := 0
	totalMatched := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		albums, err := store.AlbumsWithoutMbzID(ctx, lastID)
		if err != nil {
			l.Err(err).Msg("BackfillMbzMatching: Failed to get albums without MBZ ID")
			return err
		}

		if len(albums) == 0 {
			break
		}

		for _, album := range albums {
			lastID = album.ID

			// Get artist name from first artist
			if len(album.Artists) == 0 {
				l.Debug().Msgf("BackfillMbzMatching: Skipping album %d - no artists", album.ID)
				continue
			}

			artistName := album.Artists[0].Name
			if artistName == "" {
				l.Debug().Msgf("BackfillMbzMatching: Skipping album %d - empty artist name", album.ID)
				continue
			}

			// Search for release
			result, err := mbzc.SearchRelease(ctx, artistName, album.Title)
			if err != nil {
				l.Warn().Err(err).Msgf("BackfillMbzMatching: Error searching for album %d", album.ID)
				continue
			}

			if result == nil || len(result.Releases) == 0 {
				l.Debug().Msgf("BackfillMbzMatching: No search results for album %d (%s - %s)", album.ID, artistName, album.Title)
				store.MarkMbzSearched(ctx, album.ID)
				continue
			}

			// Find release with score == 100
			var matchedRelease *mbz.MusicBrainzSearchRelease
			for _, release := range result.Releases {
				if release.Score == 100 {
					matchedRelease = &release
					break
				}
			}

			if matchedRelease == nil {
				l.Debug().Msgf("BackfillMbzMatching: No score 100 match for album %d (%s - %s)", album.ID, artistName, album.Title)
				store.MarkMbzSearched(ctx, album.ID)
				continue
			}

			// Parse release UUID (use Parse, not MustParse to avoid panic)
			releaseID, err := uuid.Parse(matchedRelease.ID)
			if err != nil {
				l.Warn().Err(err).Msgf("BackfillMbzMatching: Invalid release ID %s for album %d", matchedRelease.ID, album.ID)
				store.MarkMbzSearched(ctx, album.ID)
				continue
			}

			// Get release to retrieve release group
			release, err := mbzc.GetRelease(ctx, releaseID)
			if err != nil {
				l.Warn().Err(err).Msgf("BackfillMbzMatching: Failed to get release %s for album %d", releaseID, album.ID)
				store.MarkMbzSearched(ctx, album.ID)
				continue
			}

			if release.ReleaseGroup == nil {
				l.Warn().Msgf("BackfillMbzMatching: No release group for release %s", releaseID)
				store.MarkMbzSearched(ctx, album.ID)
				continue
			}

			// Update album with MBZ ID
			err = store.UpdateAlbum(ctx, db.UpdateAlbumOpts{
				ID:            album.ID,
				MusicBrainzID: releaseID,
			})
			if err != nil {
				l.Warn().Err(err).Msgf("BackfillMbzMatching: Failed to update album %d with MBZ ID %s", album.ID, releaseID)
				continue
			}

		// Get genres from release group
		releaseGroupID, err := uuid.Parse(release.ReleaseGroup.ID)
		if err != nil {
			l.Warn().Err(err).Msgf("BackfillMbzMatching: Invalid release group ID %s for album %d", release.ReleaseGroup.ID, album.ID)
			store.MarkMbzSearched(ctx, album.ID)
			continue
		}
		rg, err := mbzc.GetReleaseGroup(ctx, releaseGroupID)
		if err != nil {
			l.Warn().Err(err).Msgf("BackfillMbzMatching: Failed to get release group %s for album %d", releaseGroupID, album.ID)
			store.MarkMbzSearched(ctx, album.ID)
			continue
		}

			genres := mbz.ReleaseGroupToGenres(rg)
			if len(genres) > 0 {
				err = store.SaveAlbumGenres(ctx, album.ID, genres)
				if err != nil {
					l.Warn().Err(err).Msgf("BackfillMbzMatching: Failed to save genres for album %d", album.ID)
				}
			}

			// Mark as searched
			store.MarkMbzSearched(ctx, album.ID)

			totalMatched++
			totalProcessed++

			if totalProcessed%10 == 0 {
				l.Info().Msgf("BackfillMbzMatching: Matched album %s by %s → MBZ release %s", album.Title, artistName, matchedRelease.ID)
			}
		}
	}

	l.Info().Msgf("BackfillMbzMatching: Completed. Matched %d out of %d albums processed", totalMatched, totalProcessed)
	return nil
}
