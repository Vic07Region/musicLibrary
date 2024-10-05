-- +goose Up
-- +goose StatementBegin
-- Включите расширение "uuid-ossp", если оно еще не установлено
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE songs (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    group_name VARCHAR(255) NOT NULL,
    song_name VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    text TEXT DEFAULT '',
    link VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- SET ID PRIMARY KEY
ALTER TABLE songs ADD PRIMARY KEY (id);

-- Создание индекса для полей group_name, song_name и text
CREATE INDEX songs_group_name_idx ON songs (group_name);
CREATE INDEX songs_song_name_idx ON songs (song_name);
CREATE INDEX songs_text_idx ON songs (text);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  songs;
-- +goose StatementEnd
