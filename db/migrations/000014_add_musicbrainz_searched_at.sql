-- +goose Up
-- +goose StatementBegin

ALTER TABLE releases ADD COLUMN musicbrainz_searched_at TIMESTAMPTZ DEFAULT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE releases DROP COLUMN musicbrainz_searched_at;

-- +goose StatementEnd
