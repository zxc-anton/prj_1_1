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
		WHERE "track_ID" = ANY(SELECT "track_ID" FROM track_artist WHERE "author_ID" = $1)`,
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
		`SELECT name, photo_url FROM author WHERE "ID" = $1`,
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
			SELECT DISTINCT aa."album_ID"
		FROM album_author aa
		WHERE aa."author_ID" = $1
		ORDER BY aa."album_ID"
		OFFSET $2
		LIMIT $3
		)
		SELECT
			a."ID",
			a.title,
			a.cover_url,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', at."ID",
					'name', at.name,
					'photo_url', at.photo_url
				)
			)
			 AS authors
		FROM author_albums aa
		JOIN album a ON a."ID" = aa."album_ID"
		JOIN author at ON aa."author_ID" = at."ID"
		WHERE at."ID" = $1
		GROUP BY a."ID", a.title, a.cover_url
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
			t."ID"
		FROM track t
		JOIN track_artist ta ON t."ID" = ta."track_ID"
		WHERE ta."author_ID" = $1
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
