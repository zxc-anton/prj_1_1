package artist

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArtistService struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *ArtistService {
	return &ArtistService{
		pool: pool,
	}
}

func (a *ArtistService) GetArtistLastMonthListeners(ctx context.Context, id int) (int, error) {
	row := a.pool.QueryRow(
		ctx,
		`SELECT 
			COUNT(*) lastMonthListeners
		FROM listen_history 
		WHERE track_id = ANY(SELECT track_id FROM track_author WHERE author_id = $1)`,
		id,
	)
	var LastMonthListeners int
	if err := row.Scan(&LastMonthListeners); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return LastMonthListeners, nil
}

func (a *ArtistService) GetArtistById(ctx context.Context, id int) (ArtistResponse, error) {
	row := a.pool.QueryRow(
		ctx,
		`SELECT name, photo_url FROM authors WHERE id = $1`,
		id,
	)

	var artistResponse ArtistResponse

	if err := row.Scan(&artistResponse.Artist.Name, &artistResponse.Artist.Cover); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ArtistResponse{}, ErrNotFound
		}
		return ArtistResponse{}, err
	}

	lastMonthListeners, err := a.GetArtistLastMonthListeners(ctx, id)

	if err != nil {
		return ArtistResponse{}, err
	}

	artistResponse.Artist.Id = id
	artistResponse.LastMonthListeners = lastMonthListeners

	return artistResponse, nil
}

func (a *ArtistService) GetArtistAlbums(ctx context.Context, id, limit, offset int) ([]Album, error) {
	rows, err := a.pool.Query(
		ctx,
		`
		WITH author_albums AS (
			SELECT DISTINCT aa.album_id
			FROM album_author aa
			WHERE aa.author_id = $1
			ORDER BY aa.album_id
			OFFSET $2
			LIMIT $3
		)
		SELECT
			a.id,
			a.title,
			a.cover_url,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', at.id,
					'name', at.name,
					'photo_url', at.photo_url
				)
			)
			 AS authors
		FROM album_author aa
		JOIN albums a ON a.id = aa.album_id
		JOIN authors at ON aa.author_id = at.id
		WHERE at.id = $1
		GROUP BY a.id, a.title, a.cover_url
		OFFSET $2
		LIMIT $3`,
		id,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	albums := make([]Album, 0)
	authorBytes := make([]byte, 0)

	for rows.Next() {
		var album Album

		if err := rows.Scan(&album.Id, &album.Title, &album.Cover, &authorBytes); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(authorBytes, &album.Artists); err != nil {
			return nil, err
		}

		albums = append(albums, album)

		clear(authorBytes)
	}

	return albums, nil
}

func (a *ArtistService) GetArtistPopularTracksIDS(ctx context.Context, id, limit, offset int) ([]int, error) {
	rows, err := a.pool.Query(
		ctx,
		`
		SELECT
			t.id
		FROM tracks t
		JOIN track_author ta ON t.id = ta.track_id
		WHERE ta.author_id = $1
		OFFSET $2
		LIMIT $3
		`,
		id,
		offset,
		limit,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ids := make([]int, 0)

	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
