-- +goose Up
-- +goose StatementBegin

DROP VIEW IF EXISTS releases_with_title;
CREATE VIEW releases_with_title AS
SELECT r.id,
   r.musicbrainz_id,
   r.image,
   r.various_artists,
   r.image_source,
   r.musicbrainz_searched_at,
   ra.alias AS title
FROM releases r
JOIN release_aliases ra ON ra.release_id = r.id
WHERE ra.is_primary = true;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP VIEW IF EXISTS releases_with_title;
CREATE VIEW releases_with_title AS
SELECT r.id,
   r.musicbrainz_id,
   r.image,
   r.various_artists,
   r.image_source,
   ra.alias AS title
FROM releases r
JOIN release_aliases ra ON ra.release_id = r.id
WHERE ra.is_primary = true;


-- +goose StatementEnd