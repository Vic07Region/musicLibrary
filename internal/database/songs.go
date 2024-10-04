package database

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"time"
)

func ILikeAny(column string, value string) squirrel.Sqlizer {
	return squirrel.ILike{column: fmt.Sprintf("%%%s%%", value)}
}
func LikeAny(column string, value string) squirrel.Sqlizer {
	return squirrel.Like{column: fmt.Sprintf("%%%s%%", value)}
}

func (q *Queries) GetSong(ctx context.Context, song_id uuid.UUID) (*Song, error) {
	sql_query := squirrel.Select("*").From("songs").
		Where(squirrel.Eq{"id": song_id.String()}).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)

	row := sql_query.RunWith(q.db).QueryRowContext(ctx)

	var song Song

	if err := row.Scan(
		&song.ID,
		&song.GroupName,
		&song.SongName,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
		&song.CreatedAt,
		&song.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &song, nil
}

type GetSongsParam struct {
	Group_name string
	Song_name  string
	Text       string
	Limit      int64
	Offset     int64
}

func (q *Queries) GetSongs(ctx context.Context, params GetSongsParam) ([]Song, error) {
	sql_query := squirrel.Select("*").From("songs")

	sql_query = sql_query.PlaceholderFormat(squirrel.Dollar)

	if params.Group_name != "" {
		sql_query = sql_query.Where(ILikeAny("group_name", params.Group_name))
	}

	if params.Song_name != "" {
		sql_query = sql_query.Where(ILikeAny("song_name", params.Song_name))
	}

	if params.Text != "" {
		sql_query = sql_query.Where(ILikeAny("text", params.Text))
	}

	if params.Offset > 0 {
		sql_query = sql_query.Offset(uint64(params.Offset))
	}

	if params.Limit > 0 {
		sql_query = sql_query.Limit(uint64(params.Limit))
	} else {
		sql_query = sql_query.Limit(10)
	}

	rows, err := sql_query.RunWith(q.db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []Song

	for rows.Next() {
		var i Song
		if err := rows.Scan(
			&i.ID,
			&i.GroupName,
			&i.SongName,
			&i.ReleaseDate,
			&i.Text,
			&i.Link,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		songs = append(songs, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return songs, nil
}

type GetSongTextParam struct {
	Song_id uuid.UUID
	Limit   int64
	Offset  int64
}

func (q *Queries) GetSongText(ctx context.Context, song_id uuid.UUID) (string, error) {
	sql_query := squirrel.Select("text").From("songs").
		Where(squirrel.Eq{"id": song_id.String()})

	sql_query = sql_query.PlaceholderFormat(squirrel.Dollar)

	row := sql_query.RunWith(q.db).QueryRowContext(ctx)

	var song_text string

	if err := row.Scan(&song_text); err != nil {
		return "", nil
	}
	return song_text, nil
}

func (q *Queries) DeleteSong(ctx context.Context, id uuid.UUID) error {
	sql_query := squirrel.Delete("songs").Where(squirrel.Eq{"id": id.String()})

	sql_query = sql_query.PlaceholderFormat(squirrel.Dollar)

	_, err := sql_query.RunWith(q.db).ExecContext(ctx)

	return err
}

type UpdateSongParam struct {
	Song_id     uuid.UUID
	GroupName   string
	SongName    string
	ReleaseDate time.Time
	Text        string
	Link        string
}

func (q *Queries) UpdateSong(ctx context.Context, params UpdateSongParam) error {
	sql_query := squirrel.Update("songs")

	sql_query = sql_query.PlaceholderFormat(squirrel.Dollar)

	if params.GroupName != "" {
		sql_query = sql_query.Set("group_name", params.GroupName)
	}

	if params.SongName != "" {
		sql_query = sql_query.Set("song_name", params.SongName)
	}

	if params.ReleaseDate.IsZero() {
		sql_query = sql_query.Set("release_date", params.ReleaseDate)
	}

	if params.Text != "" {
		sql_query = sql_query.Set("text", params.Text)
	}

	if params.Link != "" {
		sql_query = sql_query.Set("link", params.Link)
	}
	sql_query = sql_query.Set("updated_at", time.Now()).Where(squirrel.Eq{"id": params.Song_id.String()})
	_, err := sql_query.RunWith(q.db).ExecContext(ctx)

	return err
}

type CreateSongParam struct {
	GroupName   string
	SongName    string
	ReleaseDate time.Time
	Text        string
	Link        string
}

type CreateSongResult struct {
	ID        uuid.UUID
	GroupName string
	SongName  string
	CreatedAt time.Time
}

func (q *Queries) CreateSong(ctx context.Context, params CreateSongParam) (*CreateSongResult, error) {
	sql_query := squirrel.Insert("songs").
		Columns("group_name", "song_name", "release_date", "text", "link").
		Values(
			params.GroupName,
			params.SongName,
			params.ReleaseDate,
			params.Text,
			params.Link,
		).Suffix("RETURNING id, group_name, song_name, created_at").
		PlaceholderFormat(squirrel.Dollar)

	row := sql_query.RunWith(q.db).QueryRowContext(ctx)

	var song CreateSongResult

	if err := row.Scan(&song.ID, &song.GroupName, &song.SongName, &song.CreatedAt); err != nil {
		return nil, err
	}

	return &song, nil
}
