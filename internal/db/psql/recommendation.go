package psql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/repository"
)

func (d *Psql) GetTracksToRevisit(ctx context.Context, opts db.GetRecommendationsOpts) ([]db.TrackRecommendation, error) {
	l := logger.FromContext(ctx)
	l.Debug().Msg("Fetching tracks to revisit")

	rows, err := d.q.GetTracksToRevisit(ctx, repository.GetTracksToRevisitParams{
		ListenedAt:   opts.PastWindowStart,
		ListenedAt_2: opts.PastWindowEnd,
		Limit:        int32(opts.Limit),
	})
	if err != nil {
		return nil, fmt.Errorf("GetTracksToRevisit: repository.GetTracksToRevisit: %w", err)
	}

	recommendations := make([]db.TrackRecommendation, len(rows))
	for i, row := range rows {
		track := models.Track{
			ID:          row.TrackID,
			Title:       row.Title,
			ListenCount: 0,
			AlbumID:     row.ReleaseID,
		}

		err = json.Unmarshal(row.Artists, &track.Artists)
		if err != nil {
			return nil, fmt.Errorf("GetTracksToRevisit: json.Unmarshal artists for track %d: %w", row.TrackID, err)
		}

		track.Image = row.ReleaseImage

		var lastListened time.Time
		if row.LastListenedAt != nil {
			if t, ok := row.LastListenedAt.(time.Time); ok {
				lastListened = t
			} else {
				l.Warn().Msgf("GetTracksToRevisit: unexpected type for LastListenedAt: %T", row.LastListenedAt)
			}
		}

		recommendations[i] = db.TrackRecommendation{
			Track:           &track,
			PastListenCount: row.PastListenCount,
			LastListenedAt:  lastListened,
		}
	}

	return recommendations, nil
}
