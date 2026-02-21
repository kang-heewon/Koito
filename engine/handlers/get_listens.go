package handlers

import (
	"net/http"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

func GetListensHandler(store db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msg("GetListensHandler: Received request to retrieve listens")

		opts, err := OptsFromRequest(r)
		if err != nil {
			l.Err(err).Msg("GetListensHandler: Invalid query parameters")
			utils.WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}

		l.Debug().Msgf("GetListensHandler: Retrieving listens with options: %+v", opts)

		listens, err := store.GetListensPaginated(ctx, opts)
		if err != nil {
			l.Err(err).Msg("GetListensHandler: Failed to retrieve listens")
			if isDateRangeValidationError(err) {
				utils.WriteError(w, err.Error(), http.StatusBadRequest)
				return
			}

			utils.WriteError(w, "failed to get listens", http.StatusInternalServerError)
			return
		}

		l.Debug().Msg("GetListensHandler: Successfully retrieved listens")
		utils.WriteJSON(w, http.StatusOK, listens)
	}
}
