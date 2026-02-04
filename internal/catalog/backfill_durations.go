package catalog

import (
	"context"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/mbz"
)

func BackfillTrackDurations(ctx context.Context, store db.DB, mbzc mbz.MusicBrainzCaller) error {
	l := logger.FromContext(ctx)
	l.Info().Msg("BackfillTrackDurations: Starting track duration backfill")

	var lastID int32 = 0
	totalProcessed := 0

	for {
		tracks, err := store.TracksWithoutDuration(ctx, lastID)
		if err != nil {
			l.Err(err).Msg("BackfillTrackDurations: Failed to get tracks without duration")
			return err
		}

		if len(tracks) == 0 {
			break
		}

		for _, track := range tracks {
			lastID = track.ID

			mbzTrack, err := mbzc.GetTrack(ctx, track.MbzID)
			if err != nil {
				l.Debug().Err(err).Msgf("BackfillTrackDurations: Failed to get track %d from MusicBrainz", track.ID)
				continue
			}

			if mbzTrack.LengthMs == 0 {
				continue
			}

			// Convert ms to seconds for storage
			durationSeconds := int32(mbzTrack.LengthMs / 1000)

			err = store.UpdateTrackDuration(ctx, track.ID, durationSeconds)
			if err != nil {
				l.Warn().Err(err).Msgf("BackfillTrackDurations: Failed to update duration for track %d", track.ID)
				continue
			}

			l.Debug().Msgf("BackfillTrackDurations: Updated track %d with duration %ds", track.ID, durationSeconds)
			totalProcessed++
		}
	}

	l.Info().Msgf("BackfillTrackDurations: Completed. Updated %d tracks with duration", totalProcessed)
	return nil
}
