-- +goose Up
DROP INDEX IF EXISTS idx_artist_aliases_alias_trgm;
DROP INDEX IF EXISTS idx_release_aliases_alias_trgm;
DROP INDEX IF EXISTS idx_track_aliases_alias_trgm;

DROP EXTENSION IF EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS pg_bigm WITH SCHEMA public;

CREATE INDEX IF NOT EXISTS idx_artist_aliases_alias_bigm ON artist_aliases USING gin (alias gin_bigm_ops);
CREATE INDEX IF NOT EXISTS idx_release_aliases_alias_bigm ON release_aliases USING gin (alias gin_bigm_ops);
CREATE INDEX IF NOT EXISTS idx_track_aliases_alias_bigm ON track_aliases USING gin (alias gin_bigm_ops);

-- +goose Down
DROP INDEX IF EXISTS idx_artist_aliases_alias_bigm;
DROP INDEX IF EXISTS idx_release_aliases_alias_bigm;
DROP INDEX IF EXISTS idx_track_aliases_alias_bigm;

DROP EXTENSION IF EXISTS pg_bigm;
CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;

CREATE INDEX IF NOT EXISTS idx_artist_aliases_alias_trgm ON artist_aliases USING gin (alias gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_release_aliases_alias_trgm ON release_aliases USING gin (alias gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_track_aliases_alias_trgm ON track_aliases USING gin (alias gin_trgm_ops);
