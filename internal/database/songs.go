package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt" //nolint:gci
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

var (
	ErrDuplicateKey = fmt.Errorf("duplicate key value violates uniqueness constraint")
)

type Storage interface {
	GetGroupID(ctx context.Context, groupName string) (int64, error)
	CountSongs(ctx context.Context) (int, error)
	GetSongs(ctx context.Context, request GetSongsRequest) ([]Song, error)
	CountVerses(ctx context.Context, SongID int) (int, error)
	GetVerses(ctx context.Context, request GetVersesRequest) ([]VerseSmall, error)
	AddSong(ctx context.Context, request AddSongRequest) (*AddSongResponse, error)
	UpdateSong(ctx context.Context, request UpdateSongRequest) error
	UpdateVerse(ctx context.Context, request UpdateVerseRequest) error
	DeleteSong(ctx context.Context, SongID int) error
}

func ILikeAny(column string, value string) sq.Sqlizer {
	return sq.ILike{column: fmt.Sprintf("%%%s%%", value)}
}

func (q *Queries) GetGroupID(ctx context.Context, groupName string) (int64, error) {
	sqlQuery := sq.Select("group_id").
		From("groups").Where(sq.Eq{"name": groupName}).
		PlaceholderFormat(sq.Dollar)
	var groupID int64
	err := sqlQuery.RunWith(q.db).QueryRowContext(ctx).Scan(&groupID)
	if err != nil {
		if q.debug {
			q.log.Error("database.GetGroupID | QueryRowContext", "error", err.Error())
		}
		return 0, err
	}
	return groupID, nil
}

func (q *Queries) CountSongs(ctx context.Context) (int, error) {
	var countSongs int
	sqlQuery := sq.Select("COUNT(song_id)").From("songs")
	err := sqlQuery.RunWith(q.db).QueryRowContext(ctx).Scan(&countSongs)
	if err != nil {
		if q.debug {
			q.log.Error("database.CountSongs | QueryRowContext", "error", err.Error())
		}
		return 0, err
	}
	return countSongs, nil
}

type GetSongsRequest struct {
	GroupName   *string    `json:"group_name,omitempty" form:"group_name"`
	SongName    *string    `json:"song_name,omitempty" form:"song_name"`
	ReleaseDate *time.Time `json:"release_date,omitempty" form:"release_date"`
	SongText    *string    `json:"song_text,omitempty" form:"song_text"`
	Limit       uint64     `json:"limit" form:"limit"`
	Offset      uint64     `json:"offset" form:"offset"`
}

func (q *Queries) GetSongs(ctx context.Context, request GetSongsRequest) ([]Song, error) {
	sqlQuery := sq.Select("DISTINCT song_id", "name", "song", "releaseDate", "link").
		From("songs").
		InnerJoin("groups USING(group_id)").
		InnerJoin("verses USING(song_id)").PlaceholderFormat(sq.Dollar)

	if request.SongText != nil {
		sqlQuery = sqlQuery.Where(ILikeAny("verse_text", *request.SongText))
	}

	if request.GroupName != nil {
		sqlQuery = sqlQuery.Where(ILikeAny("name", *request.GroupName))
	}

	if request.SongName != nil {
		sqlQuery = sqlQuery.Where(ILikeAny("song", *request.SongName))
	}

	if request.ReleaseDate != nil {
		sqlQuery = sqlQuery.Where(sq.Eq{"releaseDate": *request.ReleaseDate})
	}

	sqlQuery = sqlQuery.OrderBy("releaseDate DESC")

	if request.Limit > 0 && request.Limit <= 100 {
		sqlQuery = sqlQuery.Limit(request.Limit)
	} else {
		sqlQuery = sqlQuery.Limit(10)
	}

	if request.Offset > 0 {
		sqlQuery = sqlQuery.Offset(request.Offset)
	}

	rows, err := sqlQuery.RunWith(q.db).QueryContext(ctx)
	if err != nil {
		if q.debug {
			q.log.Error("database.GetSongs | QueryContext", "error", err.Error())
		}
		//todo заменить на error.wrap +-
		return nil, err
	}
	defer rows.Close()
	var songList []Song
	for rows.Next() {
		var i Song
		if err := rows.Scan(
			&i.SongID,
			&i.GroupName,
			&i.SongName,
			&i.ReleaseDate,
			&i.Link,
		); err != nil {
			if q.debug {
				q.log.Error("database.GetSongs | row.Scan", "error", err.Error())
			}
			return nil, err
		}
		songList = append(songList, i)
	}
	if err := rows.Close(); err != nil {
		if q.debug {
			q.log.Error("database.GetSongs | rows.Close", "error", err.Error())
		}
		return nil, err
	}
	if err := rows.Err(); err != nil {
		if q.debug {
			q.log.Error("database.GetSongs | rows.Err", "error", err.Error())
		}
		return nil, err
	}
	return songList, nil
}

