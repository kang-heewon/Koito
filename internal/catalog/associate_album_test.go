package catalog

import (
	"context"
	"testing"

	"github.com/gabehf/koito/internal/models"
	"github.com/google/uuid"
)

func TestAssociateAlbum_EmptyArtistsPanic(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		opts    AssociateAlbumOpts
		wantErr bool
	}{
		{
			name: "empty Artists slice should return error",
			opts: AssociateAlbumOpts{
				Artists:      []*models.Artist{},
				TrackName:    "Test Track",
				ReleaseMbzID: uuid.Nil,
			},
			wantErr: true,
		},
		{
			name: "nil Artists slice should return error",
			opts: AssociateAlbumOpts{
				Artists:      nil,
				TrackName:    "Test Track",
				ReleaseMbzID: uuid.Nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Before fix: this will panic with "index out of range"
			// After fix: this should return an error
			_, err := AssociateAlbum(ctx, nil, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("AssociateAlbum() error = %v, wantErr %v", err, tt.wantErr)
			}

			// If panic occurred, test will fail before reaching here
		})
	}
}
