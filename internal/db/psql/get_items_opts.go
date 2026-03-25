package psql

import (
	"errors"

	"github.com/gabehf/koito/internal/db"
)

func normalizePagedGetItemsOpts(opts db.GetItemsOpts) (db.GetItemsOpts, error) {
	if opts.Month != 0 && opts.Year == 0 {
		return opts, errors.New("year must be specified with month")
	}
	if opts.Limit < 0 {
		return opts, errors.New("limit must be greater than or equal to 0")
	}
	if opts.Page < 0 {
		return opts, errors.New("page must be greater than or equal to 0")
	}
	if opts.Limit == 0 {
		opts.Limit = DefaultItemsPerPage
	}
	if opts.Page == 0 {
		opts.Page = 1
	}
	if hasLegacyItemDateFilters(opts) {
		opts.Timeframe = db.Timeframe{
			Year:     opts.Year,
			Month:    opts.Month,
			Week:     opts.Week,
			FromUnix: int64(opts.From),
			ToUnix:   int64(opts.To),
		}
	}
	if opts.Timeframe == (db.Timeframe{}) {
		opts.Timeframe = db.PeriodToTimeframe(db.PeriodDay)
	}

	return opts, nil
}

func hasLegacyItemDateFilters(opts db.GetItemsOpts) bool {
	return opts.Year != 0 || opts.Month != 0 || opts.Week != 0 || opts.From != 0 || opts.To != 0
}
