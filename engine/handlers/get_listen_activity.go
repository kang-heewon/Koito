package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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
			Timezone: parseTZ(r),
			AlbumID:  int32(albumId),
			ArtistID: int32(artistId),
			TrackID:  int32(trackId),
		}

		if strings.ToLower(opts.Timezone.String()) == "local" {
			opts.Timezone, _ = time.LoadLocation("UTC")
			l.Warn().Msg("GetListenActivityHandler: Timezone is unset, using UTC")
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

		activity = processActivity(activity, opts)

		l.Debug().Msg("GetListenActivityHandler: Successfully retrieved listen activity")
		utils.WriteJSON(w, http.StatusOK, activity)
	}
}

// ngl i hate this
func processActivity(
	items []db.ListenActivityItem,
	opts db.ListenActivityOpts,
) []db.ListenActivityItem {
	from, to := db.ListenActivityOptsToTimes(opts)

	buckets := make(map[string]int64)

	for _, item := range items {
		bucketStart := normalizeToStep(item.Start, opts.Step)
		key := bucketStart.Format("2006-01-02")
		buckets[key] += item.Listens
	}

	var result []db.ListenActivityItem

	for t := normalizeToStep(from, opts.Step); t.Before(to); t = addStep(t, opts.Step) {
		key := t.Format("2006-01-02")

		result = append(result, db.ListenActivityItem{
			Start:   t,
			Listens: buckets[key],
		})
	}

	return result
}

func normalizeToStep(t time.Time, step db.StepInterval) time.Time {
	switch step {
	case db.StepDay:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	case db.StepWeek:
		weekday := int(t.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := t.AddDate(0, 0, -(weekday - 1))
		return time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, t.Location())

	case db.StepMonth:
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	default:
		return t
	}
}

func addStep(t time.Time, step db.StepInterval) time.Time {
	switch step {
	case db.StepDay:
		return t.AddDate(0, 0, 1)
	case db.StepWeek:
		return t.AddDate(0, 0, 7)
	case db.StepMonth:
		return t.AddDate(0, 1, 0)
	default:
		return t.AddDate(0, 0, 1)
	}
}
