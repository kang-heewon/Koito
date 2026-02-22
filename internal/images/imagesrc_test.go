package images

import (
	"context"
	"testing"
)

func TestGetArtistImage_EmptyAliasesPanic(t *testing.T) {
	ctx := context.Background()
	Initialize(ImageSourceOpts{})

	tests := []struct {
		name    string
		opts    ArtistImageOpts
		wantErr bool
	}{
		{
			name: "empty Aliases slice should not panic",
			opts: ArtistImageOpts{
				Aliases: []string{},
			},
			wantErr: false, // Should return empty string, not panic
		},
		{
			name: "nil Aliases slice should not panic",
			opts: ArtistImageOpts{
				Aliases: nil,
			},
			wantErr: false, // Should return empty string, not panic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Before fix: this will panic with "index out of range" at line 92
			// After fix: this should handle gracefully
			result, err := GetArtistImage(ctx, tt.opts)

			// Should not panic - if we reach here, bounds check worked
			if tt.wantErr && err == nil {
				t.Errorf("GetArtistImage() expected error, got nil")
			}
			if !tt.wantErr && result != "" {
				t.Errorf("GetArtistImage() expected empty result, got %s", result)
			}
		})
	}
}

func TestGetAlbumImage_EmptyArtistsPanic(t *testing.T) {
	ctx := context.Background()
	Initialize(ImageSourceOpts{})

	tests := []struct {
		name    string
		opts    AlbumImageOpts
		wantErr bool
	}{
		{
			name: "empty Artists slice should not panic",
			opts: AlbumImageOpts{
				Artists: []string{},
				Album:   "Test Album",
			},
			wantErr: false, // Should return empty string, not panic
		},
		{
			name: "nil Artists slice should not panic",
			opts: AlbumImageOpts{
				Artists: nil,
				Album:   "Test Album",
			},
			wantErr: false, // Should return empty string, not panic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Before fix: this will panic with "index out of range" at line 123
			// After fix: this should handle gracefully
			result, err := GetAlbumImage(ctx, tt.opts)

			// Should not panic - if we reach here, bounds check worked
			if tt.wantErr && err == nil {
				t.Errorf("GetAlbumImage() expected error, got nil")
			}
			if !tt.wantErr && result != "" {
				t.Errorf("GetAlbumImage() expected empty result, got %s", result)
			}
		})
	}
}
