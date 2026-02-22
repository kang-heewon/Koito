// package handlers implements route handlers
package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

const defaultLimitSize = 100
const maximumLimit = 500

func OptsFromRequest(r *http.Request) (db.GetItemsOpts, error) {
	l := logger.FromContext(r.Context())

	l.Debug().Msg("OptsFromRequest: Parsing query parameters")

	parseOptionalInt := func(key string) (int, error) {
		value := strings.TrimSpace(r.URL.Query().Get(key))
		if value == "" {
			return 0, nil
		}

		parsed, err := strconv.Atoi(value)
		if err != nil {
			l.Debug().Msgf("OptsFromRequest: Invalid integer for query parameter '%s': %q", key, value)
			return 0, fmt.Errorf("invalid %s parameter", key)
		}

		if parsed < 0 {
			l.Debug().Msgf("OptsFromRequest: Negative integer for query parameter '%s': %q", key, value)
			return 0, fmt.Errorf("invalid %s parameter", key)
		}

		return parsed, nil
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

	week, err := parseOptionalInt("week")
	if err != nil {
		return db.GetItemsOpts{}, err
	}

	month, err := parseOptionalInt("month")
	if err != nil {
		return db.GetItemsOpts{}, err
	}

	year, err := parseOptionalInt("year")
	if err != nil {
		return db.GetItemsOpts{}, err
	}

	from, err := parseOptionalInt("from")
	if err != nil {
		return db.GetItemsOpts{}, err
	}

	to, err := parseOptionalInt("to")
	if err != nil {
		return db.GetItemsOpts{}, err
	}

	artistId, err := parseOptionalInt("artist_id")
	if err != nil {
		return db.GetItemsOpts{}, err
	}

	albumId, err := parseOptionalInt("album_id")
	if err != nil {
		return db.GetItemsOpts{}, err
	}

	trackId, err := parseOptionalInt("track_id")
	if err != nil {
		return db.GetItemsOpts{}, err
	}

	if (from == 0) != (to == 0) {
		return db.GetItemsOpts{}, fmt.Errorf("from and to must be specified together")
	}

	if from != 0 && to != 0 && from > to {
		return db.GetItemsOpts{}, fmt.Errorf("from must be less than or equal to to")
	}

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
	}, nil
}

func isDateRangeValidationError(err error) bool {
	var dateRangeErr *utils.DateRangeValidationError
	return errors.As(err, &dateRangeErr)
}
