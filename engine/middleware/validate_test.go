package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/models"
	"github.com/google/uuid"
)

// mockDBForValidateSession is a mock for testing ValidateSession
type mockDBForValidateSession struct {
	user              *models.User
	session           *models.Session
	sessionID         uuid.UUID
	capturedExpiresAt time.Time
	refreshCalled     bool
	getSessionCalled  bool
}

func (m *mockDBForValidateSession) GetUserBySession(ctx context.Context, sid uuid.UUID) (*models.User, error) {
	if m.user != nil && m.sessionID == sid {
		return m.user, nil
	}
	return nil, nil
}

func (m *mockDBForValidateSession) GetSession(ctx context.Context, sid uuid.UUID) (*models.Session, error) {
	m.getSessionCalled = true
	if m.session != nil && m.sessionID == sid {
		return m.session, nil
	}
	return nil, nil
}

func (m *mockDBForValidateSession) RefreshSession(ctx context.Context, sid uuid.UUID, expiresAt time.Time) error {
	m.refreshCalled = true
	m.capturedExpiresAt = expiresAt
	return nil
}

// Stub methods required by db.DB interface
func (m *mockDBForValidateSession) GetArtist(ctx context.Context, opts db.GetArtistOpts) (*models.Artist, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetAlbum(ctx context.Context, opts db.GetAlbumOpts) (*models.Album, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetTrack(ctx context.Context, opts db.GetTrackOpts) (*models.Track, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetArtistsForAlbum(ctx context.Context, id int32) ([]*models.Artist, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetArtistsForTrack(ctx context.Context, id int32) ([]*models.Artist, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetTopTracksPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Track], error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetTopArtistsPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Artist], error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetTopAlbumsPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Album], error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetListensPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Listen], error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetListenActivity(ctx context.Context, opts db.ListenActivityOpts) ([]db.ListenActivityItem, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetAllArtistAliases(ctx context.Context, id int32) ([]models.Alias, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetAllAlbumAliases(ctx context.Context, id int32) ([]models.Alias, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetAllTrackAliases(ctx context.Context, id int32) ([]models.Alias, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetApiKeysByUserID(ctx context.Context, id int32) ([]models.ApiKey, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetUserByApiKey(ctx context.Context, key string) (*models.User, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) SaveArtist(ctx context.Context, opts db.SaveArtistOpts) (*models.Artist, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) SaveArtistAliases(ctx context.Context, id int32, aliases []string, source string) error {
	return nil
}

func (m *mockDBForValidateSession) SaveArtistGenres(ctx context.Context, id int32, genres []string) error {
	return nil
}

func (m *mockDBForValidateSession) SaveAlbum(ctx context.Context, opts db.SaveAlbumOpts) (*models.Album, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) SaveAlbumAliases(ctx context.Context, id int32, aliases []string, source string) error {
	return nil
}

func (m *mockDBForValidateSession) SaveAlbumGenres(ctx context.Context, id int32, genres []string) error {
	return nil
}

func (m *mockDBForValidateSession) SaveTrack(ctx context.Context, opts db.SaveTrackOpts) (*models.Track, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) SaveTrackAliases(ctx context.Context, id int32, aliases []string, source string) error {
	return nil
}

func (m *mockDBForValidateSession) SaveListen(ctx context.Context, opts db.SaveListenOpts) error {
	return nil
}

func (m *mockDBForValidateSession) SaveUser(ctx context.Context, opts db.SaveUserOpts) (*models.User, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) SaveApiKey(ctx context.Context, opts db.SaveApiKeyOpts) (*models.ApiKey, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) SaveSession(ctx context.Context, userId int32, expiresAt time.Time, persistent bool) (*models.Session, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) UpdateArtist(ctx context.Context, opts db.UpdateArtistOpts) error {
	return nil
}

func (m *mockDBForValidateSession) UpdateTrack(ctx context.Context, opts db.UpdateTrackOpts) error {
	return nil
}

func (m *mockDBForValidateSession) UpdateAlbum(ctx context.Context, opts db.UpdateAlbumOpts) error {
	return nil
}

func (m *mockDBForValidateSession) AddArtistsToAlbum(ctx context.Context, opts db.AddArtistsToAlbumOpts) error {
	return nil
}

func (m *mockDBForValidateSession) UpdateUser(ctx context.Context, opts db.UpdateUserOpts) error {
	return nil
}

func (m *mockDBForValidateSession) UpdateApiKeyLabel(ctx context.Context, opts db.UpdateApiKeyLabelOpts) error {
	return nil
}

func (m *mockDBForValidateSession) SetPrimaryArtistAlias(ctx context.Context, id int32, alias string) error {
	return nil
}

func (m *mockDBForValidateSession) SetPrimaryAlbumAlias(ctx context.Context, id int32, alias string) error {
	return nil
}

func (m *mockDBForValidateSession) SetPrimaryTrackAlias(ctx context.Context, id int32, alias string) error {
	return nil
}

func (m *mockDBForValidateSession) SetPrimaryAlbumArtist(ctx context.Context, id int32, artistId int32, value bool) error {
	return nil
}

func (m *mockDBForValidateSession) SetPrimaryTrackArtist(ctx context.Context, id int32, artistId int32, value bool) error {
	return nil
}

func (m *mockDBForValidateSession) DeleteArtist(ctx context.Context, id int32) error {
	return nil
}

func (m *mockDBForValidateSession) DeleteAlbum(ctx context.Context, id int32) error {
	return nil
}

func (m *mockDBForValidateSession) DeleteTrack(ctx context.Context, id int32) error {
	return nil
}

func (m *mockDBForValidateSession) DeleteListen(ctx context.Context, trackId int32, listenedAt time.Time) error {
	return nil
}

func (m *mockDBForValidateSession) DeleteArtistAlias(ctx context.Context, id int32, alias string) error {
	return nil
}

func (m *mockDBForValidateSession) DeleteAlbumAlias(ctx context.Context, id int32, alias string) error {
	return nil
}

func (m *mockDBForValidateSession) DeleteTrackAlias(ctx context.Context, id int32, alias string) error {
	return nil
}

func (m *mockDBForValidateSession) DeleteSession(ctx context.Context, sid uuid.UUID) error {
	return nil
}

func (m *mockDBForValidateSession) DeleteApiKey(ctx context.Context, id int32) error {
	return nil
}

func (m *mockDBForValidateSession) CountListens(ctx context.Context, period db.Period) (int64, error) {
	return 0, nil
}

func (m *mockDBForValidateSession) CountTracks(ctx context.Context, period db.Period) (int64, error) {
	return 0, nil
}

func (m *mockDBForValidateSession) CountAlbums(ctx context.Context, period db.Period) (int64, error) {
	return 0, nil
}

func (m *mockDBForValidateSession) CountArtists(ctx context.Context, period db.Period) (int64, error) {
	return 0, nil
}

func (m *mockDBForValidateSession) CountTimeListened(ctx context.Context, period db.Period) (int64, error) {
	return 0, nil
}

func (m *mockDBForValidateSession) CountTimeListenedToItem(ctx context.Context, opts db.TimeListenedOpts) (int64, error) {
	return 0, nil
}

func (m *mockDBForValidateSession) CountUsers(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockDBForValidateSession) GetGenreStatsByListenCount(ctx context.Context, period db.Period) ([]db.GenreStat, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetGenreStatsByTimeListened(ctx context.Context, period db.Period) ([]db.GenreStat, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetWrappedStats(ctx context.Context, year int, userID int32) (*db.WrappedStats, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) GetTracksToRevisit(ctx context.Context, opts db.GetRecommendationsOpts) ([]db.TrackRecommendation, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) SearchArtists(ctx context.Context, q string) ([]*models.Artist, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) SearchAlbums(ctx context.Context, q string) ([]*models.Album, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) SearchTracks(ctx context.Context, q string) ([]*models.Track, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) MergeTracks(ctx context.Context, fromId, toId int32) error {
	return nil
}

func (m *mockDBForValidateSession) MergeAlbums(ctx context.Context, fromId, toId int32, replaceImage bool) error {
	return nil
}

func (m *mockDBForValidateSession) MergeArtists(ctx context.Context, fromId, toId int32, replaceImage bool) error {
	return nil
}

func (m *mockDBForValidateSession) ImageHasAssociation(ctx context.Context, image uuid.UUID) (bool, error) {
	return false, nil
}

func (m *mockDBForValidateSession) GetImageSource(ctx context.Context, image uuid.UUID) (string, error) {
	return "", nil
}

func (m *mockDBForValidateSession) AlbumsWithoutImages(ctx context.Context, from int32) ([]*models.Album, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) AlbumsWithoutGenres(ctx context.Context, from int32) ([]db.ItemWithMbzID, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) ArtistsWithoutGenres(ctx context.Context, from int32) ([]db.ItemWithMbzID, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) TracksWithoutDuration(ctx context.Context, lastID int32) ([]db.TrackWithMbzID, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) UpdateTrackDuration(ctx context.Context, id int32, duration int32) error {
	return nil
}

func (m *mockDBForValidateSession) GetExportPage(ctx context.Context, opts db.GetExportPageOpts) ([]*db.ExportItem, error) {
	return nil, nil
}

func (m *mockDBForValidateSession) Ping(ctx context.Context) error {
	return nil
}

func (m *mockDBForValidateSession) Close(ctx context.Context) {}

func TestValidateSession_NonPersistent(t *testing.T) {
	sessionID := uuid.New()
	now := time.Now()

	mock := &mockDBForValidateSession{
		user: &models.User{
			ID:       1,
			Username: "testuser",
		},
		session: &models.Session{
			ID:         sessionID,
			UserID:     1,
			CreatedAt:  now,
			ExpiresAt:  now.Add(24 * time.Hour),
			Persistent: false,
		},
		sessionID: sessionID,
	}

	handler := ValidateSession(mock)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "koito_session",
		Value: sessionID.String(),
	})

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !mock.getSessionCalled {
		t.Error("GetSession was not called")
	}

	if !mock.refreshCalled {
		t.Error("RefreshSession was not called")
	}

	expectedMinExpires := now.Add(24 * time.Hour)
	expectedMaxExpires := now.Add(24*time.Hour + time.Second)
	if mock.capturedExpiresAt.Before(expectedMinExpires) || mock.capturedExpiresAt.After(expectedMaxExpires) {
		t.Errorf("non-persistent session should extend by ~24 hours, got %v", mock.capturedExpiresAt.Sub(now))
	}
}

func TestValidateSession_Persistent(t *testing.T) {
	sessionID := uuid.New()
	now := time.Now()

	mock := &mockDBForValidateSession{
		user: &models.User{
			ID:       1,
			Username: "testuser",
		},
		session: &models.Session{
			ID:         sessionID,
			UserID:     1,
			CreatedAt:  now,
			ExpiresAt:  now.Add(30 * 24 * time.Hour),
			Persistent: true,
		},
		sessionID: sessionID,
	}

	handler := ValidateSession(mock)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "koito_session",
		Value: sessionID.String(),
	})

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !mock.getSessionCalled {
		t.Error("GetSession was not called")
	}

	if !mock.refreshCalled {
		t.Error("RefreshSession was not called")
	}

	expectedMinExpires := now.Add(30 * 24 * time.Hour)
	expectedMaxExpires := now.Add(30*24*time.Hour + time.Second)
	if mock.capturedExpiresAt.Before(expectedMinExpires) || mock.capturedExpiresAt.After(expectedMaxExpires) {
		t.Errorf("persistent session should extend by ~30 days, got %v", mock.capturedExpiresAt.Sub(now))
	}
}
