package psql

import (
	"context"
	"fmt"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/repository"
)

func (d *Psql) SaveAlbumGenres(ctx context.Context, id int32, genres []string) error {
	l := logger.FromContext(ctx)
	if id == 0 {
		return fmt.Errorf("SaveAlbumGenres: album id not specified")
	}
	if len(genres) == 0 {
		return nil
	}

	tx, qtx, ownsTx, err := d.withTx(ctx)
	if err != nil {
		l.Err(err).Msg("Failed to begin transaction")
		return fmt.Errorf("SaveAlbumGenres: BeginTx: %w", err)
	}
	if ownsTx {
		defer tx.Rollback(ctx)
	}

	for _, genreName := range genres {
		genre, err := qtx.InsertGenre(ctx, genreName)
		if err != nil {
			return fmt.Errorf("SaveAlbumGenres: InsertGenre: %w", err)
		}
		err = qtx.AssociateGenreToRelease(ctx, repository.AssociateGenreToReleaseParams{
			ReleaseID: id,
			GenreID:   genre.ID,
		})
		if err != nil {
			return fmt.Errorf("SaveAlbumGenres: AssociateGenreToRelease: %w", err)
		}
	}

	if ownsTx {
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("SaveAlbumGenres: Commit: %w", err)
		}
	}

	return nil
}

func (d *Psql) SaveArtistGenres(ctx context.Context, id int32, genres []string) error {
	l := logger.FromContext(ctx)
	if id == 0 {
		return fmt.Errorf("SaveArtistGenres: artist id not specified")
	}
	if len(genres) == 0 {
		return nil
	}

	tx, qtx, ownsTx, err := d.withTx(ctx)
	if err != nil {
		l.Err(err).Msg("Failed to begin transaction")
		return fmt.Errorf("SaveArtistGenres: BeginTx: %w", err)
	}
	if ownsTx {
		defer tx.Rollback(ctx)
	}

	for _, genreName := range genres {
		genre, err := qtx.InsertGenre(ctx, genreName)
		if err != nil {
			return fmt.Errorf("SaveArtistGenres: InsertGenre: %w", err)
		}
		err = qtx.AssociateGenreToArtist(ctx, repository.AssociateGenreToArtistParams{
			ArtistID: id,
			GenreID:  genre.ID,
		})
		if err != nil {
			return fmt.Errorf("SaveArtistGenres: AssociateGenreToArtist: %w", err)
		}
	}

	if ownsTx {
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("SaveArtistGenres: Commit: %w", err)
		}
	}

	return nil
}

func (d *Psql) getGenresForRelease(ctx context.Context, id int32) ([]string, error) {
	rows, err := d.q.GetGenresForRelease(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getGenresForRelease: %w", err)
	}
	genres := make([]string, len(rows))
	for i, row := range rows {
		genres[i] = row.Name
	}
	return genres, nil
}

func (d *Psql) getGenresForArtist(ctx context.Context, id int32) ([]string, error) {
	rows, err := d.q.GetGenresForArtist(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getGenresForArtist: %w", err)
	}
	genres := make([]string, len(rows))
	for i, row := range rows {
		genres[i] = row.Name
	}
	return genres, nil
}

func (d *Psql) AlbumsWithoutGenres(ctx context.Context, from int32) ([]db.ItemWithMbzID, error) {
	rows, err := d.q.GetReleasesWithoutGenres(ctx, repository.GetReleasesWithoutGenresParams{
		Limit: 100,
		ID:    from,
	})
	if err != nil {
		return nil, fmt.Errorf("AlbumsWithoutGenres: %w", err)
	}
	items := make([]db.ItemWithMbzID, len(rows))
	for i, row := range rows {
		items[i] = db.ItemWithMbzID{
			ID:    row.ID,
			MbzID: *row.MusicBrainzID,
		}
	}
	return items, nil
}

func (d *Psql) ArtistsWithoutGenres(ctx context.Context, from int32) ([]db.ItemWithMbzID, error) {
	rows, err := d.q.GetArtistsWithoutGenres(ctx, repository.GetArtistsWithoutGenresParams{
		Limit: 100,
		ID:    from,
	})
	if err != nil {
		return nil, fmt.Errorf("ArtistsWithoutGenres: %w", err)
	}
	items := make([]db.ItemWithMbzID, len(rows))
	for i, row := range rows {
		items[i] = db.ItemWithMbzID{
			ID:    row.ID,
			MbzID: *row.MusicBrainzID,
		}
	}
	return items, nil
}
