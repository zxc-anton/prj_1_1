package playlist

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PLaylistService struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PLaylistService {
	return &PLaylistService{
		pool: pool,
	}
}

func (p *PLaylistService) GetPLaylistTracks(ctx context.Context, id int, limit, offset int) ([]int, error) {
	rows, err := p.pool.Query(
		ctx,
		`
		SELECT 
			t.id
		FROM playlist_track pt
		JOIN tracks t ON pt.track_id = t.id
		WHERE pt.playlist_id = $1
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

	ids := make([]int, 0, limit)

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}
