package handlers

import (
	"net/http"
	"time"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
)

type Artist struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type RecommendationItem struct {
	ID              int32    `json:"id"`
	Title           string   `json:"title"`
	Artists         []Artist `json:"artists"`
	AlbumID         int32    `json:"album_id,omitempty"`
	Image           *string  `json:"image,omitempty"`
	PastListenCount int64    `json:"past_listen_count"`
	LastListenedAt  string   `json:"last_listened_at"`
}

type RecommendationsResponse struct {
	Tracks []RecommendationItem `json:"tracks"`
}

func RecommendationsHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		now := time.Now()
		opts := db.GetRecommendationsOpts{
			PastWindowStart: now.AddDate(0, 0, -90),
			PastWindowEnd:   now.AddDate(0, 0, -30),
			MinPastListens:  5,
			Limit:           20,
		}

		recommendations, err := store.GetTracksToRevisit(r.Context(), opts)
		if err != nil {
			l.Err(err).Msg("RecommendationsHandler: Failed to fetch recommendations")
			utils.WriteError(w, "failed to get recommendations: "+err.Error(), http.StatusInternalServerError)
			return
		}

		items := make([]RecommendationItem, len(recommendations))
		for i, rec := range recommendations {
			item := RecommendationItem{
				ID:              rec.Track.ID,
				Title:           rec.Track.Title,
				PastListenCount: rec.PastListenCount,
				LastListenedAt:  rec.LastListenedAt.Format(time.RFC3339),
			}

			if rec.Track.Artists != nil {
				artists := make([]Artist, len(rec.Track.Artists))
				for j, a := range rec.Track.Artists {
					artists[j] = Artist{ID: a.ID, Name: a.Name}
				}
				item.Artists = artists
			}

			if rec.Track.AlbumID != 0 {
				item.AlbumID = rec.Track.AlbumID
			}

			if rec.Track.Image != nil {
				img := rec.Track.Image.String()
				item.Image = &img
			}

			items[i] = item
		}

		utils.WriteJSON(w, http.StatusOK, RecommendationsResponse{Tracks: items})
	}
}
