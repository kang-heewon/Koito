package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/utils"
)

type WrappedResponse struct {
	Year                   int                        `json:"year"`
	TotalListens           int64                      `json:"total_listens"`
	TotalSecondsListened   int64                      `json:"total_seconds_listened"`
	UniqueArtists          int64                      `json:"unique_artists"`
	UniqueTracks           int64                      `json:"unique_tracks"`
	UniqueAlbums           int64                      `json:"unique_albums"`
	TopTracks              []*models.Track            `json:"top_tracks"`
	TopArtists             []*models.Artist           `json:"top_artists"`
	TopAlbums              []*models.Album            `json:"top_albums"`
	TopNewArtists          []*models.Artist           `json:"top_new_artists"`
	MostReplayedTrack      *TrackStreakResponse       `json:"most_replayed_track"`
	ListeningHours         []HourDistributionResponse `json:"listening_hours"`
	BusiestWeek            *WeekStatsResponse         `json:"busiest_week"`
	FirstListen            *models.Listen             `json:"first_listen"`
	TracksPlayedEveryMonth []*models.Track            `json:"tracks_played_every_month"`
	ArtistConcentration    float64                    `json:"artist_concentration"`
	TrackConcentration     float64                    `json:"track_concentration"`
}

type TrackStreakResponse struct {
	Track       *models.Track `json:"track"`
	StreakCount int           `json:"streak_count"`
}

type HourDistributionResponse struct {
	Hour        int   `json:"hour"`
	ListenCount int64 `json:"listen_count"`
}

type WeekStatsResponse struct {
	WeekStart   string `json:"week_start"`
	ListenCount int64  `json:"listen_count"`
}

func WrappedHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		l.Debug().Msg("WrappedHandler: Received request to retrieve wrapped stats")

		yearStr := r.URL.Query().Get("year")
		year := time.Now().Year()
		if yearStr != "" {
			if y, err := strconv.Atoi(yearStr); err == nil {
				year = y
			}
		}

		l.Debug().Msgf("WrappedHandler: Fetching wrapped stats for year '%d'", year)

		userID := int32(1)

		stats, err := store.GetWrappedStats(r.Context(), year, userID)
		if err != nil {
			l.Err(err).Msg("WrappedHandler: Failed to fetch wrapped stats")
			utils.WriteError(w, "failed to get wrapped stats: "+err.Error(), http.StatusInternalServerError)
			return
		}

		l.Debug().Msg("WrappedHandler: Successfully fetched wrapped stats")
		utils.WriteJSON(w, http.StatusOK, mapWrappedStatsToResponse(stats))
	}
}

func mapWrappedStatsToResponse(stats *db.WrappedStats) *WrappedResponse {
	if stats == nil {
		return nil
	}

	response := &WrappedResponse{
		Year:                   stats.Year,
		TotalListens:           stats.TotalListens,
		TotalSecondsListened:   stats.TotalSecondsListened,
		UniqueArtists:          stats.UniqueArtists,
		UniqueTracks:           stats.UniqueTracks,
		UniqueAlbums:           stats.UniqueAlbums,
		TopTracks:              stats.TopTracks,
		TopArtists:             stats.TopArtists,
		TopAlbums:              stats.TopAlbums,
		TopNewArtists:          stats.TopNewArtists,
		ListeningHours:         mapHourDistribution(stats.ListeningHours),
		FirstListen:            stats.FirstListen,
		TracksPlayedEveryMonth: stats.TracksPlayedEveryMonth,
		ArtistConcentration:    stats.ArtistConcentration,
		TrackConcentration:     stats.TrackConcentration,
	}

	if stats.MostReplayedTrack != nil {
		response.MostReplayedTrack = &TrackStreakResponse{
			Track:       stats.MostReplayedTrack.Track,
			StreakCount: stats.MostReplayedTrack.StreakCount,
		}
	}

	if stats.BusiestWeek != nil {
		response.BusiestWeek = &WeekStatsResponse{
			WeekStart:   stats.BusiestWeek.WeekStart.Format(time.RFC3339),
			ListenCount: stats.BusiestWeek.ListenCount,
		}
	}

	return response
}

func mapHourDistribution(hours []db.HourDistribution) []HourDistributionResponse {
	if hours == nil {
		return nil
	}

	result := make([]HourDistributionResponse, len(hours))
	for i, h := range hours {
		result[i] = HourDistributionResponse{
			Hour:        h.Hour,
			ListenCount: h.ListenCount,
		}
	}
	return result
}
