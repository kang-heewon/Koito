package psql

import (
	"context"
	"fmt"
	"time"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/repository"
)

func (d *Psql) GetGenreStatsByListenCount(ctx context.Context, period db.Period) ([]db.GenreStat, error) {
	startTime := db.StartTimeFromPeriod(period)
	endTime := time.Now()

	rows, err := d.q.GetGenreStatsByListenCount(ctx, repository.GetGenreStatsByListenCountParams{
		ListenedAt:   startTime,
		ListenedAt_2: endTime,
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

func (d *Psql) GetGenreStatsByTimeListened(ctx context.Context, period db.Period) ([]db.GenreStat, error) {
	startTime := db.StartTimeFromPeriod(period)
	endTime := time.Now()

	rows, err := d.q.GetGenreStatsByTimeListened(ctx, repository.GetGenreStatsByTimeListenedParams{
		ListenedAt:   startTime,
		ListenedAt_2: endTime,
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
