package handlers

import (
	"net/http"

	"github.com/gabehf/koito/engine/middleware"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/summary"
	"github.com/gabehf/koito/internal/utils"
)

func SummaryHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msg("SummaryHandler: Received request to retrieve summary")

		userID := int32(1)
		if u := middleware.GetUserFromContext(ctx); u != nil {
			userID = u.ID
		}

		timeframe := TimeframeFromRequest(r)
		summaryData, err := summary.GenerateSummary(ctx, store, userID, timeframe, "")
		if err != nil {
			l.Err(err).Int32("user_id", userID).Any("timeframe", timeframe).Msg("SummaryHandler: Failed to generate summary")
			utils.WriteError(w, "failed to generate summary", http.StatusInternalServerError)
			return
		}

		utils.WriteJSON(w, http.StatusOK, summaryData)
	}
}
