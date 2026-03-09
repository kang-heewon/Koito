package summary_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/gabehf/koito/internal/db"
	"github.com/gabehf/koito/internal/models"
	"github.com/gabehf/koito/internal/summary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSummary(t *testing.T) {
	timeframe := db.Timeframe{
		From: "2024-01-01T00:00:00Z",
		To:   "2024-01-05T00:00:00Z",
	}

	store := &mockSummaryStore{
		topArtists: []db.RankedItem[*models.Artist]{
			{Item: &models.Artist{ID: 1, Name: "Artist One"}, Rank: 1},
			{Item: &models.Artist{ID: 2, Name: "Artist Two"}, Rank: 2},
		},
		topAlbums: []db.RankedItem[*models.Album]{
			{Item: &models.Album{ID: 10, Title: "Album One"}, Rank: 1},
		},
		topTracks: []db.RankedItem[*models.Track]{
			{Item: &models.Track{ID: 20, Title: "Track One"}, Rank: 1},
		},
		itemTimeListened: map[string]int64{
			"artist:1": 600,
			"artist:2": 300,
			"album:10": 1200,
			"track:20": 180,
		},
		itemListens: map[string]int64{
			"artist:1": 11,
			"artist:2": 7,
			"album:10": 9,
			"track:20": 3,
		},
		totalTimeListened: 7200,
		totalListens:      20,
		totalTracks:       8,
		totalAlbums:       5,
		totalArtists:      4,
		newTracks:         2,
		newAlbums:         1,
		newArtists:        3,
	}

	got, err := summary.GenerateSummary(context.Background(), store, 1, timeframe, "2024 Rewind")
	require.NoError(t, err)

	assert.Equal(t, "2024 Rewind", got.Title)
	assert.Equal(t, 120, got.MinutesListened)
	assert.Equal(t, 30, got.AvgMinutesPerDay)
	assert.Equal(t, 20, got.Plays)
	assert.Equal(t, float32(5), got.AvgPlaysPerDay)
	assert.Equal(t, 8, got.UniqueTracks)
	assert.Equal(t, 5, got.UniqueAlbums)
	assert.Equal(t, 4, got.UniqueArtists)
	assert.Equal(t, 2, got.NewTracks)
	assert.Equal(t, 1, got.NewAlbums)
	assert.Equal(t, 3, got.NewArtists)

	require.Len(t, got.TopArtists, 2)
	assert.Equal(t, int64(600), got.TopArtists[0].Item.TimeListened)
	assert.Equal(t, int64(11), got.TopArtists[0].Item.ListenCount)
	assert.Equal(t, int64(300), got.TopArtists[1].Item.TimeListened)
	assert.Equal(t, int64(7), got.TopArtists[1].Item.ListenCount)

	require.Len(t, got.TopAlbums, 1)
	assert.Equal(t, int64(1200), got.TopAlbums[0].Item.TimeListened)
	assert.Equal(t, int64(9), got.TopAlbums[0].Item.ListenCount)

	require.Len(t, got.TopTracks, 1)
	assert.Equal(t, int64(180), got.TopTracks[0].Item.TimeListened)
	assert.Equal(t, int64(3), got.TopTracks[0].Item.ListenCount)

	assert.Equal(t, []db.GetItemsOpts{
		{Page: 1, Limit: 5, Timeframe: timeframe},
		{Page: 1, Limit: 5, Timeframe: timeframe},
		{Page: 1, Limit: 5, Timeframe: timeframe},
	}, store.itemRequests)
	assert.Len(t, store.timeListenedRequests, 4)
	assert.Len(t, store.listenCountRequests, 4)
	for _, req := range store.timeListenedRequests {
		assert.Equal(t, timeframe, req.Timeframe)
	}
	for _, req := range store.listenCountRequests {
		assert.Equal(t, timeframe, req.Timeframe)
	}
	assert.Equal(t, []db.Timeframe{timeframe, timeframe, timeframe, timeframe, timeframe, timeframe, timeframe, timeframe}, store.aggregateTimeframes)
}

func TestGenerateSummaryReturnsWrappedError(t *testing.T) {
	expectedErr := errors.New("boom")
	store := &mockSummaryStore{
		getTopArtistsErr: expectedErr,
	}

	got, err := summary.GenerateSummary(context.Background(), store, 1, db.Timeframe{Period: string(db.PeriodYear)}, "2024 Rewind")
	require.Error(t, err)
	assert.Nil(t, got)
	assert.ErrorIs(t, err, expectedErr)
	assert.EqualError(t, err, "GenerateSummary: boom")
}

