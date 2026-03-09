package psql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/repository"
)

const (
	countNewTracksQuery = `
		SELECT COUNT(*)
		FROM (
			SELECT l.track_id
			FROM listens l
			GROUP BY l.track_id
			HAVING MIN(l.listened_at) BETWEEN $1 AND $2
		) new_tracks;
	`
	countNewAlbumsQuery = `
		SELECT COUNT(*)
		FROM (
			SELECT t.release_id
			FROM listens l
			JOIN tracks t ON t.id = l.track_id
			GROUP BY t.release_id
			HAVING MIN(l.listened_at) BETWEEN $1 AND $2
		) new_albums;
	`
	countNewArtistsQuery = `
		SELECT COUNT(*)
		FROM (
			SELECT at.artist_id
			FROM listens l
			JOIN artist_tracks at ON at.track_id = l.track_id
			GROUP BY at.artist_id
			HAVING MIN(l.listened_at) BETWEEN $1 AND $2
		) new_artists;
	`
)

func (p *Psql) CountListens(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	t1, t2 := db.TimeframeToTimeRange(timeframe)
	count, err := p.q.CountListens(ctx, repository.CountListensParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return 0, fmt.Errorf("CountListens: %w", err)
	}
	return count, nil
}

func (p *Psql) CountTracks(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	t1, t2 := db.TimeframeToTimeRange(timeframe)
	count, err := p.q.CountTopTracks(ctx, repository.CountTopTracksParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return 0, fmt.Errorf("CountTracks: %w", err)
	}
	return count, nil
}

func (p *Psql) CountAlbums(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	t1, t2 := db.TimeframeToTimeRange(timeframe)
	count, err := p.q.CountTopReleases(ctx, repository.CountTopReleasesParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return 0, fmt.Errorf("CountAlbums: %w", err)
	}
	return count, nil
}

func (p *Psql) CountArtists(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	t1, t2 := db.TimeframeToTimeRange(timeframe)
	count, err := p.q.CountTopArtists(ctx, repository.CountTopArtistsParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return 0, fmt.Errorf("CountArtists: %w", err)
	}
	return count, nil
}

func (p *Psql) CountTimeListened(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	t1, t2 := db.TimeframeToTimeRange(timeframe)
	count, err := p.q.CountTimeListened(ctx, repository.CountTimeListenedParams{
		ListenedAt:   t1,
		ListenedAt_2: t2,
	})
	if err != nil {
		return 0, fmt.Errorf("CountTimeListened: %w", err)
	}
	return count, nil
}

func (p *Psql) CountTimeListenedToItem(ctx context.Context, opts db.TimeListenedOpts) (int64, error) {
	t1, t2 := db.TimeframeToTimeRange(opts.Timeframe)

	if opts.ArtistID > 0 {
		count, err := p.q.CountTimeListenedToArtist(ctx, repository.CountTimeListenedToArtistParams{
			ListenedAt:   t1,
			ListenedAt_2: t2,
			ArtistID:     opts.ArtistID,
		})
		if err != nil {
			return 0, fmt.Errorf("CountTimeListenedToItem (Artist): %w", err)
		}
		return count, nil
	} else if opts.AlbumID > 0 {
		count, err := p.q.CountTimeListenedToRelease(ctx, repository.CountTimeListenedToReleaseParams{
			ListenedAt:   t1,
			ListenedAt_2: t2,
			ReleaseID:    opts.AlbumID,
		})
		if err != nil {
			return 0, fmt.Errorf("CountTimeListenedToItem (Album): %w", err)
		}
		return count, nil
	} else if opts.TrackID > 0 {
		count, err := p.q.CountTimeListenedToTrack(ctx, repository.CountTimeListenedToTrackParams{
			ListenedAt:   t1,
			ListenedAt_2: t2,
			ID:           opts.TrackID,
		})
		if err != nil {
			return 0, fmt.Errorf("CountTimeListenedToItem (Track): %w", err)
		}
		return count, nil
	}
	return 0, errors.New("CountTimeListenedToItem: an id must be provided")
}

func (p *Psql) CountListensToItem(ctx context.Context, opts db.TimeListenedOpts) (int64, error) {
	t1, t2 := db.TimeframeToTimeRange(opts.Timeframe)

	if opts.ArtistID > 0 {
		count, err := p.q.CountListensFromArtist(ctx, repository.CountListensFromArtistParams{
			ListenedAt:   t1,
			ListenedAt_2: t2,
			ArtistID:     opts.ArtistID,
		})
		if err != nil {
			return 0, fmt.Errorf("CountListensToItem (Artist): %w", err)
		}
		return count, nil
	} else if opts.AlbumID > 0 {
		count, err := p.q.CountListensFromRelease(ctx, repository.CountListensFromReleaseParams{
			ListenedAt:   t1,
			ListenedAt_2: t2,
			ReleaseID:    opts.AlbumID,
		})
		if err != nil {
			return 0, fmt.Errorf("CountListensToItem (Album): %w", err)
		}
		return count, nil
	} else if opts.TrackID > 0 {
		count, err := p.q.CountListensFromTrack(ctx, repository.CountListensFromTrackParams{
			ListenedAt:   t1,
			ListenedAt_2: t2,
			TrackID:      opts.TrackID,
		})
		if err != nil {
			return 0, fmt.Errorf("CountListensToItem (Track): %w", err)
		}
		return count, nil
	}
	return 0, errors.New("CountListensToItem: an id must be provided")
}

func (p *Psql) CountNewTracks(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	t1, t2 := db.TimeframeToTimeRange(timeframe)
	count, err := p.countWithTimeRange(ctx, countNewTracksQuery, t1, t2)
	if err != nil {
		return 0, fmt.Errorf("CountNewTracks: %w", err)
	}
	return count, nil
}

func (p *Psql) CountNewAlbums(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	t1, t2 := db.TimeframeToTimeRange(timeframe)
	count, err := p.countWithTimeRange(ctx, countNewAlbumsQuery, t1, t2)
	if err != nil {
		return 0, fmt.Errorf("CountNewAlbums: %w", err)
	}
	return count, nil
}

func (p *Psql) CountNewArtists(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	t1, t2 := db.TimeframeToTimeRange(timeframe)
	count, err := p.countWithTimeRange(ctx, countNewArtistsQuery, t1, t2)
	if err != nil {
		return 0, fmt.Errorf("CountNewArtists: %w", err)
	}
	return count, nil
}

func (p *Psql) countWithTimeRange(ctx context.Context, query string, t1 time.Time, t2 time.Time) (int64, error) {
	var count int64
	if p.tx != nil {
		if err := p.tx.QueryRow(ctx, query, t1, t2).Scan(&count); err != nil {
			return 0, err
		}
		return count, nil
	}
	if err := p.conn.QueryRow(ctx, query, t1, t2).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
