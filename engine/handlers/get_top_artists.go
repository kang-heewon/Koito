package handlers

import (
	"net/http"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

func GetTopArtistsHandler(store db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msg("GetTopArtistsHandler: Received request to retrieve top artists")

		opts, err := OptsFromRequest(r)
		if err != nil {
			l.Err(err).Msg("GetTopArtistsHandler: Invalid query parameters")
			utils.WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}

		l.Debug().Msgf("GetTopArtistsHandler: Retrieving top artists with options: %+v", opts)

		artists, err := store.GetTopArtistsPaginated(ctx, opts)
		if err != nil {
			l.Err(err).Msg("GetTopArtistsHandler: Failed to retrieve top artists")
			if isDateRangeValidationError(err) {
				utils.WriteError(w, err.Error(), http.StatusBadRequest)
				return
			}

			utils.WriteError(w, "failed to get artists", http.StatusInternalServerError)
			return
		}

		l.Debug().Msg("GetTopArtistsHandler: Successfully retrieved top artists")
		utils.WriteJSON(w, http.StatusOK, artists)
	}
}
