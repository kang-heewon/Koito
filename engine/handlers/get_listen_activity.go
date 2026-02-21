package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

func GetListenActivityHandler(store db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msg("GetListenActivityHandler: Received request to retrieve listen activity")

		parseOptionalInt := func(key string, invalidMessage string) (int, bool) {
			value := strings.TrimSpace(r.URL.Query().Get(key))
			if value == "" {
				return 0, true
			}

			parsed, err := strconv.Atoi(value)
			if err != nil {
				l.Debug().AnErr("error", err).Msgf("GetListenActivityHandler: Invalid %s parameter", key)
				utils.WriteError(w, invalidMessage, http.StatusBadRequest)
				return 0, false
			}

			return parsed, true
		}

		_range, ok := parseOptionalInt("range", "invalid range parameter")
		if !ok {
			return
		}

		month, ok := parseOptionalInt("month", "invalid month parameter")
		if !ok {
			return
		}

		year, ok := parseOptionalInt("year", "invalid year parameter")
		if !ok {
			return
		}

		artistId, ok := parseOptionalInt("artist_id", "invalid artist ID parameter")
		if !ok {
			return
		}

		albumId, ok := parseOptionalInt("album_id", "invalid album ID parameter")
		if !ok {
			return
		}

		trackId, ok := parseOptionalInt("track_id", "invalid track ID parameter")
		if !ok {
			return
		}

		var step db.StepInterval
		switch strings.ToLower(r.URL.Query().Get("step")) {
		case "day":
			step = db.StepDay
		case "week":
			step = db.StepWeek
		case "month":
			step = db.StepMonth
		case "year":
			step = db.StepYear
		default:
			l.Debug().Msgf("GetListenActivityHandler: Using default value '%s' for step", db.StepDefault)
			step = db.StepDay
		}

		opts := db.ListenActivityOpts{
			Step:     step,
			Range:    _range,
			Month:    month,
			Year:     year,
			AlbumID:  int32(albumId),
			ArtistID: int32(artistId),
			TrackID:  int32(trackId),
		}

		l.Debug().Msgf("GetListenActivityHandler: Retrieving listen activity with options: %+v", opts)

		activity, err := store.GetListenActivity(ctx, opts)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "year must be specified with month") {
				l.Debug().AnErr("error", err).Msg("GetListenActivityHandler: Invalid month/year combination")
				utils.WriteError(w, err.Error(), http.StatusBadRequest)
				return
			}

			l.Err(err).Msg("GetListenActivityHandler: Failed to retrieve listen activity")
			utils.WriteError(w, "failed to retrieve listen activity", http.StatusInternalServerError)
			return
		}

		l.Debug().Msg("GetListenActivityHandler: Successfully retrieved listen activity")
		utils.WriteJSON(w, http.StatusOK, activity)
	}
}
