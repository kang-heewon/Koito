-- name: SearchArtists :many
SELECT id, name, musicbrainz_id, image, score
FROM (
    SELECT
        a.id,
        a.name,
        a.musicbrainz_id,
        a.image,
        bigm_similarity(aa.alias, $1) AS score,
        ROW_NUMBER() OVER (PARTITION BY a.id ORDER BY bigm_similarity(aa.alias, $1) DESC) AS rn
    FROM artist_aliases aa
    JOIN artists_with_name a ON aa.artist_id = a.id
    WHERE bigm_similarity(aa.alias, $1) > 0.22
) ranked
WHERE rn = 1
ORDER BY score DESC
LIMIT $2;

-- name: SearchArtistsBySubstring :many
SELECT id, name, musicbrainz_id, image, score
FROM (
    SELECT
        a.id,
        a.name,
        a.musicbrainz_id,
        a.image,
        1.0 AS score, -- why
        ROW_NUMBER() OVER (PARTITION BY a.id ORDER BY aa.alias) AS rn
    FROM artist_aliases aa
    JOIN artists_with_name a ON aa.artist_id = a.id
    WHERE aa.alias ILIKE $1 || '%'
) ranked
WHERE rn = 1
ORDER BY score DESC
LIMIT $2;

-- name: SearchTracks :many
SELECT
    ranked.id,
    ranked.title,
    ranked.musicbrainz_id,
    ranked.release_id,
    ranked.image,
    ranked.score,
    artists.artists
FROM (
    SELECT
        t.id,
        t.title,
        t.musicbrainz_id,
        t.release_id,
        r.image,
        bigm_similarity(ta.alias, $1) AS score,
        ROW_NUMBER() OVER (PARTITION BY t.id ORDER BY bigm_similarity(ta.alias, $1) DESC) AS rn
    FROM track_aliases ta
    JOIN tracks_with_title t ON ta.track_id = t.id
    JOIN releases r ON t.release_id = r.id
    WHERE bigm_similarity(ta.alias, $1) > 0.22
) ranked
CROSS JOIN LATERAL (
    SELECT json_agg(
        jsonb_build_object('id', a.id, 'name', a.name)
        ORDER BY at.is_primary DESC, a.name
    ) AS artists
    FROM artist_tracks at
    JOIN artists_with_name a ON a.id = at.artist_id
    WHERE at.track_id = ranked.id
) artists
WHERE ranked.rn = 1
ORDER BY ranked.score DESC, ranked.title
LIMIT $2;

-- name: SearchTracksBySubstring :many
SELECT
    ranked.id,
    ranked.title,
    ranked.musicbrainz_id,
    ranked.release_id,
    ranked.image,
    ranked.score,
    artists.artists
FROM (
    SELECT
        t.id,
        t.title,
        t.musicbrainz_id,
        t.release_id,
        r.image,
        1.0 AS score,
        ROW_NUMBER() OVER (PARTITION BY t.id ORDER BY ta.alias) AS rn
    FROM track_aliases ta
    JOIN tracks_with_title t ON ta.track_id = t.id
    JOIN releases r ON t.release_id = r.id
    WHERE ta.alias ILIKE $1 || '%'
) ranked
CROSS JOIN LATERAL (
    SELECT json_agg(
        jsonb_build_object('id', a.id, 'name', a.name)
        ORDER BY at.is_primary DESC, a.name
    ) AS artists
    FROM artist_tracks at
    JOIN artists_with_name a ON a.id = at.artist_id
    WHERE at.track_id = ranked.id
) artists
WHERE ranked.rn = 1
ORDER BY ranked.score DESC, ranked.title
LIMIT $2;

-- name: SearchReleases :many
SELECT
    ranked.id,
    ranked.title,
    ranked.musicbrainz_id,
    ranked.image,
    ranked.various_artists,
    ranked.score,
    artists.artists
FROM (
    SELECT
        r.id,
        r.title,
        r.musicbrainz_id,
        r.image,
        r.various_artists,
        bigm_similarity(ra.alias, $1) AS score,
        ROW_NUMBER() OVER (PARTITION BY r.id ORDER BY bigm_similarity(ra.alias, $1) DESC) AS rn
    FROM release_aliases ra
    JOIN releases_with_title r ON ra.release_id = r.id
    WHERE bigm_similarity(ra.alias, $1) > 0.22
) ranked
CROSS JOIN LATERAL (
    SELECT json_agg(
        jsonb_build_object('id', a.id, 'name', a.name)
        ORDER BY ar.is_primary DESC, a.name
    ) AS artists
    FROM artist_releases ar
    JOIN artists_with_name a ON a.id = ar.artist_id
    WHERE ar.release_id = ranked.id
) artists
WHERE ranked.rn = 1
ORDER BY ranked.score DESC, ranked.title
LIMIT $2;

-- name: SearchReleasesBySubstring :many
SELECT
    ranked.id,
    ranked.title,
    ranked.musicbrainz_id,
    ranked.image,
    ranked.various_artists,
    ranked.score,
    artists.artists
FROM (
    SELECT
        r.id,
        r.title,
        r.musicbrainz_id,
        r.image,
        r.various_artists,
        1.0 AS score,
        ROW_NUMBER() OVER (PARTITION BY r.id ORDER BY ra.alias) AS rn
    FROM release_aliases ra
    JOIN releases_with_title r ON ra.release_id = r.id
    WHERE ra.alias ILIKE $1 || '%'
) ranked
CROSS JOIN LATERAL (
    SELECT json_agg(
        jsonb_build_object('id', a.id, 'name', a.name)
        ORDER BY ar.is_primary DESC, a.name
    ) AS artists
    FROM artist_releases ar
    JOIN artists_with_name a ON a.id = ar.artist_id
    WHERE ar.release_id = ranked.id
) artists
WHERE ranked.rn = 1
ORDER BY ranked.score DESC, ranked.title
LIMIT $2;
