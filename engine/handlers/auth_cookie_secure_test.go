//go:build auth_cookie_secure
// +build auth_cookie_secure

package handlers

// Run with: go test -tags=auth_cookie_secure ./engine/handlers/ -run "TestLoginHandlerCookieSecureFlagTrue|TestLogoutHandlerCookieSecureFlagTrue" -v

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

func TestLoginHandlerCookieSecureFlagTrue(t *testing.T) {
	os.Setenv("KOITO_DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Setenv("KOITO_SECURE_COOKIES", "true")
	defer os.Unsetenv("KOITO_DATABASE_URL")
	defer os.Unsetenv("KOITO_SECURE_COOKIES")

	cfg.Load(os.Getenv, "test")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	store := &mockSecureAuthDB{
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
	if cookie.Secure != true {
		t.Errorf("cookie.Secure = %v, want true", cookie.Secure)
	}
}

func TestLogoutHandlerCookieSecureFlagTrue(t *testing.T) {
	os.Setenv("KOITO_DATABASE_URL", "postgres://test:test@localhost:5432/test")
	os.Setenv("KOITO_SECURE_COOKIES", "true")
	defer os.Unsetenv("KOITO_DATABASE_URL")
	defer os.Unsetenv("KOITO_SECURE_COOKIES")

	cfg.Load(os.Getenv, "test")

	store := &mockSecureAuthDB{}
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
	if cookie.Secure != true {
		t.Errorf("cookie.Secure = %v, want true", cookie.Secure)
	}
}

type mockSecureAuthDB struct {
	user      *models.User
	sessionID uuid.UUID
}

func (m *mockSecureAuthDB) GetUserByUsername(_ context.Context, username string) (*models.User, error) {
	if m.user != nil && m.user.Username == username {
		return m.user, nil
	}
	return nil, nil
}
func (m *mockSecureAuthDB) SaveSession(_ context.Context, userID int32, expiresAt time.Time, rememberMe bool) (*models.Session, error) {
	return &models.Session{ID: m.sessionID, UserID: userID, ExpiresAt: expiresAt}, nil
}
func (m *mockSecureAuthDB) DeleteSession(_ context.Context, sid uuid.UUID) error { return nil }
func (m *mockSecureAuthDB) UpdateUser(_ context.Context, opts db.UpdateUserOpts) error { return nil }
func (m *mockSecureAuthDB) GetArtist(ctx context.Context, opts db.GetArtistOpts) (*models.Artist, error) { return nil, nil }
func (m *mockSecureAuthDB) GetAlbum(ctx context.Context, opts db.GetAlbumOpts) (*models.Album, error) { return nil, nil }
func (m *mockSecureAuthDB) GetTrack(ctx context.Context, opts db.GetTrackOpts) (*models.Track, error) { return nil, nil }
func (m *mockSecureAuthDB) GetArtistsForAlbum(ctx context.Context, id int32) ([]*models.Artist, error) { return nil, nil }
func (m *mockSecureAuthDB) GetArtistsForTrack(ctx context.Context, id int32) ([]*models.Artist, error) { return nil, nil }
func (m *mockSecureAuthDB) GetTopTracksPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Track], error) { return nil, nil }
func (m *mockSecureAuthDB) GetTopArtistsPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Artist], error) { return nil, nil }
func (m *mockSecureAuthDB) GetTopAlbumsPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Album], error) { return nil, nil }
func (m *mockSecureAuthDB) GetListensPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Listen], error) { return nil, nil }
func (m *mockSecureAuthDB) GetListenActivity(ctx context.Context, opts db.ListenActivityOpts) ([]db.ListenActivityItem, error) { return nil, nil }
func (m *mockSecureAuthDB) GetAllArtistAliases(ctx context.Context, id int32) ([]models.Alias, error) { return nil, nil }
func (m *mockSecureAuthDB) GetAllAlbumAliases(ctx context.Context, id int32) ([]models.Alias, error) { return nil, nil }
func (m *mockSecureAuthDB) GetAllTrackAliases(ctx context.Context, id int32) ([]models.Alias, error) { return nil, nil }
func (m *mockSecureAuthDB) GetApiKeysByUserID(ctx context.Context, id int32) ([]models.ApiKey, error) { return nil, nil }
func (m *mockSecureAuthDB) GetUserBySession(ctx context.Context, sessionId uuid.UUID) (*models.User, error) { return nil, nil }
func (m *mockSecureAuthDB) GetSession(ctx context.Context, sessionId uuid.UUID) (*models.Session, error) { return nil, nil }
func (m *mockSecureAuthDB) GetUserByApiKey(ctx context.Context, key string) (*models.User, error) { return nil, nil }
func (m *mockSecureAuthDB) SaveArtist(ctx context.Context, opts db.SaveArtistOpts) (*models.Artist, error) { return nil, nil }
func (m *mockSecureAuthDB) SaveArtistAliases(ctx context.Context, id int32, aliases []string, source string) error { return nil }
func (m *mockSecureAuthDB) SaveArtistGenres(ctx context.Context, id int32, genres []string) error { return nil }
func (m *mockSecureAuthDB) SaveAlbum(ctx context.Context, opts db.SaveAlbumOpts) (*models.Album, error) { return nil, nil }
func (m *mockSecureAuthDB) SaveAlbumAliases(ctx context.Context, id int32, aliases []string, source string) error { return nil }
func (m *mockSecureAuthDB) SaveAlbumGenres(ctx context.Context, id int32, genres []string) error { return nil }
func (m *mockSecureAuthDB) SaveTrack(ctx context.Context, opts db.SaveTrackOpts) (*models.Track, error) { return nil, nil }
func (m *mockSecureAuthDB) SaveTrackAliases(ctx context.Context, id int32, aliases []string, source string) error { return nil }
func (m *mockSecureAuthDB) SaveListen(ctx context.Context, opts db.SaveListenOpts) error { return nil }
func (m *mockSecureAuthDB) SaveUser(ctx context.Context, opts db.SaveUserOpts) (*models.User, error) { return nil, nil }
func (m *mockSecureAuthDB) SaveApiKey(ctx context.Context, opts db.SaveApiKeyOpts) (*models.ApiKey, error) { return nil, nil }
func (m *mockSecureAuthDB) UpdateArtist(ctx context.Context, opts db.UpdateArtistOpts) error { return nil }
func (m *mockSecureAuthDB) UpdateTrack(ctx context.Context, opts db.UpdateTrackOpts) error { return nil }
func (m *mockSecureAuthDB) UpdateAlbum(ctx context.Context, opts db.UpdateAlbumOpts) error { return nil }
func (m *mockSecureAuthDB) AddArtistsToAlbum(ctx context.Context, opts db.AddArtistsToAlbumOpts) error { return nil }
func (m *mockSecureAuthDB) UpdateApiKeyLabel(ctx context.Context, opts db.UpdateApiKeyLabelOpts) error { return nil }
func (m *mockSecureAuthDB) RefreshSession(ctx context.Context, sessionId uuid.UUID, expiresAt time.Time) error { return nil }
func (m *mockSecureAuthDB) SetPrimaryArtistAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockSecureAuthDB) SetPrimaryAlbumAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockSecureAuthDB) SetPrimaryTrackAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockSecureAuthDB) SetPrimaryAlbumArtist(ctx context.Context, id int32, artistId int32, value bool) error { return nil }
func (m *mockSecureAuthDB) SetPrimaryTrackArtist(ctx context.Context, id int32, artistId int32, value bool) error { return nil }
func (m *mockSecureAuthDB) DeleteArtist(ctx context.Context, id int32) error { return nil }
func (m *mockSecureAuthDB) DeleteAlbum(ctx context.Context, id int32) error { return nil }
func (m *mockSecureAuthDB) DeleteTrack(ctx context.Context, id int32) error { return nil }
func (m *mockSecureAuthDB) DeleteListen(ctx context.Context, trackId int32, listenedAt time.Time) error { return nil }
func (m *mockSecureAuthDB) DeleteArtistAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockSecureAuthDB) DeleteAlbumAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockSecureAuthDB) DeleteTrackAlias(ctx context.Context, id int32, alias string) error { return nil }
func (m *mockSecureAuthDB) DeleteApiKey(ctx context.Context, id int32) error { return nil }
func (m *mockSecureAuthDB) CountListens(ctx context.Context, period db.Period) (int64, error) { return 0, nil }
func (m *mockSecureAuthDB) CountTracks(ctx context.Context, period db.Period) (int64, error) { return 0, nil }
func (m *mockSecureAuthDB) CountAlbums(ctx context.Context, period db.Period) (int64, error) { return 0, nil }
func (m *mockSecureAuthDB) CountArtists(ctx context.Context, period db.Period) (int64, error) { return 0, nil }
func (m *mockSecureAuthDB) CountTimeListened(ctx context.Context, period db.Period) (int64, error) { return 0, nil }
func (m *mockSecureAuthDB) CountTimeListenedToItem(ctx context.Context, opts db.TimeListenedOpts) (int64, error) { return 0, nil }
func (m *mockSecureAuthDB) CountUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockSecureAuthDB) GetGenreStatsByListenCount(ctx context.Context, period db.Period) ([]db.GenreStat, error) { return nil, nil }
func (m *mockSecureAuthDB) GetGenreStatsByTimeListened(ctx context.Context, period db.Period) ([]db.GenreStat, error) { return nil, nil }
func (m *mockSecureAuthDB) GetWrappedStats(ctx context.Context, year int, userID int32) (*db.WrappedStats, error) { return nil, nil }
func (m *mockSecureAuthDB) GetTracksToRevisit(ctx context.Context, opts db.GetRecommendationsOpts) ([]db.TrackRecommendation, error) { return nil, nil }
func (m *mockSecureAuthDB) SearchArtists(ctx context.Context, q string) ([]*models.Artist, error) { return nil, nil }
func (m *mockSecureAuthDB) SearchAlbums(ctx context.Context, q string) ([]*models.Album, error) { return nil, nil }
func (m *mockSecureAuthDB) SearchTracks(ctx context.Context, q string) ([]*models.Track, error) { return nil, nil }
func (m *mockSecureAuthDB) MergeTracks(ctx context.Context, fromId, toId int32) error { return nil }
func (m *mockSecureAuthDB) MergeAlbums(ctx context.Context, fromId, toId int32, replaceImage bool) error { return nil }
func (m *mockSecureAuthDB) MergeArtists(ctx context.Context, fromId, toId int32, replaceImage bool) error { return nil }
func (m *mockSecureAuthDB) ImageHasAssociation(ctx context.Context, image uuid.UUID) (bool, error) { return false, nil }
func (m *mockSecureAuthDB) GetImageSource(ctx context.Context, image uuid.UUID) (string, error) { return "", nil }
func (m *mockSecureAuthDB) AlbumsWithoutImages(ctx context.Context, from int32) ([]*models.Album, error) { return nil, nil }
func (m *mockSecureAuthDB) AlbumsWithoutGenres(ctx context.Context, from int32) ([]db.ItemWithMbzID, error) { return nil, nil }
func (m *mockSecureAuthDB) ArtistsWithoutGenres(ctx context.Context, from int32) ([]db.ItemWithMbzID, error) { return nil, nil }
func (m *mockSecureAuthDB) TracksWithoutDuration(ctx context.Context, lastID int32) ([]db.TrackWithMbzID, error) { return nil, nil }
func (m *mockSecureAuthDB) UpdateTrackDuration(ctx context.Context, id int32, duration int32) error { return nil }
func (m *mockSecureAuthDB) GetExportPage(ctx context.Context, opts db.GetExportPageOpts) ([]*db.ExportItem, error) { return nil, nil }
func (m *mockSecureAuthDB) Ping(ctx context.Context) error { return nil }
func (m *mockSecureAuthDB) Close(ctx context.Context) {}
