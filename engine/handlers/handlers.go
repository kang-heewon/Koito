// package handlers implements route handlers
package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/gabehf/koito/internal/cfg"
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

	tf := TimeframeFromRequest(r)
	if week != 0 {
		tf.Week = week
	}
	if month != 0 {
		tf.Month = month
	}
	if year != 0 {
		tf.Year = year
	}
	if from != 0 {
		tf.FromUnix = int64(from)
	}
	if to != 0 {
		tf.ToUnix = int64(to)
	}

	l.Debug().Msgf("OptsFromRequest: Parsed options: limit=%d, page=%d, week=%d, month=%d, year=%d, from=%d, to=%d, artist_id=%d, album_id=%d, track_id=%d, period=%s",
		limit, page, tf.Week, tf.Month, tf.Year, tf.FromUnix, tf.ToUnix, artistId, albumId, trackId, tf.Period)

	return db.GetItemsOpts{
		Limit:     limit,
		Timeframe: tf,
		Page:      page,
		Week:      week,
		Month:     month,
		Year:      year,
		From:      from,
		To:        to,
		ArtistID:  artistId,
		AlbumID:   albumId,
		TrackID:   trackId,
	}, nil
}

func isDateRangeValidationError(err error) bool {
	var dateRangeErr *utils.DateRangeValidationError
	return errors.As(err, &dateRangeErr)
}

type rankedResponseItem[T any] struct {
	Item         T     `json:"item"`
	Rank         int64 `json:"rank"`
	ListenCount  int64 `json:"listen_count"`
	TimeListened int64 `json:"time_listened"`
}

type rankedPaginatedResponse[T any] struct {
	Items        []rankedResponseItem[T] `json:"items"`
	TotalCount   int64                   `json:"total_record_count"`
	ItemsPerPage int32                   `json:"items_per_page"`
	HasNextPage  bool                    `json:"has_next_page"`
	CurrentPage  int32                   `json:"current_page"`
}

func rankedPaginatedResponseFrom[T any](response *db.PaginatedResponse[db.RankedItem[T]]) *rankedPaginatedResponse[T] {
	if response == nil {
		return nil
	}

	items := make([]rankedResponseItem[T], len(response.Items))
	for i, item := range response.Items {
		items[i] = rankedResponseItem[T]{
			Item:         item.Item,
			Rank:         item.Rank,
			ListenCount:  item.ListenCount,
			TimeListened: item.TimeListened,
		}
	}

	return &rankedPaginatedResponse[T]{
		Items:        items,
		TotalCount:   response.TotalCount,
		ItemsPerPage: response.ItemsPerPage,
		HasNextPage:  response.HasNextPage,
		CurrentPage:  response.CurrentPage,
	}
}

func TimeframeFromRequest(r *http.Request) db.Timeframe {
	q := r.URL.Query()

	parseInt := func(key string) int {
		v := q.Get(key)
		if v == "" {
			return 0
		}
		i, _ := strconv.Atoi(v)
		return i
	}

	parseInt64 := func(key string) int64 {
		v := q.Get(key)
		if v == "" {
			return 0
		}
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	}

	return db.Timeframe{
		Period:   db.Period(q.Get("period")),
		Year:     parseInt("year"),
		Month:    parseInt("month"),
		Week:     parseInt("week"),
		FromUnix: parseInt64("from"),
		ToUnix:   parseInt64("to"),
		Timezone: parseTZ(r),
	}
}

