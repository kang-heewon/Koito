package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

// MergeHandler creates a handler for merge operations.
// name: handler name for logging
// mergeFn: function that performs the actual merge operation
// hasReplaceImage: whether the merge operation supports the replace_image parameter
func MergeHandler(
	name string,
	mergeFn func(ctx context.Context, fromId, toId int32, replaceImage bool) error,
	hasReplaceImage bool,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msgf("%s: Received request", name)

		// Parse from_id parameter
		fromidStr := r.URL.Query().Get("from_id")
		fromId, err := strconv.Atoi(fromidStr)
		if err != nil {
			l.Debug().AnErr("error", err).Msgf("%s: Invalid from_id parameter", name)
			utils.WriteError(w, "from_id is invalid", http.StatusBadRequest)
			return
		}

		// Parse to_id parameter
		toidStr := r.URL.Query().Get("to_id")
		toId, err := strconv.Atoi(toidStr)
		if err != nil {
			l.Debug().AnErr("error", err).Msgf("%s: Invalid to_id parameter", name)
			utils.WriteError(w, "to_id is invalid", http.StatusBadRequest)
			return
		}

		// Parse replace_image parameter (optional)
		var replaceImage bool
		if hasReplaceImage {
			replaceImgStr := r.URL.Query().Get("replace_image")
			if strings.ToLower(replaceImgStr) == "true" {
				l.Debug().Msgf("%s: Merge will replace image", name)
				replaceImage = true
			}
		}

		l.Debug().Msgf("%s: Merging from ID %d to ID %d", name, fromId, toId)

		// Execute merge function
		err = mergeFn(ctx, int32(fromId), int32(toId), replaceImage)
		if err != nil {
			l.Err(err).Msgf("%s: Failed to merge", name)
			utils.WriteError(w, name+" failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		l.Debug().Msgf("%s: Successfully merged from ID %d to ID %d", name, fromId, toId)
		w.WriteHeader(http.StatusNoContent)
	}
}

func MergeTracksHandler(store db.DB) http.HandlerFunc {
	// Adapter to make MergeTracks compatible with the factory signature
	adapter := func(ctx context.Context, fromId, toId int32, replaceImage bool) error {
		return store.MergeTracks(ctx, fromId, toId)
	}
	return MergeHandler("MergeTracksHandler", adapter, false)
}

func MergeReleaseGroupsHandler(store db.DB) http.HandlerFunc {
	return MergeHandler("MergeReleaseGroupsHandler", store.MergeAlbums, true)
}

func MergeArtistsHandler(store db.DB) http.HandlerFunc {
	return MergeHandler("MergeArtistsHandler", store.MergeArtists, true)
}

func UpdateAlbumHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		l.Debug().Msg("UpdateAlbumHandler: Received request")

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)

		valStr := r.URL.Query().Get("is_various_artists")
		var variousArists bool
		var updateVariousArtists = false
		if strings.ToLower(valStr) == "true" {
			variousArists = true
			updateVariousArtists = true
		} else if strings.ToLower(valStr) == "false" {
			variousArists = false
			updateVariousArtists = true
		}
		if err != nil {
			l.Debug().AnErr("error", err).Msg("UpdateAlbumHandler: Invalid id parameter")
			utils.WriteError(w, "id is invalid", http.StatusBadRequest)
			return
		}

		err = store.UpdateAlbum(ctx, db.UpdateAlbumOpts{
			ID:                   int32(id),
			VariousArtistsUpdate: updateVariousArtists,
			VariousArtistsValue:  variousArists,
		})
		if err != nil {
			l.Debug().AnErr("error", err).Msg("UpdateAlbumHandler: Failed to update album")
			utils.WriteError(w, "failed to update album", http.StatusBadRequest)
			return
		}

		l.Debug().Msg("UpdateAlbumHandler: Successfully updated album")

		w.WriteHeader(http.StatusNoContent)
	}
}