type mockSummaryStore struct {
	db.DB

	topArtists []db.RankedItem[*models.Artist]
	topAlbums  []db.RankedItem[*models.Album]
	topTracks  []db.RankedItem[*models.Track]

	itemTimeListened map[string]int64
	itemListens      map[string]int64

	totalTimeListened int64
	totalListens      int64
	totalTracks       int64
	totalAlbums       int64
	totalArtists      int64
	newTracks         int64
	newAlbums         int64
	newArtists        int64

	getTopArtistsErr error

	itemRequests         []db.GetItemsOpts
	timeListenedRequests []db.TimeListenedOpts
	listenCountRequests  []db.TimeListenedOpts
	aggregateTimeframes  []db.Timeframe
}

func (m *mockSummaryStore) GetTopArtistsPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[db.RankedItem[*models.Artist]], error) {
	if m.getTopArtistsErr != nil {
		return nil, m.getTopArtistsErr
	}
	m.itemRequests = append(m.itemRequests, opts)
	return &db.PaginatedResponse[db.RankedItem[*models.Artist]]{Items: m.topArtists}, nil
}

func (m *mockSummaryStore) GetTopAlbumsPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[db.RankedItem[*models.Album]], error) {
	m.itemRequests = append(m.itemRequests, opts)
	return &db.PaginatedResponse[db.RankedItem[*models.Album]]{Items: m.topAlbums}, nil
}

func (m *mockSummaryStore) GetTopTracksPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[db.RankedItem[*models.Track]], error) {
	m.itemRequests = append(m.itemRequests, opts)
	return &db.PaginatedResponse[db.RankedItem[*models.Track]]{Items: m.topTracks}, nil
}

func (m *mockSummaryStore) CountTimeListenedToItem(ctx context.Context, opts db.TimeListenedOpts) (int64, error) {
	m.timeListenedRequests = append(m.timeListenedRequests, opts)
	return m.itemTimeListened[itemKey(opts)], nil
}

func (m *mockSummaryStore) CountListensToItem(ctx context.Context, opts db.TimeListenedOpts) (int64, error) {
	m.listenCountRequests = append(m.listenCountRequests, opts)
	return m.itemListens[itemKey(opts)], nil
}

func (m *mockSummaryStore) CountTimeListened(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	m.aggregateTimeframes = append(m.aggregateTimeframes, timeframe)
	return m.totalTimeListened, nil
}

func (m *mockSummaryStore) CountListens(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	m.aggregateTimeframes = append(m.aggregateTimeframes, timeframe)
	return m.totalListens, nil
}

func (m *mockSummaryStore) CountTracks(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	m.aggregateTimeframes = append(m.aggregateTimeframes, timeframe)
	return m.totalTracks, nil
}

func (m *mockSummaryStore) CountAlbums(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	m.aggregateTimeframes = append(m.aggregateTimeframes, timeframe)
	return m.totalAlbums, nil
}

func (m *mockSummaryStore) CountArtists(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	m.aggregateTimeframes = append(m.aggregateTimeframes, timeframe)
	return m.totalArtists, nil
}

func (m *mockSummaryStore) CountNewTracks(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	m.aggregateTimeframes = append(m.aggregateTimeframes, timeframe)
	return m.newTracks, nil
}

func (m *mockSummaryStore) CountNewAlbums(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	m.aggregateTimeframes = append(m.aggregateTimeframes, timeframe)
	return m.newAlbums, nil
}

func (m *mockSummaryStore) CountNewArtists(ctx context.Context, timeframe db.Timeframe) (int64, error) {
	m.aggregateTimeframes = append(m.aggregateTimeframes, timeframe)
	return m.newArtists, nil
}

func itemKey(opts db.TimeListenedOpts) string {
	switch {
	case opts.ArtistID != 0:
		return "artist:" + int32ToString(opts.ArtistID)
	case opts.AlbumID != 0:
		return "album:" + int32ToString(opts.AlbumID)
	default:
		return "track:" + int32ToString(opts.TrackID)
	}
}

func int32ToString(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}
