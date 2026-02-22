// package imagesrc defines interfaces for album and artist image providers
package images

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gabehf/koito/internal/cover"
	"github.com/gabehf/koito/internal/logger"
	"github.com/google/uuid"
)

type ImageSource struct {
	deezerEnabled   bool
	deezerC         *DeezerClient
	subsonicEnabled bool
	subsonicC       *SubsonicClient
	spotifyEnabled  bool
	spotifyC        *SpotifyClient
	caaEnabled      bool
}
type ImageSourceOpts struct {
	UserAgent      string
	EnableCAA      bool
	EnableDeezer   bool
	EnableSubsonic bool
	EnableSpotify  bool
}

var once sync.Once
var imgsrc ImageSource

type ArtistImageOpts struct {
	Aliases []string
}

type AlbumImageOpts struct {
	Artists           []string
	Album             string
	ReleaseMbzID      *uuid.UUID
	ReleaseGroupMbzID *uuid.UUID
}

const caaBaseUrl = "https://coverartarchive.org"

var caaClient = &http.Client{Timeout: 15 * time.Second}

type caaResponse struct {
	Images []caaImage `json:"images"`
}

type caaImage struct {
	Image      string            `json:"image"`
	Front      bool              `json:"front"`
	Back       bool              `json:"back"`
	Width      int               `json:"width"`
	Height     int               `json:"height"`
	Thumbnails map[string]string `json:"thumbnails"`
}

func caaCoverImageExtract(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("caaCoverImageExtract: %w", err)
	}

	resp, err := caaClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("caaCoverImageExtract: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	parsed := new(caaResponse)
	if err := json.NewDecoder(resp.Body).Decode(parsed); err != nil {
		return "", fmt.Errorf("caaCoverImageExtract: %w", err)
	}

	images := make([]cover.Image, 0, len(parsed.Images))
	for _, image := range parsed.Images {
		images = append(images, cover.Image{
			URL:        image.Image,
			Front:      image.Front,
			Back:       image.Back,
			Width:      image.Width,
			Height:     image.Height,
			Thumbnails: image.Thumbnails,
		})
	}

	return cover.CoverImageExtract(images), nil
}