func parseTZ(r *http.Request) *time.Location {

	// this map is obviously AI.
	// i manually referenced as many links as I could and couldn't find any
	// incorrect entries here so hopefully it is all correct.
	overrides := map[string]string{
		// --- North America ---
		"America/Indianapolis":  "America/Indiana/Indianapolis",
		"America/Knoxville":     "America/Indiana/Knoxville",
		"America/Louisville":    "America/Kentucky/Louisville",
		"America/Montreal":      "America/Toronto",
		"America/Shiprock":      "America/Denver",
		"America/Fort_Wayne":    "America/Indiana/Indianapolis",
		"America/Virgin":        "America/Port_of_Spain",
		"America/Santa_Isabel":  "America/Tijuana",
		"America/Ensenada":      "America/Tijuana",
		"America/Rosario":       "America/Argentina/Cordoba",
		"America/Jujuy":         "America/Argentina/Jujuy",
		"America/Mendoza":       "America/Argentina/Mendoza",
		"America/Catamarca":     "America/Argentina/Catamarca",
		"America/Cordoba":       "America/Argentina/Cordoba",
		"America/Buenos_Aires":  "America/Argentina/Buenos_Aires",
		"America/Coral_Harbour": "America/Atikokan",
		"America/Atka":          "America/Adak",
		"US/Alaska":             "America/Anchorage",
		"US/Aleutian":           "America/Adak",
		"US/Arizona":            "America/Phoenix",
		"US/Central":            "America/Chicago",
		"US/Eastern":            "America/New_York",
		"US/East-Indiana":       "America/Indiana/Indianapolis",
		"US/Hawaii":             "Pacific/Honolulu",
		"US/Indiana-Starke":     "America/Indiana/Knoxville",
		"US/Michigan":           "America/Detroit",
		"US/Mountain":           "America/Denver",
		"US/Pacific":            "America/Los_Angeles",
		"US/Samoa":              "Pacific/Pago_Pago",
		"Canada/Atlantic":       "America/Halifax",
		"Canada/Central":        "America/Winnipeg",
		"Canada/Eastern":        "America/Toronto",
		"Canada/Mountain":       "America/Edmonton",
		"Canada/Newfoundland":   "America/St_Johns",
		"Canada/Pacific":        "America/Vancouver",

		// --- Asia ---
		"Asia/Calcutta":      "Asia/Kolkata",
		"Asia/Saigon":        "Asia/Ho_Chi_Minh",
		"Asia/Katmandu":      "Asia/Kathmandu",
		"Asia/Rangoon":       "Asia/Yangon",
		"Asia/Ulan_Bator":    "Asia/Ulaanbaatar",
		"Asia/Macao":         "Asia/Macau",
		"Asia/Tel_Aviv":      "Asia/Jerusalem",
		"Asia/Ashkhabad":     "Asia/Ashgabat",
		"Asia/Chungking":     "Asia/Chongqing",
		"Asia/Dacca":         "Asia/Dhaka",
		"Asia/Istanbul":      "Europe/Istanbul",
		"Asia/Kashgar":       "Asia/Urumqi",
		"Asia/Thimbu":        "Asia/Thimphu",
		"Asia/Ujung_Pandang": "Asia/Makassar",
		"ROC":                "Asia/Taipei",
		"Iran":               "Asia/Tehran",
		"Israel":             "Asia/Jerusalem",
		"Japan":              "Asia/Tokyo",
		"Singapore":          "Asia/Singapore",
		"Hongkong":           "Asia/Hong_Kong",

		// --- Europe ---
		"Europe/Kiev":     "Europe/Kyiv",
		"Europe/Belfast":  "Europe/London",
		"Europe/Tiraspol": "Europe/Chisinau",
		"Europe/Nicosia":  "Asia/Nicosia",
		"Europe/Moscow":   "Europe/Moscow",
		"W-SU":            "Europe/Moscow",
		"GB":              "Europe/London",
		"GB-Eire":         "Europe/London",
		"Eire":            "Europe/Dublin",
		"Poland":          "Europe/Warsaw",
		"Portugal":        "Europe/Lisbon",
		"Turkey":          "Europe/Istanbul",

		// --- Australia / Pacific ---
		"Australia/ACT":        "Australia/Sydney",
		"Australia/Canberra":   "Australia/Sydney",
		"Australia/LHI":        "Australia/Lord_Howe",
		"Australia/North":      "Australia/Darwin",
		"Australia/NSW":        "Australia/Sydney",
		"Australia/Queensland": "Australia/Brisbane",
		"Australia/South":      "Australia/Adelaide",
		"Australia/Tasmania":   "Australia/Hobart",
		"Australia/Victoria":   "Australia/Melbourne",
		"Australia/West":       "Australia/Perth",
		"Australia/Yancowinna": "Australia/Broken_Hill",
		"Pacific/Samoa":        "Pacific/Pago_Pago",
		"Pacific/Yap":          "Pacific/Chuuk",
		"Pacific/Truk":         "Pacific/Chuuk",
		"Pacific/Ponape":       "Pacific/Pohnpei",
		"NZ":                   "Pacific/Auckland",
		"NZ-CHAT":              "Pacific/Chatham",

		// --- Africa ---
		"Africa/Asmera":   "Africa/Asmara",
		"Africa/Timbuktu": "Africa/Bamako",
		"Egypt":           "Africa/Cairo",
		"Libya":           "Africa/Tripoli",

		// --- Atlantic ---
		"Atlantic/Faeroe":    "Atlantic/Faroe",
		"Atlantic/Jan_Mayen": "Europe/Oslo",
		"Iceland":            "Atlantic/Reykjavik",

		// --- Etc / Misc ---
		"UTC":       "UTC",
		"Etc/UTC":   "UTC",
		"Etc/GMT":   "UTC",
		"GMT":       "UTC",
		"Zulu":      "UTC",
		"Universal": "UTC",
	}

	if cfg.ForceTZ() != nil {
		return cfg.ForceTZ()
	}

	if tz := r.URL.Query().Get("tz"); tz != "" {
		if fixedTz, exists := overrides[tz]; exists {
			tz = fixedTz
		}
		if loc, err := time.LoadLocation(tz); err == nil {
			return loc
		}
	}

	if c, err := r.Cookie("tz"); err == nil {
		var tz string
		if fixedTz, exists := overrides[c.Value]; exists {
			tz = fixedTz
		} else {
			tz = c.Value
		}
		if loc, err := time.LoadLocation(tz); err == nil {
			return loc
		}
	}

	return time.Now().Location()
}
