package psql

import (
	"context"
	"fmt"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/repository"
)

func (d *Psql) GetGenreStatsByListenCount(ctx context.Context, timeframe db.Timeframe) ([]db.GenreStat, error) {
	t1, t2 := db.TimeframeToTimeRange(timeframe)

	rows, err := d.q.GetGenreStatsByListenCount(ctx, repository.GetGenreStatsByListenCountParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return nil, fmt.Errorf("GetGenreStatsByListenCount: %w", err)
	}

	stats := make([]db.GenreStat, len(rows))
	for i, row := range rows {
		stats[i] = db.GenreStat{
			Name:  row.Name,
			Value: row.ListenCount,
		}
	}
	return stats, nil
}

func (d *Psql) GetGenreStatsByTimeListened(ctx context.Context, timeframe db.Timeframe) ([]db.GenreStat, error) {
	t1, t2 := db.TimeframeToTimeRange(timeframe)

	rows, err := d.q.GetGenreStatsByTimeListened(ctx, repository.GetGenreStatsByTimeListenedParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return nil, fmt.Errorf("GetGenreStatsByTimeListened: %w", err)
	}

	stats := make([]db.GenreStat, len(rows))
	for i, row := range rows {
		stats[i] = db.GenreStat{
			Name:  row.Name,
			Value: row.SecondsListened,
		}
	}
	return stats, nil
}
