-- +goose Up
-- +goose StatementBegin
-- Table: Groups
CREATE TABLE groups (
    group_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Table: Songs
CREATE TABLE songs (
    song_id SERIAL PRIMARY KEY,
    group_id INT NOT NULL,
    song VARCHAR(255) NOT NULL,
    releaseDate DATE,
    link VARCHAR(255),
    FOREIGN KEY (group_id) REFERENCES Groups(group_id) ON DELETE CASCADE,
    CONSTRAINT unique_group_song UNIQUE (group_id, song)
);

-- Table: Verses
CREATE TABLE verses (
    verse_id SERIAL PRIMARY KEY,
    song_id INT NOT NULL,
    verse_number INT NOT NULL,
    verse_text TEXT,
    FOREIGN KEY (song_id) REFERENCES Songs(song_id) ON DELETE CASCADE,
    CONSTRAINT unique_song_verse UNIQUE (song_id, verse_number)
);

-- Indexes
CREATE INDEX idx_songs_group_id ON Songs(group_id);
CREATE INDEX idx_songs_song ON Songs(song);
CREATE INDEX idx_songs_releaseDate ON Songs(releaseDate);
CREATE INDEX idx_verses_song_id ON Verses(song_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  groups;
DROP TABLE  songs;
DROP TABLE  verses;
-- +goose StatementEnd
