package psql

import (
	"context"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/repository"
	"github.com/google/uuid"
)

func (d *Psql) TracksWithoutDuration(ctx context.Context, lastID int32) ([]db.TrackWithMbzID, error) {
	rows, err := d.q.TracksWithoutDuration(ctx, lastID)
	if err != nil {
		return nil, err
	}

	result := make([]db.TrackWithMbzID, len(rows))
	for i, row := range rows {
		// row.MusicBrainzID is *uuid.UUID, but query ensures IS NOT NULL
		var mbzID uuid.UUID
		if row.MusicBrainzID != nil {
			mbzID = *row.MusicBrainzID
		}

		result[i] = db.TrackWithMbzID{
			ID:    row.ID,
			MbzID: mbzID,
		}
	}
	return result, nil
}

func (d *Psql) UpdateTrackDuration(ctx context.Context, id int32, duration int32) error {
	return d.q.UpdateTrackDuration(ctx, repository.UpdateTrackDurationParams{
		ID:       id,
		Duration: duration,
	})
}
