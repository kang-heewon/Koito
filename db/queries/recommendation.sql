-- name: GetTracksToRevisit :many
WITH past_stats AS (
    SELECT 
        track_id,
        COUNT(*) AS past_listen_count,
        MAX(listened_at) AS last_listened_at
    FROM listens l
    WHERE l.listened_at BETWEEN $1 AND $2
    GROUP BY track_id
    HAVING COUNT(*) >= 5
),
recent_listens AS (
    SELECT DISTINCT track_id
    FROM listens
    WHERE listened_at > $2
)
SELECT 
    t.id AS track_id,
    t.title,
    t.release_id,
    r.image AS release_image,
    get_artists_for_track(t.id) AS artists,
    p.past_listen_count,
    p.last_listened_at
FROM past_stats p
JOIN tracks_with_title t ON p.track_id = t.id
JOIN releases r ON t.release_id = r.id
LEFT JOIN recent_listens recent ON p.track_id = recent.track_id
WHERE recent.track_id IS NULL
ORDER BY p.past_listen_count DESC
LIMIT $3;
