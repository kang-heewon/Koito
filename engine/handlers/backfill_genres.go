package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gabehf/koito/internal/catalog"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/mbz"
	"github.com/gabehf/koito/internal/utils"
)

type BackfillController struct {
	running int32
	cancel  context.CancelFunc
	mu      sync.Mutex
	appCtx  context.Context
}

func NewBackfillController(appCtx context.Context) *BackfillController {
	return &BackfillController{
		running: 0,
		appCtx:  appCtx,
	}
}

func (c *BackfillController) Begin() (ctx context.Context, release func(), ok bool) {
	if !atomic.CompareAndSwapInt32(&c.running, 0, 1) {
		return nil, nil, false
	}

	c.mu.Lock()
	backfillCtx, cancel := context.WithCancel(c.appCtx)
	c.cancel = cancel
	c.mu.Unlock()

	release = func() {
		atomic.StoreInt32(&c.running, 0)
		cancel()
	}

	return backfillCtx, release, true
}

func (c *BackfillController) Cancel() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cancel != nil {
		c.cancel()
	}
}

func BackfillGenresHandler(store db.DB, mbzC mbz.MusicBrainzCaller, discogsC catalog.DiscogsCaller, lastfmC catalog.LastFmCaller, spotifyC catalog.SpotifyCaller, controller *BackfillController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logger.FromContext(ctx)

		if r.Method != http.MethodPost {
			l.Warn().Msg("BackfillGenresHandler: Method not allowed")
			utils.WriteError(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		l.Info().Msg("BackfillGenresHandler: Received manual backfill request")

		backfillCtx, release, ok := controller.Begin()
		if !ok {
			l.Warn().Msg("BackfillGenresHandler: Backfill already running")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "backfill already running",
			})
			return
		}

		l.Info().Msg("BackfillGenresHandler: Starting backfill in goroutine")

		go func() {
			defer release()
			fetcher := catalog.NewHybridGenreFetcher(mbzC, discogsC, lastfmC, spotifyC)
			catalog.BackfillGenres(backfillCtx, store, fetcher)
		}()

		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "backfill started",
		})
	}
}
