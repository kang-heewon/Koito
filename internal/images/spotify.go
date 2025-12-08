package images

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/logger"
	"github.com/gabehf/koito/internal/utils"
	"github.com/gabehf/koito/queue"
)

type SpotifyClient struct {
	clientID     string
	clientSecret string
	userAgent    string
	requestQueue *queue.RequestQueue
	token        string
	tokenExpiry  time.Time
	tokenMu      sync.Mutex
	httpClient   *http.Client
}

type spotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type spotifySearchResponse struct {
	Albums struct {
		Items []spotifyAlbum `json:"items"`
	} `json:"albums"`
	Artists struct {
		Items []spotifyArtistFull `json:"items"`
	} `json:"artists"`
}

type spotifyAlbum struct {
	Name    string          `json:"name"`
	Images  []spotifyImage  `json:"images"`
	Artists []spotifyArtist `json:"artists"`
}

type spotifyImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type spotifyArtist struct {
	Name string `json:"name"`
}

type spotifyArtistFull struct {
	Name   string         `json:"name"`
	Images []spotifyImage `json:"images"`
}

const (
	spotifyTokenURL          = "https://accounts.spotify.com/api/token"
	spotifySearchFmt         = "https://api.spotify.com/v1/search?type=album&limit=5&market=KR&q=%s"
	spotifyArtistSearchFmt   = "https://api.spotify.com/v1/search?type=artist&limit=5&market=KR&q=%s"
	tokenExpiryPadding       = 60 * time.Second
)

func NewSpotifyClient() *SpotifyClient {
	return &SpotifyClient{
		clientID:     cfg.SpotifyClientID(),
		clientSecret: cfg.SpotifyClientSecret(),
		userAgent:    cfg.UserAgent(),
		requestQueue: queue.NewRequestQueue(5, 5),
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *SpotifyClient) Shutdown() {
	c.requestQueue.Shutdown()
}

func (c *SpotifyClient) queue(ctx context.Context, req *http.Request) ([]byte, int, error) {
	l := logger.FromContext(ctx)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	status := 0
	resultChan := c.requestQueue.Enqueue(func(client *http.Client, done chan<- queue.RequestResult) {
		resp, err := client.Do(req)
		if err != nil {
			l.Debug().Err(err).Str("url", req.RequestURI).Msg("Failed to contact Spotify")
			done <- queue.RequestResult{Err: err}
			return
		}
		status = resp.StatusCode
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if resp.StatusCode >= 300 || resp.StatusCode < 200 {
			err = fmt.Errorf("received non-ok status from Spotify: %s", resp.Status)
		}
		done <- queue.RequestResult{Body: body, Err: err}
	})

	result := <-resultChan
	return result.Body, status, result.Err
}

func (c *SpotifyClient) getToken(ctx context.Context) (string, error) {
	c.tokenMu.Lock()
	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		token := c.token
		c.tokenMu.Unlock()
		return token, nil
	}
	c.tokenMu.Unlock()

	l := logger.FromContext(ctx)
	req, err := http.NewRequestWithContext(ctx, "POST", spotifyTokenURL, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", fmt.Errorf("getToken: %w", err)
	}
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.clientID+":"+c.clientSecret)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("getToken: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("getToken: %w", err)
	}
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return "", fmt.Errorf("getToken: received non-ok status from Spotify: %s", resp.Status)
	}

	tokenResp := new(spotifyTokenResponse)
	err = json.Unmarshal(body, tokenResp)
	if err != nil {
		return "", fmt.Errorf("getToken: %w", err)
	}

	c.tokenMu.Lock()
	c.token = tokenResp.AccessToken
	expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	if tokenResp.ExpiresIn > int(tokenExpiryPadding.Seconds()) {
		expiry = expiry.Add(-tokenExpiryPadding)
	}
	c.tokenExpiry = expiry
	c.tokenMu.Unlock()

	l.Debug().Msg("SpotifyClient: refreshed token")
	return tokenResp.AccessToken, nil
}

