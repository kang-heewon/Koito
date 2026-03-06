package psql

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/repository"
)

func (d *Psql) AlbumsWithoutMbzID(ctx context.Context, from int32) ([]*models.Album, error) {
	l := logger.FromContext(ctx)
	rows, err := d.q.GetReleasesWithoutMbzID(ctx, repository.GetReleasesWithoutMbzIDParams{
		Limit: 20,
		ID:    from,
	})
	if err != nil {
		return nil, fmt.Errorf("AlbumsWithoutMbzID: GetReleasesWithoutMbzID: %w", err)
	}
	albums := make([]*models.Album, len(rows))
	for i, row := range rows {
		var artists []models.SimpleArtist
		if err := json.Unmarshal(row.Artists, &artists); err != nil {
			l.Err(err).Msgf("AlbumsWithoutMbzID: error unmarshalling artists for release %d", row.ID)
			artists = nil
		}
		albums[i] = &models.Album{
			ID:      row.ID,
			Title:   row.Title,
			Artists: artists,
		}
	}
	return albums, nil
}

func (d *Psql) MarkMbzSearched(ctx context.Context, id int32) error {
	return d.q.MarkMbzSearched(ctx, id)
}
