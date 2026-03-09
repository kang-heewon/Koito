package db

import (
	"strconv"
	"time"
)

// should this be in db package ???

// Deprecated: Use Timeframe instead.
type Period string

const (
	PeriodDay     Period = "day"
	PeriodWeek    Period = "week"
	PeriodMonth   Period = "month"
	PeriodYear    Period = "year"
	PeriodAllTime Period = "all_time"
	PeriodDefault Period = "day"
)

type Timeframe struct {
	Period   string
	Year     string
	Month    string
	Week     string
	FromUnix string
	ToUnix   string
	From     string
	To       string
	Timezone string
}

func PeriodToTimeframe(p Period) Timeframe {
	return Timeframe{Period: string(p)}
}

func TimeframeToTimeRange(tf Timeframe) (time.Time, time.Time) {
	now := time.Now()
	loc := now.Location()
	if tf.Timezone != "" {
		if timezone, err := time.LoadLocation(tf.Timezone); err == nil {
			loc = timezone
		}
	}

	if from, ok := parseTimeframeTime(tf.From, loc); ok {
		if to, ok := parseTimeframeTime(tf.To, loc); ok {
			return from, to
		}
		return from, now
	}

	if fromUnix, ok := parseTimeframeUnix(tf.FromUnix, loc); ok {
		if toUnix, ok := parseTimeframeUnix(tf.ToUnix, loc); ok {
			return fromUnix, toUnix
		}
		return fromUnix, now
	}

	if year, ok := parseTimeframeInt(tf.Year); ok {
		if month, ok := parseTimeframeInt(tf.Month); ok {
			start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
			end := endOfMonth(year, time.Month(month), loc)
			return start, end
		}

		if week, ok := parseTimeframeInt(tf.Week); ok {
			jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, loc)
			week1Start := startOfWeek(jan4)
			start := week1Start.AddDate(0, 0, (week-1)*7)
			return start, endOfWeek(start)
		}

		start := time.Date(year, 1, 1, 0, 0, 0, 0, loc)
		end := time.Date(year+1, 1, 1, 0, 0, 0, 0, loc).Add(-time.Second)
		return start, end
	}

	if month, ok := parseTimeframeInt(tf.Month); ok {
		year := now.In(loc).Year()
		if int(now.In(loc).Month()) < month {
			year--
		}

		start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
		end := endOfMonth(year, time.Month(month), loc)
		return start, end
	}

	if week, ok := parseTimeframeInt(tf.Week); ok {
		current := now.In(loc)
		year := current.Year()
		_, currentWeek := current.ISOWeek()
		if currentWeek < week {
			year--
		}

		jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, loc)
		week1Start := startOfWeek(jan4)
		start := week1Start.AddDate(0, 0, (week-1)*7)
		return start, endOfWeek(start)
	}

	switch tf.Period {
	case "all", string(PeriodAllTime):
		return time.Time{}, now
	case "":
		return time.Time{}, time.Time{}
	default:
		return StartTimeFromPeriod(Period(tf.Period)), now
	}
}

func parseTimeframeInt(value string) (int, bool) {
	if value == "" {
		return 0, false
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}

	return parsed, true
}

func parseTimeframeUnix(value string, loc *time.Location) (time.Time, bool) {
	if value == "" {
		return time.Time{}, false
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}, false
	}

	return time.Unix(parsed, 0).In(loc), true
}

func parseTimeframeTime(value string, loc *time.Location) (time.Time, bool) {
	if value == "" {
		return time.Time{}, false
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		var (
			parsed time.Time
			err    error
		)

		if layout == "2006-01-02 15:04:05" || layout == "2006-01-02" {
			parsed, err = time.ParseInLocation(layout, value, loc)
		} else {
			parsed, err = time.Parse(layout, value)
		}

		if err == nil {
			return parsed.In(loc), true
		}
	}

	return time.Time{}, false
}

// Deprecated: Use TimeframeToTimeRange instead.
func StartTimeFromPeriod(p Period) time.Time {
	now := time.Now()
	switch p {
	case "day":
		return now.AddDate(0, 0, -1)
	case "week":
		return now.AddDate(0, 0, -7)
	case "month":
		return now.AddDate(0, -1, 0)
	case "year":
		return now.AddDate(-1, 0, 0)
	case "all_time":
		return time.Time{}
	default:
		// default 1 day
		return now.AddDate(0, 0, -1)
	}
}

func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	return time.Date(t.Year(), t.Month(), t.Day()-weekday+1, 0, 0, 0, 0, t.Location())
}

func endOfWeek(t time.Time) time.Time {
	return startOfWeek(t).AddDate(0, 0, 7).Add(-time.Second)
}

func endOfMonth(year int, month time.Month, loc *time.Location) time.Time {
	return time.Date(year, month+1, 1, 0, 0, 0, 0, loc).Add(-time.Second)
}

type StepInterval string

const (
	StepDay     StepInterval = "day"
	StepWeek    StepInterval = "week"
	StepMonth   StepInterval = "month"
	StepYear    StepInterval = "year"
	StepDefault StepInterval = "day"

	DefaultRange int = 12
)

// start is the time of 00:00 at the beginning of opts.Range opts.Steps ago,
// end is the end time of the current opts.Step.
// E.g. if step is StepWeek and range is 4, start will be the time 00:00 on Sunday on the 4th week ago,
// and end will be 23:59:59 on Saturday at the end of the current week.
// If opts.Year (or opts.Year + opts.Month) is provided, start and end will simply by the start and end times of that year/month.
func ListenActivityOptsToTimes(opts ListenActivityOpts) (start, end time.Time) {
	now := time.Now()

	// If Year (and optionally Month) are specified, use calendar boundaries
	if opts.Year != 0 {
		if opts.Month != 0 {
			// Specific month of a specific year
			start = time.Date(opts.Year, time.Month(opts.Month), 1, 0, 0, 0, 0, now.Location())
			end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		} else {
			// Whole year
			start = time.Date(opts.Year, 1, 1, 0, 0, 0, 0, now.Location())
			end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)
		}
		return start, end
	}

	// X days ago + today = range
	opts.Range = opts.Range - 1

	// Determine step and align accordingly
	switch opts.Step {
	case StepDay:
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		start = today.AddDate(0, 0, -opts.Range)
		end = today.AddDate(0, 0, 1).Add(-time.Nanosecond)

	case StepWeek:
		// Align to most recent Sunday
		weekday := int(now.Weekday()) // Sunday = 0
		startOfThisWeek := time.Date(now.Year(), now.Month(), now.Day()-weekday, 0, 0, 0, 0, now.Location())
		start = startOfThisWeek.AddDate(0, 0, -7*opts.Range)
		end = startOfThisWeek.AddDate(0, 0, 7).Add(-time.Nanosecond)

	case StepMonth:
		firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		start = firstOfThisMonth.AddDate(0, -opts.Range, 0)
		end = firstOfThisMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	case StepYear:
		firstOfThisYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		start = firstOfThisYear.AddDate(-opts.Range, 0, 0)
		end = firstOfThisYear.AddDate(1, 0, 0).Add(-time.Nanosecond)

	default:
		// Default to daily
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		start = today.AddDate(0, 0, -opts.Range)
		end = today.AddDate(0, 0, 1).Add(-time.Nanosecond)
	}

	return start, end
}