func (q *Queries) CountVerses(ctx context.Context, SongID int) (int, error) {
	var verseCount int
	sqlQuery := sq.Select("COUNT(verse_id)").
		From("verses").Where(sq.Eq{"song_id": SongID}).PlaceholderFormat(sq.Dollar)
	if err := sqlQuery.RunWith(q.db).QueryRowContext(ctx).Scan(&verseCount); err != nil {
		if q.debug {
			q.log.Error("database.CountVerses | QueryRowContext", "error", err.Error())
		}
		return 0, err
	}
	return verseCount, nil
}

type GetVersesRequest struct {
	SongID int `json:"song_id" form:"song_id"`
	Limit  int `json:"limit" form:"limit"`
	Offset int `json:"offset" form:"offset"`
}

func (q *Queries) GetVerses(ctx context.Context, request GetVersesRequest) ([]VerseSmall, error) {
	sqlQuery := sq.Select("verse_number", "verse_text").
		From("verses").
		Where(sq.Eq{"song_id": request.SongID}).
		OrderBy("verse_number").PlaceholderFormat(sq.Dollar)

	if request.Limit > 0 {
		sqlQuery = sqlQuery.Limit(uint64(request.Limit))
	} else {
		sqlQuery = sqlQuery.Limit(2)
	}

	if request.Offset > 0 {
		sqlQuery = sqlQuery.Offset(uint64(request.Offset))
	}

	rows, err := sqlQuery.RunWith(q.db).QueryContext(ctx)
	if err != nil {
		if q.debug {
			q.log.Error("database.GetVerses | QueryContext", "error", err.Error())
		}
		return nil, err
	}
	defer rows.Close()

	var verses []VerseSmall

	for rows.Next() {
		var i VerseSmall
		if err := rows.Scan(
			&i.VerseNumber,
			&i.VerseText,
		); err != nil {
			if q.debug {
				q.log.Error("database.GetVerses | row.Scan", "error", err.Error())
			}
			return nil, err
		}
		verses = append(verses, i)
	}

	if err := rows.Close(); err != nil {
		if q.debug {
			q.log.Error("database.GetVerses | rows.Close", "error", err.Error())
		}
		return nil, err
	}

	if err := rows.Err(); err != nil {
		if q.debug {
			q.log.Error("database.GetVerses | rows.Err", "error", err.Error())
		}
		return nil, err
	}

	return verses, nil
}

type AddSongRequest struct {
	GroupName   string       `json:"group_name"`
	SongName    string       `json:"song_name"`
	ReleaseDate time.Time    `json:"release_date,omitempty"`
	Link        string       `json:"link,omitempty"`
	Verses      []VerseSmall `json:"verses"`
}

type AddSongResponse struct {
	SongID int64 `json:"song_id"`
}

func (q *Queries) AddSong(ctx context.Context, request AddSongRequest) (*AddSongResponse, error) {
	var groupID int64
	var songID int64
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, err := q.db.BeginTx(ctxWithTimeout, nil)
	if err != nil {
		if q.debug {
			q.log.Error("database.AddSong | BeginTx", "error", err.Error())
		}
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	insertGroup := psql.Insert("groups").Columns("name").
		Values(request.GroupName).
		Suffix("ON CONFLICT (name) DO NOTHING RETURNING group_id")

	err = insertGroup.RunWith(tx).QueryRowContext(ctxWithTimeout).Scan(&groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = psql.Select("group_id").From("groups").
				Where(sq.Eq{"name": request.GroupName}).RunWith(tx).QueryRowContext(ctx).Scan(&groupID)
			if err != nil {
				if q.debug {
					q.log.Error("database.AddSong | insertGroup GetGroupID", "error", err.Error())
				}
				return nil, err
			}
		} else {
			if q.debug {
				q.log.Error("database.AddSong | insertGroup.QueryRowContext", "error", err.Error())
			}
			return nil, err
		}
	}

	insertSong := psql.Insert("songs").Columns("group_id", "song", "releaseDate", "link").
		Values(groupID, request.SongName, request.ReleaseDate, request.Link).
		Suffix("RETURNING song_id")

	err = insertSong.RunWith(tx).QueryRowContext(ctx).Scan(&songID)
	if err != nil {
		if q.debug {
			q.log.Error("database.AddSong | insertSong.QueryRowContext", "error", err.Error())
		}
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, ErrDuplicateKey
			}
		}
		return nil, err
	}

	insertVerses := psql.Insert("verses").Columns("song_id", "verse_number", "verse_text")
	for _, verse := range request.Verses {
		insertVerses = insertVerses.Values(songID, verse.VerseNumber, verse.VerseText)
	}
	_, err = insertVerses.RunWith(tx).ExecContext(ctxWithTimeout)
	if err != nil {
		if q.debug {
			q.log.Error("database.AddSong | insertVerses.ExecContext", "error", err.Error())
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		if q.debug {
			q.log.Error("database.AddSong | Commit", "error", err.Error())
		}
		return nil, err
	}
	return &AddSongResponse{SongID: songID}, nil
}