func caaFrontImage(url string) (string, string, error) {
	resp, err := caaClient.Head(url)
	if err != nil {
		return "", "", fmt.Errorf("caaFrontImage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return url, resp.Status, nil
	}

	return "", resp.Status, nil
}

// all functions are no-op if no providers are enabled
func Initialize(opts ImageSourceOpts) {
	once.Do(func() {
		if opts.EnableCAA {
			imgsrc.caaEnabled = true
		}
		if opts.EnableDeezer {
			imgsrc.deezerEnabled = true
			imgsrc.deezerC = NewDeezerClient()
		}
		if opts.EnableSubsonic {
			imgsrc.subsonicEnabled = true
			imgsrc.subsonicC = NewSubsonicClient()
		}
		if opts.EnableSpotify {
			imgsrc.spotifyEnabled = true
			imgsrc.spotifyC = NewSpotifyClient()
		}
	})
}

func Shutdown() {
	if imgsrc.deezerC != nil {
		imgsrc.deezerC.Shutdown()
	}
	if imgsrc.subsonicC != nil {
		imgsrc.subsonicC.Shutdown()
	}
	if imgsrc.spotifyC != nil {
		imgsrc.spotifyC.Shutdown()
	}
}

func GetArtistImage(ctx context.Context, opts ArtistImageOpts) (string, error) {
	l := logger.FromContext(ctx)
	if imgsrc.spotifyEnabled {
		l.Debug().Msg("Attempting to find artist image from Spotify")
		img, err := imgsrc.spotifyC.GetArtistImage(ctx, opts.Aliases)
		if err != nil {
			l.Debug().Err(err).Msg("Could not find artist image from Spotify")
		} else if img != "" {
			return img, nil
		}
	}
	if imgsrc.subsonicEnabled {
		if len(opts.Aliases) == 0 {
			l.Debug().Msg("GetArtistImage: no aliases provided, skipping Subsonic")
		} else {
			img, err := imgsrc.subsonicC.GetArtistImage(ctx, opts.Aliases[0])
			if err != nil {
				l.Debug().Err(err).Msg("Could not find artist image from Subsonic")
			} else if img != "" {
				return img, nil
			}
		}
	}
	if imgsrc.deezerEnabled {
		l.Debug().Msg("Attempting to find artist image from Deezer")
		img, err := imgsrc.deezerC.GetArtistImages(ctx, opts.Aliases)
		if err != nil {
			return "", err
		}
		return img, nil
	}
	l.Warn().Msg("GetArtistImage: No image providers are enabled")
	return "", nil
}
func GetAlbumImage(ctx context.Context, opts AlbumImageOpts) (string, error) {
	l := logger.FromContext(ctx)
	if imgsrc.spotifyEnabled {
		l.Debug().Msg("Attempting to find album image from Spotify")
		img, err := imgsrc.spotifyC.GetAlbumImage(ctx, opts.Artists, opts.Album)
		if err != nil {
			return "", err
		}
		if img != "" {
			return img, nil
		}
	}
	if imgsrc.subsonicEnabled {
		if len(opts.Artists) == 0 {
			l.Debug().Msg("GetAlbumImage: no artists provided, skipping Subsonic")
		} else {
			img, err := imgsrc.subsonicC.GetAlbumImage(ctx, opts.Artists[0], opts.Album)
			if err != nil {
				return "", err
			}
			if img != "" {
				return img, nil
			}
			l.Debug().Msg("Could not find album cover from Subsonic")
		}
	}
	if imgsrc.caaEnabled {
		l.Debug().Msg("Attempting to find album image from CoverArtArchive")
		if opts.ReleaseMbzID != nil && *opts.ReleaseMbzID != uuid.Nil {
			url := fmt.Sprintf(caaBaseUrl+"/release/%s", opts.ReleaseMbzID.String())
			img, err := caaCoverImageExtract(ctx, url)
			if err != nil {
				return "", err
			}
			if img != "" {
				return img, nil
			}

			frontURL := fmt.Sprintf(caaBaseUrl+"/release/%s/front", opts.ReleaseMbzID.String())
			img, status, err := caaFrontImage(frontURL)
			if err != nil {
				return "", err
			}
			if img != "" {
				return img, nil
			}
			l.Debug().Str("url", frontURL).Str("status", status).Msg("Could not find album cover from CoverArtArchive with MusicBrainz release ID")
		}
		if opts.ReleaseGroupMbzID != nil && *opts.ReleaseGroupMbzID != uuid.Nil {
			url := fmt.Sprintf(caaBaseUrl+"/release-group/%s", opts.ReleaseGroupMbzID.String())
			img, err := caaCoverImageExtract(ctx, url)
			if err != nil {
				return "", err
			}
			if img != "" {
				return img, nil
			}

			frontURL := fmt.Sprintf(caaBaseUrl+"/release-group/%s/front", opts.ReleaseGroupMbzID.String())
			img, status, err := caaFrontImage(frontURL)
			if err != nil {
				return "", err
			}
			if img != "" {
				return img, nil
			}
			l.Debug().Str("url", frontURL).Str("status", status).Msg("Could not find album cover from CoverArtArchive with MusicBrainz release group ID")
		}
	}
	if imgsrc.deezerEnabled {
		l.Debug().Msg("Attempting to find album image from Deezer")
		img, err := imgsrc.deezerC.GetAlbumImages(ctx, opts.Artists, opts.Album)
		if err != nil {
			return "", err
		}
		return img, nil
	}
	l.Warn().Msg("GetAlbumImage: No image providers are enabled")
	return "", nil
}
