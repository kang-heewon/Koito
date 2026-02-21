package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/utils"
	"github.com/jackc/pgx/v5"
)

type SearchResults struct {
	Artists []*models.Artist `json:"artists"`
	Albums  []*models.Album  `json:"albums"`
	Tracks  []*models.Track  `json:"tracks"`
}

func SearchHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)
		q := strings.TrimSpace(r.URL.Query().Get("q"))

		l.Debug().Msgf("SearchHandler: Received search with query: %s", r.URL.Query().Encode())

		var artists []*models.Artist
		var albums []*models.Album
		var tracks []*models.Track

		qLower := strings.ToLower(q)
		if strings.HasPrefix(qLower, "id:") {
			idStr := strings.TrimSpace(q[3:])
			id, err := strconv.Atoi(idStr)
			if err != nil || id < 1 {
				l.Debug().Msgf("SearchHandler: Invalid id query %q", q)
				utils.WriteError(w, "id search must use a positive integer (id:<number>)", http.StatusBadRequest)
				return
			}

			artist, err := store.GetArtist(ctx, db.GetArtistOpts{ID: int32(id)})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					l.Debug().Msgf("SearchHandler: No artists found with id=%d", id)
				} else {
					l.Err(err).Msgf("SearchHandler: Failed to search artist by id=%d", id)
					utils.WriteError(w, "failed to search in database", http.StatusInternalServerError)
					return
				}
			}
			if artist != nil {
				artists = append(artists, artist)
			}

			album, err := store.GetAlbum(ctx, db.GetAlbumOpts{ID: int32(id)})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					l.Debug().Msgf("SearchHandler: No albums found with id=%d", id)
				} else {
					l.Err(err).Msgf("SearchHandler: Failed to search album by id=%d", id)
					utils.WriteError(w, "failed to search in database", http.StatusInternalServerError)
					return
				}
			}
			if album != nil {
				albums = append(albums, album)
			}

			track, err := store.GetTrack(ctx, db.GetTrackOpts{ID: int32(id)})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					l.Debug().Msgf("SearchHandler: No tracks found with id=%d", id)
				} else {
					l.Err(err).Msgf("SearchHandler: Failed to search track by id=%d", id)
					utils.WriteError(w, "failed to search in database", http.StatusInternalServerError)
					return
				}
			}
			if track != nil {
				tracks = append(tracks, track)
			}
		} else {
			var err error
			artists, err = store.SearchArtists(ctx, q)
			if err != nil {
				l.Err(err).Msg("Failed to search for artists")
				utils.WriteError(w, "failed to search in database", http.StatusInternalServerError)
				return
			}
			albums, err = store.SearchAlbums(ctx, q)
			if err != nil {
				l.Err(err).Msg("Failed to search for albums")
				utils.WriteError(w, "failed to search in database", http.StatusInternalServerError)
				return
			}
			tracks, err = store.SearchTracks(ctx, q)
			if err != nil {
				l.Err(err).Msg("Failed to search for tracks")
				utils.WriteError(w, "failed to search in database", http.StatusInternalServerError)
				return
			}
		}

		utils.WriteJSON(w, http.StatusOK, SearchResults{
			Artists: artists,
			Albums:  albums,
			Tracks:  tracks,
		})
	}
}