type UpdateSongRequest struct {
	SongID      int        `json:"song_id"`
	GroupID     *int64     `json:"group_id"`
	SongName    *string    `json:"song_name,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Link        *string    `json:"link,omitempty"`
}

func (q *Queries) UpdateSong(ctx context.Context, request UpdateSongRequest) error {
	sqlQury := sq.Update("songs").PlaceholderFormat(sq.Dollar)

	if request.GroupID != nil {
		sqlQury = sqlQury.Set("group_id", *request.GroupID)
	}

	if request.SongName != nil {
		sqlQury = sqlQury.Set("song", *request.SongName)
	}

	if request.ReleaseDate != nil {
		sqlQury = sqlQury.Set("releaseDate", *request.ReleaseDate)
	}

	if request.Link != nil {
		sqlQury = sqlQury.Set("link", *request.Link)
	}

	result, err := sqlQury.Where(sq.Eq{"song_id": request.SongID}).
		RunWith(q.db).ExecContext(ctx)
	if err != nil {
		if q.debug {
			q.log.Error("database.UpdateSong | ExecContext", "error", err.Error())
		}
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		if q.debug {
			q.log.Error("database.UpdateSong | RowsAffected", "error", err.Error())
		}
		return err
	}

	if affected == 0 {
		if q.debug {
			q.log.Warn("database.UpdateSong | RowsAffected 0",
				"error", sql.ErrNoRows.Error(),
				"song_id", request.SongID)
		}
		return sql.ErrNoRows
	}

	return nil

}

type UpdateVerseRequest struct {
	SongID      int    `json:"song_id"`
	VerseNumber int    `json:"verse_number"`
	VerseText   string `json:"verse_text"`
}

func (q *Queries) UpdateVerse(ctx context.Context, request UpdateVerseRequest) error {
	sqlQuery := sq.Update("verses").
		Set("verse_text", request.VerseText).
		Where(sq.Eq{
			"song_id":      request.SongID,
			"verse_number": request.VerseNumber,
		}).PlaceholderFormat(sq.Dollar)

	result, err := sqlQuery.RunWith(q.db).ExecContext(ctx)
	if err != nil {
		if q.debug {
			q.log.Error("database.UpdateVerse | ExecContext", "error", err.Error())
		}
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		if q.debug {
			q.log.Error("database.UpdateVerse | RowsAffected", "error", err.Error())
		}
		return err
	}

	if affected == 0 {
		if q.debug {
			q.log.Warn("database.UpdateVerse | RowsAffected 0",
				"error", sql.ErrNoRows.Error(),
				"song_id", request.SongID)
		}
		return sql.ErrNoRows
	}

	return nil
}

func (q *Queries) DeleteSong(ctx context.Context, SongID int) error {
	sqlQuery := sq.Delete("songs").Where(sq.Eq{"song_id": SongID}).PlaceholderFormat(sq.Dollar)

	result, err := sqlQuery.RunWith(q.db).ExecContext(ctx)
	if err != nil {
		if q.debug {
			q.log.Error("database.DeleteSong | ExecContext", "error", err.Error())
		}
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		if q.debug {
			q.log.Error("database.DeleteSong | RowsAffected", "error", err.Error())
		}
		return err
	}

	if affected == 0 {
		if q.debug {
			q.log.Warn("database.DeleteSong | RowsAffected 0",
				"error", sql.ErrNoRows.Error(),
				"song_id", SongID)
		}
		return sql.ErrNoRows
	}
	return nil
}
