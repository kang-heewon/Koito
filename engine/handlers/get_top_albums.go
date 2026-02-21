package handlers

import (
	"net/http"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

func GetTopAlbumsHandler(store db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msg("GetTopAlbumsHandler: Received request to retrieve top albums")

		opts, err := OptsFromRequest(r)
		if err != nil {
			l.Err(err).Msg("GetTopAlbumsHandler: Invalid query parameters")
			utils.WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}

		l.Debug().Msgf("GetTopAlbumsHandler: Retrieving top albums with options: %+v", opts)

		albums, err := store.GetTopAlbumsPaginated(ctx, opts)
		if err != nil {
			l.Err(err).Msg("GetTopAlbumsHandler: Failed to retrieve top albums")
			if isDateRangeValidationError(err) {
				utils.WriteError(w, err.Error(), http.StatusBadRequest)
				return
			}

			utils.WriteError(w, "failed to get albums", http.StatusInternalServerError)
			return
		}

		l.Debug().Msg("GetTopAlbumsHandler: Successfully retrieved top albums")
		utils.WriteJSON(w, http.StatusOK, albums)
	}
}
