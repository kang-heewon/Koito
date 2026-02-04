package handlers

import (
	"net/http"
	"strings"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

type GenreStatsResponse struct {
	Stats []GenreStatItem `json:"stats"`
}

type GenreStatItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func GenreStatsHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

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
			period = db.PeriodMonth
		}

		metric := strings.ToLower(r.URL.Query().Get("metric"))

		var stats []db.GenreStat
		var err error

		if metric == "time" {
			stats, err = store.GetGenreStatsByTimeListened(r.Context(), period)
		} else {
			stats, err = store.GetGenreStatsByListenCount(r.Context(), period)
		}

		if err != nil {
			l.Err(err).Msg("GenreStatsHandler: Failed to fetch genre stats")
			utils.WriteError(w, "failed to get genre stats: "+err.Error(), http.StatusInternalServerError)
			return
		}

		items := make([]GenreStatItem, len(stats))
		for i, s := range stats {
			items[i] = GenreStatItem{Name: s.Name, Value: s.Value}
		}

		utils.WriteJSON(w, http.StatusOK, GenreStatsResponse{Stats: items})
	}
}
