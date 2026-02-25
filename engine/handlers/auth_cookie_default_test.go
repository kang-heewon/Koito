//go:build auth_cookie_default
// +build auth_cookie_default

package handlers

// Run with: go test -tags=auth_cookie_default ./engine/handlers/ -run TestLoginHandlerCookieSecureFlagDefault -v

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gabehf/koito/internal/cfg"
	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandlerCookieSecureFlagDefault(t *testing.T) {
	os.Setenv("KOITO_DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Unsetenv("KOITO_SECURE_COOKIES")
	defer os.Unsetenv("KOITO_DATABASE_URL")
	defer os.Unsetenv("KOITO_SECURE_COOKIES")

	cfg.Load(os.Getenv, "test")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	store := &mockAuthDB{
		user: &models.User{
			ID:       1,
			Username: "testuser",
			Password: hashedPassword,
		},
		sessionID: uuid.New(),
	}

	req := httptest.NewRequest("POST", "/login", strings.NewReader("username=testuser&password=password"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	LoginHandler(store)(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	cookies := rr.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Secure != false {
		t.Errorf("cookie.Secure = %v, want false (default)", cookie.Secure)
	}
}

func TestLogoutHandlerCookieSecureFlagDefault(t *testing.T) {
	os.Setenv("KOITO_DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Unsetenv("KOITO_SECURE_COOKIES")
	defer os.Unsetenv("KOITO_DATABASE_URL")
	defer os.Unsetenv("KOITO_SECURE_COOKIES")

	cfg.Load(os.Getenv, "test")

	store := &mockAuthDB{}
	req := httptest.NewRequest("POST", "/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "koito_session",
		Value: uuid.New().String(),
	})
	rr := httptest.NewRecorder()

	LogoutHandler(store)(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	cookies := rr.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Secure != false {
		t.Errorf("cookie.Secure = %v, want false (default)", cookie.Secure)
	}
}

type mockAuthDB struct {
	user      *models.User
	sessionID uuid.UUID
}

func (m *mockAuthDB) GetUserByUsername(_ context.Context, username string) (*models.User, error) {
	if m.user != nil && m.user.Username == username {
		return m.user, nil
	}
	return nil, nil
}
func (m *mockAuthDB) SaveSession(_ context.Context, userID int32, expiresAt time.Time, rememberMe bool) (*models.Session, error) {
	return &models.Session{ID: m.sessionID, UserID: userID, ExpiresAt: expiresAt}, nil
}
func (m *mockAuthDB) DeleteSession(_ context.Context, sid uuid.UUID) error { return nil }
func (m *mockAuthDB) UpdateUser(_ context.Context, opts db.UpdateUserOpts) error { return nil }
func (m *mockAuthDB) GetArtist(ctx context.Context, opts db.GetArtistOpts) (*models.Artist, error) { return nil, nil }
func (m *mockAuthDB) GetAlbum(ctx context.Context, opts db.GetAlbumOpts) (*models.Album, error) { return nil, nil }
func (m *mockAuthDB) GetTrack(ctx context.Context, opts db.GetTrackOpts) (*models.Track, error) { return nil, nil }
func (m *mockAuthDB) GetArtistsForAlbum(ctx context.Context, id int32) ([]*models.Artist, error) { return nil, nil }
func (m *mockAuthDB) GetArtistsForTrack(ctx context.Context, id int32) ([]*models.Artist, error) { return nil, nil }
func (m *mockAuthDB) GetTopTracksPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Track], error) { return nil, nil }
func (m *mockAuthDB) GetTopArtistsPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Artist], error) { return nil, nil }
func (m *mockAuthDB) GetTopAlbumsPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Album], error) { return nil, nil }
func (m *mockAuthDB) GetListensPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Listen], error) { return nil, nil }
func (m *mockAuthDB) GetListenActivity(ctx context.Context, opts db.ListenActivityOpts) ([]db.ListenActivityItem, error) { return nil, nil }
func (m *mockAuthDB) GetAllArtistAliases(ctx context.Context, id int32) ([]models.Alias, error) { return nil, nil }
func (m *mockAuthDB) GetAllAlbumAliases(ctx context.Context, id int32) ([]models.Alias, error) { return nil, nil }
func (m *mockAuthDB) GetAllTrackAliases(ctx context.Context, id int32) ([]models.Alias, error) { return nil, nil }
func (m *mockAuthDB) GetApiKeysByUserID(ctx context.Context, id int32) ([]models.ApiKey, error) { return nil, nil }
func (m *mockAuthDB) GetUserBySession(ctx context.Context, sessionId uuid.UUID) (*models.User, error) { return nil, nil }
func (m *mockAuthDB) GetSession(ctx context.Context, sessionId uuid.UUID) (*models.Session, error) { return nil, nil }
func (m *mockAuthDB) GetUserByApiKey(ctx context.Context, key string) (*models.User, error) { return nil, nil }
func (m *mockAuthDB) SaveArtist(ctx context.Context, opts db.SaveArtistOpts) (*models.Artist, error) { return nil, nil }
func (m *mockAuthDB) SaveArtistAliases(ctx context.Context, id int32, aliases []string, source string) error { return nil }
func (m *mockAuthDB) SaveArtistGenres(ctx context.Context, id int32, genres []string) error { return nil }
func (m *mockAuthDB) SaveAlbum(ctx context.Context, opts db.SaveAlbumOpts) (*models.Album, error) { return nil, nil }
func (m *mockAuthDB) SaveAlbumAliases(ctx context.Context, id int32, aliases []string, source string) error { return nil }
func (m *mockAuthDB) SaveAlbumGenres(ctx context.Context, id int32, genres []string) error { return nil }
func (m *mockAuthDB) SaveTrack(ctx context.Context, opts db.SaveTrackOpts) (*models.Track, error) { return nil, nil }
func (m *mockAuthDB) SaveTrackAliases(ctx context.Context, id int32, aliases []string, source string) error { return nil }
func (m *mockAuthDB) SaveListen(ctx context.Context, opts db.SaveListenOpts) error { return nil }
func (m *mockAuthDB) SaveUser(ctx context.Context, opts db.SaveUserOpts) (*models.User, error) { return nil, nil }
func (m *mockAuthDB) SaveApiKey(ctx context.Context, opts db.SaveApiKeyOpts) (*models.ApiKey, error) { return nil, nil }
func (m *mockAuthDB) UpdateArtist(ctx context.Context, opts db.UpdateArtistOpts) error { return nil }
func (m *mockAuthDB) UpdateTrack(ctx context.Context, opts db.UpdateTrackOpts) error { return nil }
func (m *mockAuthDB) UpdateAlbum(ctx context.Context, opts db.UpdateAlbumOpts) error { return nil }
func (m *mockAuthDB) AddArtistsToAlbum(ctx context.Context, opts db.AddArtistsToAlbumOpts) error { return nil }
func (m *mockAuthDB) UpdateApiKeyLabel(ctx context.Context, opts db.UpdateApiKeyLabelOpts) error { return nil }
func (m *mockAuthDB) RefreshSession(ctx context.Context, sessionId uuid.UUID, expiresAt time.Time) error { return nil }
func (m *mockAuthDB) SetPrimaryArtistAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockAuthDB) SetPrimaryAlbumAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockAuthDB) SetPrimaryTrackAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockAuthDB) SetPrimaryAlbumArtist(ctx context.Context, id int32, artistId int32, value bool) error { return nil }
func (m *mockAuthDB) SetPrimaryTrackArtist(ctx context.Context, id int32, artistId int32, value bool) error { return nil }
func (m *mockAuthDB) DeleteArtist(ctx context.Context, id int32) error { return nil }
func (m *mockAuthDB) DeleteAlbum(ctx context.Context, id int32) error { return nil }
func (m *mockAuthDB) DeleteTrack(ctx context.Context, id int32) error { return nil }
func (m *mockAuthDB) DeleteListen(ctx context.Context, trackId int32, listenedAt time.Time) error { return nil }
func (m *mockAuthDB) DeleteArtistAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockAuthDB) DeleteAlbumAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockAuthDB) DeleteTrackAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockAuthDB) DeleteApiKey(ctx context.Context, id int32) error { return nil }
func (m *mockAuthDB) CountListens(ctx context.Context, period db.Period) (int64, error) { return 0, nil }
func (m *mockAuthDB) CountTracks(ctx context.Context, period db.Period) (int64, error) { return 0, nil }
func (m *mockAuthDB) CountAlbums(ctx context.Context, period db.Period) (int64, error) { return 0, nil }
func (m *mockAuthDB) CountArtists(ctx context.Context, period db.Period) (int64, error) { return 0, nil }
func (m *mockAuthDB) CountTimeListened(ctx context.Context, period db.Period) (int64, error) { return 0, nil }
func (m *mockAuthDB) CountTimeListenedToItem(ctx context.Context, opts db.TimeListenedOpts) (int64, error) { return 0, nil }
func (m *mockAuthDB) CountUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockAuthDB) GetGenreStatsByListenCount(ctx context.Context, period db.Period) ([]db.GenreStat, error) { return nil, nil }
func (m *mockAuthDB) GetGenreStatsByTimeListened(ctx context.Context, period db.Period) ([]db.GenreStat, error) { return nil, nil }
func (m *mockAuthDB) GetWrappedStats(ctx context.Context, year int, userID int32) (*db.WrappedStats, error) { return nil, nil }
func (m *mockAuthDB) GetTracksToRevisit(ctx context.Context, opts db.GetRecommendationsOpts) ([]db.TrackRecommendation, error) { return nil, nil }
func (m *mockAuthDB) SearchArtists(ctx context.Context, q string) ([]*models.Artist, error) { return nil, nil }
func (m *mockAuthDB) SearchAlbums(ctx context.Context, q string) ([]*models.Album, error) { return nil, nil }
func (m *mockAuthDB) SearchTracks(ctx context.Context, q string) ([]*models.Track, error) { return nil, nil }
func (m *mockAuthDB) MergeTracks(ctx context.Context, fromId, toId int32) error { return nil }
func (m *mockAuthDB) MergeAlbums(ctx context.Context, fromId, toId int32, replaceImage bool) error { return nil }
func (m *mockAuthDB) MergeArtists(ctx context.Context, fromId, toId int32, replaceImage bool) error { return nil }
func (m *mockAuthDB) ImageHasAssociation(ctx context.Context, image uuid.UUID) (bool, error) { return false, nil }
func (m *mockAuthDB) GetImageSource(ctx context.Context, image uuid.UUID) (string, error) { return "", nil }
func (m *mockAuthDB) AlbumsWithoutImages(ctx context.Context, from int32) ([]*models.Album, error) { return nil, nil }
func (m *mockAuthDB) AlbumsWithoutGenres(ctx context.Context, from int32) ([]db.ItemWithMbzID, error) { return nil, nil }
func (m *mockAuthDB) ArtistsWithoutGenres(ctx context.Context, from int32) ([]db.ItemWithMbzID, error) { return nil, nil }
func (m *mockAuthDB) TracksWithoutDuration(ctx context.Context, lastID int32) ([]db.TrackWithMbzID, error) { return nil, nil }
func (m *mockAuthDB) UpdateTrackDuration(ctx context.Context, id int32, duration int32) error { return nil }
func (m *mockAuthDB) GetExportPage(ctx context.Context, opts db.GetExportPageOpts) ([]*db.ExportItem, error) { return nil, nil }
func (m *mockAuthDB) Ping(ctx context.Context) error { return nil }
func (m *mockAuthDB) Close(ctx context.Context) {}
