package psql

import (
	"context"
	"encoding/json"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/repository"
)

func (d *Psql) GetWrappedStats(ctx context.Context, year int, userID int32) (*db.WrappedStats, error) {
	l := logger.FromContext(ctx)
	parseArtists := func(raw []byte, source string) []models.SimpleArtist {
		artists := make([]models.SimpleArtist, 0)
		if len(raw) == 0 {
			return artists
		}
		if err := json.Unmarshal(raw, &artists); err != nil {
			l.Warn().Err(err).Str("source", source).Msg("Failed to parse artists JSON")
			return make([]models.SimpleArtist, 0)
		}
		return artists
	}

	stats := &db.WrappedStats{
		Year: year,
	}

	totalListens, err := d.q.GetTotalListenCountInYear(ctx, repository.GetTotalListenCountInYearParams{
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Err(err).Msg("Failed to get total listen count")
	} else {
		stats.TotalListens = totalListens
	}

	totalSeconds, err := d.q.GetTotalTimeListenedInYear(ctx, repository.GetTotalTimeListenedInYearParams{
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Err(err).Msg("Failed to get total time listened")
	} else {
		stats.TotalSecondsListened = totalSeconds
	}

	uniqueArtists, err := d.q.GetArtistCountInYear(ctx, repository.GetArtistCountInYearParams{
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Err(err).Msg("Failed to get unique artists")
	} else {
		stats.UniqueArtists = uniqueArtists
	}

	uniqueTracks, err := d.q.GetTrackCountInYear(ctx, repository.GetTrackCountInYearParams{
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Err(err).Msg("Failed to get unique tracks")
	} else {
		stats.UniqueTracks = uniqueTracks
	}

	uniqueAlbums, err := d.q.GetAlbumCountInYear(ctx, repository.GetAlbumCountInYearParams{
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Err(err).Msg("Failed to get unique albums")
	} else {
		stats.UniqueAlbums = uniqueAlbums
	}

	topTracks, err := d.q.GetTopTracksInYear(ctx, repository.GetTopTracksInYearParams{
		Limit:  10,
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Err(err).Msg("Failed to get top tracks")
	} else {
		stats.TopTracks = make([]*models.Track, len(topTracks))
		for i, row := range topTracks {
			artists := parseArtists(row.Artists, "GetTopTracksInYear")
			stats.TopTracks[i] = &models.Track{
				ID:          row.TrackID,
				Title:       row.Title,
				Artists:     artists,
				ListenCount: row.ListenCount,
			}
		}
	}

	topNewArtists, err := d.q.GetTopThreeNewArtistsInYear(ctx, repository.GetTopThreeNewArtistsInYearParams{
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Debug().Err(err).Msg("Failed to get top new artists")
	} else {
		stats.TopNewArtists = make([]*models.Artist, len(topNewArtists))
		for i, row := range topNewArtists {
			stats.TopNewArtists[i] = &models.Artist{
				ID:          row.ArtistID,
				Name:        row.ArtistName,
				ListenCount: row.TotalPlaysInYear,
			}
		}
	}

	topArtists, err := d.q.GetTopArtistsInYear(ctx, repository.GetTopArtistsInYearParams{
		Limit:  10,
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Err(err).Msg("Failed to get top artists")
	} else {
		stats.TopArtists = make([]*models.Artist, len(topArtists))
		for i, row := range topArtists {
			stats.TopArtists[i] = &models.Artist{
				ID:          row.ArtistID,
				Name:        row.Name,
				ListenCount: row.ListenCount,
			}
		}
	}

	topAlbums, err := d.q.GetTopAlbumsInYear(ctx, repository.GetTopAlbumsInYearParams{
		Limit:  10,
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Err(err).Msg("Failed to get top albums")
	} else {
		stats.TopAlbums = make([]*models.Album, len(topAlbums))
		for i, row := range topAlbums {
			stats.TopAlbums[i] = &models.Album{
				ID:          row.ReleaseID,
				Title:       row.Title,
				ListenCount: row.ListenCount,
			}
		}
	}

	mostReplayed, err := d.q.GetMostReplayedTrackInYear(ctx, repository.GetMostReplayedTrackInYearParams{
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Debug().Err(err).Msg("Failed to get most replayed track")
	} else {
		artists := parseArtists(mostReplayed.Artists, "GetMostReplayedTrackInYear")
		stats.MostReplayedTrack = &db.TrackStreak{
			Track: &models.Track{
				ID:       mostReplayed.ID,
				MbzID:    mostReplayed.MusicBrainzID,
				Title:    mostReplayed.Title,
				Duration: mostReplayed.Duration,
				AlbumID:  mostReplayed.ReleaseID,
				Artists:  artists,
			},
			StreakCount: int(mostReplayed.StreakLength),
		}
	}

	hours, err := d.q.GetListeningHoursDistributionInYear(ctx, repository.GetListeningHoursDistributionInYearParams{
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Err(err).Msg("Failed to get hourly distribution")
	} else {
		stats.ListeningHours = make([]db.HourDistribution, len(hours))
		for i, h := range hours {
			stats.ListeningHours[i] = db.HourDistribution{
				Hour:        int(h.Hour),
				ListenCount: h.Count,
			}
		}
	}

	busiestWeek, err := d.q.GetWeekWithMostListensInYear(ctx, repository.GetWeekWithMostListensInYearParams{
		Year:   int32(year),
		UserID: userID,
	})
	if err != nil {
		l.Debug().Err(err).Msg("Failed to get week with most listens")
	} else {
		stats.BusiestWeek = &db.WeekStats{
			WeekStart:   busiestWeek.WeekStart,
			ListenCount: busiestWeek.ListenCount,
		}
	}

	firstListen, err := d.q.GetFirstListenInYear(ctx, repository.GetFirstListenInYearParams{
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Err(err).Msg("Failed to get first listen")
	} else {
		artists := parseArtists(firstListen.Artists, "GetFirstListenInYear")
		stats.FirstListen = &models.Listen{
			Time: firstListen.ListenedAt,
			Track: models.Track{
				ID:       firstListen.TrackID,
				MbzID:    firstListen.MusicBrainzID,
				Title:    firstListen.Title.String,
				Duration: firstListen.Duration.Int32,
				AlbumID:  firstListen.ReleaseID.Int32,
				Artists:  artists,
			},
		}
	}

	everyMonth, err := d.q.GetTracksPlayedAtLeastOncePerMonthInYear(ctx, repository.GetTracksPlayedAtLeastOncePerMonthInYearParams{
		UserID: userID,
		Year:   int32(year),
	})
	if err != nil {
		l.Err(err).Msg("Failed to get tracks played every month")
	} else {
		stats.TracksPlayedEveryMonth = make([]*models.Track, len(everyMonth))
		for i, row := range everyMonth {
			stats.TracksPlayedEveryMonth[i] = &models.Track{
				ID:    row.TrackID,
				Title: row.Title,
			}
		}
	}

	artistConc, err := d.q.GetPercentageOfTotalListensFromTopArtistsInYear(ctx, repository.GetPercentageOfTotalListensFromTopArtistsInYearParams{
		Limit:  5,
		UserID: userID,
		Year:   int32(year),
	})
	if err == nil {
		f, _ := artistConc.PercentOfTotal.Float64Value()
		stats.ArtistConcentration = f.Float64
	}

	trackConc, err := d.q.GetPercentageOfTotalListensFromTopTracksInYear(ctx, repository.GetPercentageOfTotalListensFromTopTracksInYearParams{
		Limit:  5,
		UserID: userID,
		Year:   int32(year),
	})
	if err == nil {
		f, _ := trackConc.PercentOfTotal.Float64Value()
		stats.TrackConcentration = f.Float64
	}

	return stats, nil
}
