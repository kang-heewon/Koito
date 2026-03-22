package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/gabehf/koito/internal/catalog"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/mbz"
	"github.com/gabehf/koito/internal/utils"
)

var (
	backfillRunning = int32(0)
	backfillCancel  context.CancelFunc
)

func BackfillGenresHandler(store db.DB, mbzC mbz.MusicBrainzCaller, discogsC catalog.DiscogsCaller, lastfmC catalog.LastFmCaller, spotifyC catalog.SpotifyCaller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		if r.Method != http.MethodPost {
			l.Warn().Msg("BackfillGenresHandler: Method not allowed")
			utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		l.Info().Msg("BackfillGenresHandler: Received manual backfill request")

		if !atomic.CompareAndSwapInt32(&backfillRunning, 0, 1) {
			l.Warn().Msg("BackfillGenresHandler: Backfill already running")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "backfill already running",
			})
			return
		}

		backfillCtx, cancel := context.WithCancel(context.Background())
		backfillCancel = cancel

		l.Info().Msg("BackfillGenresHandler: Starting backfill in goroutine")

		go func() {
			defer atomic.StoreInt32(&backfillRunning, 0)
			defer cancel()
			fetcher := catalog.NewHybridGenreFetcher(mbzC, discogsC, lastfmC, spotifyC)
			catalog.BackfillGenres(backfillCtx, store, fetcher)
		}()

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "backfill started",
		})
	}
}