func (c *SpotifyClient) GetAlbumImage(ctx context.Context, artists []string, album string) (string, error) {
	l := logger.FromContext(ctx)
	token, err := c.getToken(ctx)
	if err != nil {
		return "", fmt.Errorf("GetAlbumImage: %w", err)
	}

	artistCandidates := utils.UniqueIgnoringCase(artists)
	queries := make([]string, 0, len(artistCandidates)+1)
	for _, artist := range artistCandidates {
		if strings.TrimSpace(artist) == "" {
			continue
		}
		queries = append(queries, fmt.Sprintf(`album:"%s" artist:"%s"`, album, artist))
	}
	queries = append(queries, fmt.Sprintf(`album:"%s"`, album))

	for _, query := range queries {
		img, unauthorized, err := c.searchAlbum(ctx, token, query, album, artistCandidates)
		if unauthorized {
			c.tokenMu.Lock()
			c.token = ""
			c.tokenExpiry = time.Time{}
			c.tokenMu.Unlock()
			token, err = c.getToken(ctx)
			if err != nil {
				return "", fmt.Errorf("GetAlbumImage: %w", err)
			}
			img, unauthorized, err = c.searchAlbum(ctx, token, query, album, artistCandidates)
		}
		if err != nil {
			return "", fmt.Errorf("GetAlbumImage: %w", err)
		}
		if unauthorized {
			return "", fmt.Errorf("GetAlbumImage: Spotify returned unauthorized for album query")
		}
		if img != "" {
			l.Debug().Str("query", query).Msg("Found album image from Spotify")
			return img, nil
		}
	}

	return "", fmt.Errorf("GetAlbumImage: album image not found")
}

func (c *SpotifyClient) GetArtistImage(ctx context.Context, aliases []string) (string, error) {
	l := logger.FromContext(ctx)
	token, err := c.getToken(ctx)
	if err != nil {
		return "", fmt.Errorf("GetArtistImage: %w", err)
	}

	artistCandidates := utils.UniqueIgnoringCase(aliases)
	for _, artist := range artistCandidates {
		if strings.TrimSpace(artist) == "" {
			continue
		}
		query := fmt.Sprintf(`artist:"%s"`, artist)
		img, unauthorized, err := c.searchArtist(ctx, token, query, artist)
		if unauthorized {
			c.tokenMu.Lock()
			c.token = ""
			c.tokenExpiry = time.Time{}
			c.tokenMu.Unlock()
			token, err = c.getToken(ctx)
			if err != nil {
				return "", fmt.Errorf("GetArtistImage: %w", err)
			}
			img, unauthorized, err = c.searchArtist(ctx, token, query, artist)
		}
		if err != nil {
			return "", fmt.Errorf("GetArtistImage: %w", err)
		}
		if unauthorized {
			return "", fmt.Errorf("GetArtistImage: Spotify returned unauthorized for artist query")
		}
		if img != "" {
			l.Debug().Str("query", query).Msg("Found artist image from Spotify")
			return img, nil
		}
	}

	return "", fmt.Errorf("GetArtistImage: artist image not found")
}

func (c *SpotifyClient) searchArtist(ctx context.Context, token, query, artist string) (string, bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(spotifyArtistSearchFmt, url.QueryEscape(query)), nil)
	if err != nil {
		return "", false, fmt.Errorf("searchArtist: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	body, status, err := c.queue(ctx, req)
	if status == http.StatusUnauthorized {
		return "", true, fmt.Errorf("searchArtist: received unauthorized status")
	}
	if err != nil {
		return "", false, fmt.Errorf("searchArtist: %w", err)
	}

	resp := new(spotifySearchResponse)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return "", false, fmt.Errorf("searchArtist: %w", err)
	}

	for _, item := range resp.Artists.Items {
		if !strings.EqualFold(item.Name, artist) {
			continue
		}
		if len(item.Images) > 0 {
			return item.Images[0].URL, false, nil
		}
	}
	return "", false, nil
}

func (c *SpotifyClient) searchAlbum(ctx context.Context, token, query, album string, artists []string) (string, bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(spotifySearchFmt, url.QueryEscape(query)), nil)
	if err != nil {
		return "", false, fmt.Errorf("searchAlbum: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	body, status, err := c.queue(ctx, req)
	if status == http.StatusUnauthorized {
		return "", true, fmt.Errorf("searchAlbum: received unauthorized status")
	}
	if err != nil {
		return "", false, fmt.Errorf("searchAlbum: %w", err)
	}

	resp := new(spotifySearchResponse)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return "", false, fmt.Errorf("searchAlbum: %w", err)
	}

	for _, item := range resp.Albums.Items {
		if !strings.EqualFold(item.Name, album) {
			continue
		}
		if len(artists) > 0 && !artistMatch(item.Artists, artists) {
			continue
		}
		if len(item.Images) > 0 {
			return item.Images[0].URL, false, nil
		}
	}
	return "", false, nil
}

func artistMatch(found []spotifyArtist, candidates []string) bool {
	for _, a := range found {
		for _, candidate := range candidates {
			if strings.EqualFold(a.Name, candidate) {
				return true
			}
		}
	}
	return false
}
