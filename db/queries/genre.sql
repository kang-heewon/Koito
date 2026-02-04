-- name: InsertGenre :one
INSERT INTO genres (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

-- name: GetGenreByName :one
SELECT * FROM genres WHERE name = $1 LIMIT 1;

-- name: GetGenresByNames :many
SELECT * FROM genres WHERE name = ANY($1::text[]);

-- name: GetGenresForRelease :many
SELECT g.*
FROM genres g
JOIN release_genres rg ON g.id = rg.genre_id
WHERE rg.release_id = $1;

-- name: GetGenresForArtist :many
SELECT g.*
FROM genres g
JOIN artist_genres ag ON g.id = ag.genre_id
WHERE ag.artist_id = $1;

-- name: AssociateGenreToRelease :exec
INSERT INTO release_genres (release_id, genre_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: AssociateGenreToArtist :exec
INSERT INTO artist_genres (artist_id, genre_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: DeleteReleaseGenres :exec
DELETE FROM release_genres WHERE release_id = $1;

-- name: DeleteArtistGenres :exec
DELETE FROM artist_genres WHERE artist_id = $1;

-- name: GetReleasesWithoutGenres :many
SELECT r.id, r.musicbrainz_id
FROM releases r
WHERE r.musicbrainz_id IS NOT NULL
  AND r.id > $2
  AND NOT EXISTS (
    SELECT 1 FROM release_genres rg WHERE rg.release_id = r.id
  )
ORDER BY r.id ASC
LIMIT $1;

-- name: GetArtistsWithoutGenres :many
SELECT a.id, a.musicbrainz_id
FROM artists a
WHERE a.musicbrainz_id IS NOT NULL
  AND a.id > $2
  AND NOT EXISTS (
    SELECT 1 FROM artist_genres ag WHERE ag.artist_id = a.id
  )
ORDER BY a.id ASC
LIMIT $1;

-- name: GetGenreStatsByListenCount :many
SELECT
    g.name,
    COUNT(l.listened_at) AS listen_count
FROM listens l
JOIN tracks t ON l.track_id = t.id
JOIN releases r ON t.release_id = r.id
JOIN release_genres rg ON r.id = rg.release_id
JOIN genres g ON rg.genre_id = g.id
WHERE l.listened_at BETWEEN $1 AND $2
GROUP BY g.name
ORDER BY listen_count DESC;

-- name: GetGenreStatsByTimeListened :many
SELECT
    g.name,
    COALESCE(SUM(t.duration), 0)::BIGINT AS seconds_listened
FROM listens l
JOIN tracks t ON l.track_id = t.id
JOIN releases r ON t.release_id = r.id
JOIN release_genres rg ON r.id = rg.release_id
JOIN genres g ON rg.genre_id = g.id
WHERE l.listened_at BETWEEN $1 AND $2
GROUP BY g.name
ORDER BY seconds_listened DESC;
