package album

import (
	"context"
	"encoding/json"
	"errors"
	"server/internal/service/track"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlbumService struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *AlbumService {
	return &AlbumService{pool: pool}
}

func (a *AlbumService) GetAlbumById(ctx context.Context, id int) (Album, error) {
	row := a.pool.QueryRow(
		ctx,
		`
		SELECT
			a.id,
			a.title,
			a.cover_url,
			a.created_at,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', at.id,
					'name', at.name,
					'cover', at.photo_url
				)
			)
			 AS authors
		FROM albums a
		JOIN album_author aa ON a.id = aa.album_id
		JOIN authors at ON aa.author_id = at.id
		WHERE a.id = $1
		GROUP BY a.id, a.title, a.cover_url, a.created_at
		`,
		id,
	)
	var album Album
	var authors []byte
	if err := row.Scan(&album.Id, &album.Title, &album.Cover, &album.CreatedAt, &authors); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Album{}, track.ErrNotFound
		}
		return Album{}, err
	}

	if err := json.Unmarshal(authors, &album.Authors); err != nil {
		return Album{}, err
	}
	return album, nil
}

func (a *AlbumService) GetAlbumTracks(ctx context.Context, id int) ([]int, error) {
	rows, err := a.pool.Query(
		ctx,
		`
		SELECT
			track_id
		FROM album_track
		WHERE album_id = $1
		`,
		id,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, track.ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
