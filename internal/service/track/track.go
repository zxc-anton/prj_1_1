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
	rows, err := t.pool.Query(ctx, `SELECT a.name, a.photo_url, a."ID"  FROM track_artist ta JOIN author a ON a."ID" = ta."author_ID"  WHERE "track_ID" = $1 `, id)

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
			t."album_ID", 
			t.title, 
			t.created_at,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', a."ID",
					'name', a.name,
					'photo_url', a.photo_url
				)
			)
			 AS authors
		FROM track t
		JOIN track_artist ta ON t."ID" = ta."track_ID"
		JOIN author a ON  ta."author_ID" = a."ID"
		WHERE t."ID" = $1
		GROUP BY t."ID", t.file_url, t."album_ID", t.title, t.created_at`,
		id,
	)

	track := Track{Id: id}
	authorsBytes := make([]byte, 0)

	if err := row.Scan(&track.FileURL, &track.AlbumId, &track.Title, &track.CreatedAt, &authorsBytes); err != nil {
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
			t."ID", 
			t.file_url, 
			t."album_ID", 
			t.title, 
			t.created_at,
			COALESCE(
				(
					SELECT JSON_AGG(
						JSON_BUILD_OBJECT(
							'id', a."ID",
							'name', a.name,
							'photo_url', a.photo_url
						) ORDER BY a.name
					)
					FROM track_artist ta
					JOIN author a ON ta."author_ID" = a."ID"
					WHERE ta."track_ID" = t."ID"
				),
				'[]'::json
			) AS authors
		FROM track t 
		LEFT JOIN listen_history lh ON t."ID" = lh."track_ID"
		GROUP BY t."ID", t.file_url, t."album_ID", t.title, t.created_at 
		ORDER BY COUNT(lh."track_ID") DESC
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
			t."ID",
			t.file_url, 
			t."album_ID", 
			t.title, 
			t.created_at,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', a."ID",
					'name', a.name,
					'photo_url', a.photo_url
				)
			)
			 AS authors
		FROM track t
		JOIN track_artist ta ON t."ID" = ta."track_ID"
		JOIN author a ON  ta."author_ID" = a."ID"
		WHERE t."ID" = ANY($1)
		GROUP BY t."ID", t.file_url, t."album_ID", t.title, t.created_at`,
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
