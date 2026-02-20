// package handlers implements route handlers
package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
)

const defaultLimitSize = 100
const maximumLimit = 500

func OptsFromRequest(r *http.Request) db.GetItemsOpts {
	l := logger.FromContext(r.Context())

	l.Debug().Msg("OptsFromRequest: Parsing query parameters")

	parseOptionalInt := func(key string) int {
		value := strings.TrimSpace(r.URL.Query().Get(key))
		if value == "" {
			return 0
		}

		parsed, err := strconv.Atoi(value)
		if err != nil {
			l.Debug().Msgf("OptsFromRequest: Invalid integer for query parameter '%s': %q", key, value)
			return 0
		}

		return parsed
	}

	limitStr := strings.TrimSpace(r.URL.Query().Get("limit"))
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		l.Debug().Msgf("OptsFromRequest: Query parameter 'limit' not specified, using default %d", defaultLimitSize)
		limit = defaultLimitSize
	}
	if limit > maximumLimit {
		l.Debug().Msgf("OptsFromRequest: Limit exceeds maximum %d, clamping to %d", maximumLimit, maximumLimit)
		limit = maximumLimit
	}

	pageStr := strings.TrimSpace(r.URL.Query().Get("page"))
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		l.Debug().Msg("OptsFromRequest: Page parameter is less than 1, defaulting to 1")
		page = 1
	}

	week := parseOptionalInt("week")
	month := parseOptionalInt("month")
	year := parseOptionalInt("year")
	from := parseOptionalInt("from")
	to := parseOptionalInt("to")

	artistId := parseOptionalInt("artist_id")
	albumId := parseOptionalInt("album_id")
	trackId := parseOptionalInt("track_id")

	var period db.Period
	switch strings.ToLower(r.URL.Query().Get("period")) {
	case "day":
		period = db.PeriodDay
	case "week":
		period = db.PeriodWeek
	case "month":
		period = db.PeriodMonth
	case "year":
		period = db.PeriodYear
	case "all_time":
		period = db.PeriodAllTime
	default:
		l.Debug().Msgf("OptsFromRequest: Using default value '%s' for period", db.PeriodDay)
		period = db.PeriodDay
	}

	l.Debug().Msgf("OptsFromRequest: Parsed options: limit=%d, page=%d, week=%d, month=%d, year=%d, from=%d, to=%d, artist_id=%d, album_id=%d, track_id=%d, period=%s",
		limit, page, week, month, year, from, to, artistId, albumId, trackId, period)

	return db.GetItemsOpts{
		Limit:    limit,
		Period:   period,
		Page:     page,
		Week:     week,
		Month:    month,
		Year:     year,
		From:     from,
		To:       to,
		ArtistID: artistId,
		AlbumID:  albumId,
		TrackID:  trackId,
	}
}
