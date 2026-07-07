package track

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackService struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *TrackService {
	return &TrackService{
		pool: pool,
	}
}

func (t *TrackService) GetTrackAuthors(ctx context.Context, id int) ([]Author, error) {
	rows, err := t.pool.Query(ctx, `SELECT a.name, a.photo_url, a.id  FROM track_author ta JOIN authors a ON a.id = ta.author_id  WHERE track_id = $1`, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	authors := make([]Author, 0)

	for rows.Next() {
		var author Author
		if err := rows.Scan(&author.Name, &author.PhotoUrl, &author.Id); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil

}

func (t *TrackService) GetTrackById(ctx context.Context, id int) (Track, error) {
	row := t.pool.QueryRow(
		ctx,
		`SELECT 
			t.file_url, 
			t.title, 
			t.created_at,
			t.length,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', a.id,
					'name', a.name,
					'photo_url', a.photo_url
				)
			)
			 AS authors
		FROM tracks t
		JOIN track_author ta ON t.id = ta.track_id
		JOIN authors a ON  ta.author_id = a.id
		WHERE t.id = $1
		GROUP BY t.id, t.file_url, t.title, t.created_at`,
		id,
	)

	track := Track{Id: id}
	authorsBytes := make([]byte, 0)

	if err := row.Scan(&track.FileURL, &track.Title, &track.CreatedAt, &track.Length, &authorsBytes); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Track{}, ErrNotFound
		}
		return Track{}, err
	}
	if err := json.Unmarshal(authorsBytes, &track.Authors); err != nil {
		return Track{}, err
	}
	return track, nil
}

// todo: paralel reques

func (t *TrackService) GetPopularTracks(ctx context.Context, limit, offset int) ([]Track, error) {
	rows, err := t.pool.Query(
		ctx,
		`SELECT 
			t.id, 
			t.file_url, 
			at.album_id, 
			t.title, 
			t.created_at,
			COALESCE(
				(
					SELECT JSON_AGG(
						JSON_BUILD_OBJECT(
							'id', a.id,
							'name', a.name,
							'photo_url', a.photo_url
						) ORDER BY a.name
					)
					FROM track_author ta
					JOIN authors a ON ta.author_id = a.id
					WHERE ta.track_id = t.id
				),
				'[]'::json
			) AS authors
		FROM tracks t 
		LEFT JOIN listen_history lh ON t.id = lh.track_id
		LEFT JOIN album_track at ON at.track_id = t.id
		GROUP BY t.id, t.file_url, at.album_id, t.title, t.created_at 
		ORDER BY COUNT(lh.track_id) DESC
		OFFSET $1 
		LIMIT $2`,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tracks := make([]Track, 0, limit)
	authorsBytes := make([]byte, 0)

	for rows.Next() {
		var track Track

		if err := rows.Scan(&track.Id, &track.FileURL, &track.AlbumId, &track.Title, &track.CreatedAt, &authorsBytes); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(authorsBytes, &track.Authors); err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
		clear(authorsBytes)
	}

	return tracks, nil
}

func (t *TrackService) GetTracksByIDS(ctx context.Context, ids []int) ([]Track, error) {
	rows, err := t.pool.Query(
		ctx,
		`SELECT 
			t.id,
			t.file_url, 
			at.album_id, 
			t.title, 
			t.created_at,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', a.id,
					'name', a.name,
					'photo_url', a.photo_url
				)
			)
			 AS authors
		FROM tracks t
		JOIN track_author ta ON t.id = ta.track_id
		JOIN authors a ON  ta.author_id = a.id
		LEFT JOIN album_track at ON at.track_id = t.id
		WHERE t.id = ANY($1)
		GROUP BY t.id, at.album_id, t.file_url, t.title, t.created_at`,
		ids,
	)

	if err != nil {
		return nil, err
	}

	tracks := make([]Track, 0, len(ids))
	authorsBytes := make([]byte, 0)

	for rows.Next() {
		var track Track

		if err := rows.Scan(&track.Id, &track.FileURL, &track.AlbumId, &track.Title, &track.CreatedAt, &authorsBytes); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(authorsBytes, &track.Authors); err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
		clear(authorsBytes)
	}

	return tracks, nil
}

func (t *TrackService) IsFollowed(ctx context.Context, trackId, userId int) (res bool, err error) {
	row := t.pool.QueryRow(ctx, `SELECT EXISTS (SELECT * FROM favorites WHERE user_id=$1 AND track_id=$2)`, userId, trackId)
	err = row.Scan(&res)
	return
}
