-- +goose Up
-- +goose StatementBegin

CREATE TABLE genres (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY (
        SEQUENCE NAME genres_id_seq
        START WITH 1
        INCREMENT BY 1
        NO MINVALUE
        NO MAXVALUE
        CACHE 1
    ),
    name text UNIQUE NOT NULL,
    CONSTRAINT genres_pkey PRIMARY KEY (id)
);

CREATE TABLE release_genres (
    release_id integer NOT NULL,
    genre_id integer NOT NULL,
    CONSTRAINT release_genres_pkey PRIMARY KEY (release_id, genre_id)
);

CREATE TABLE artist_genres (
    artist_id integer NOT NULL,
    genre_id integer NOT NULL,
    CONSTRAINT artist_genres_pkey PRIMARY KEY (artist_id, genre_id)
);

ALTER TABLE ONLY release_genres
    ADD CONSTRAINT release_genres_release_id_fkey FOREIGN KEY (release_id) REFERENCES releases(id) ON DELETE CASCADE;

ALTER TABLE ONLY release_genres
    ADD CONSTRAINT release_genres_genre_id_fkey FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE;

ALTER TABLE ONLY artist_genres
    ADD CONSTRAINT artist_genres_artist_id_fkey FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE;

ALTER TABLE ONLY artist_genres
    ADD CONSTRAINT artist_genres_genre_id_fkey FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE;

CREATE INDEX idx_release_genres_release_id ON release_genres USING btree (release_id);
CREATE INDEX idx_artist_genres_artist_id ON artist_genres USING btree (artist_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS release_genres CASCADE;
DROP TABLE IF EXISTS artist_genres CASCADE;
DROP TABLE IF EXISTS genres CASCADE;

-- +goose StatementEnd
