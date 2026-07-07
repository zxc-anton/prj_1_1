package user

import (
	"context"
	"server/internal/service/user/types"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *UserService {
	return &UserService{
		pool: pool,
	}
}

func (u *UserService) GetCurrentUser(ctx context.Context, id int) (user types.UserResponse, err error) {
	row := u.pool.QueryRow(ctx, "SELECT avatar_url, name FROM users WHERE id = $1", id)
	err = row.Scan(&user.Avatar_URL, &user.Name)
	user.Id = id
	return
}

func (u *UserService) GetMyPlaylists(ctx context.Context, id, limit, offset int) (types.MyPLaylists, error) {
	rows, err := u.pool.Query(ctx, "SELECT name, cover_url, id, is_public FROM playlists WHERE user_id = $1 OFFSET $2 LIMIT $3", id, offset, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	playlists := make(types.MyPLaylists, 0)

	for rows.Next() {
		var playlist types.Playlist
		if err := rows.Scan(&playlist.Name, &playlist.Cover_URL, &playlist.Id, &playlist.Is_public); err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (u *UserService) GetMyListenHistory(ctx context.Context, id, limit, offset int) (types.ListenHistory, error) {
	rows, err := u.pool.Query(
		ctx,
		`
		SELECT 
			t.id, t.file_url, t.title, a.id album_id, a.title album
		FROM listen_history lh 
		JOIN tracks t ON lh.track_id = t.id
		JOIN album_track at ON at.track_id = t.id
 		RIGHT JOIN albums a ON a.id = at.album_id
		WHERE lh.user_id = $1 
		OFFSET $2 
		LIMIT $3
		`,
		id, offset, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tracks := make(types.ListenHistory, 0)

	for rows.Next() {
		var track types.Track
		if err = rows.Scan(&track.Id, &track.File_url, &track.Title, &track.Album_Id, &track.Album); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}

func (u *UserService) AddInFavorits(ctx context.Context, id, trackId int) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO favorites (user_id, track_id) VALUES ($1, $2)`, id, trackId)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}

func (u *UserService) CreatePlaylist(ctx context.Context, id int, playlist types.CreatePlaylist) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)
	_, err = tx.Exec(ctx, `INSERT INTO playlist (name, is_public, user_id) VALUES ($1, $2, $3)`, playlist.Name, playlist.IsPublic, id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return nil
}

func (u *UserService) DeletePLaylist(ctx context.Context, id int, playlistId int) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM playlist WHERE user_id = $1 AND id = $2`, id, playlistId)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}

func (u *UserService) DeleteTrackFromFavorits(ctx context.Context, id, track_id int) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)
	_, err = tx.Exec(ctx, `DELETE FROM favorites WHERE user_id = $1 AND track_id = $2`, id, track_id)

	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return nil
}

func (u *UserService) GetFavorits(ctx context.Context, id, limit, offset int) ([]int, error) {
	rows, err := u.pool.Query(
		ctx,
		`
		SELECT
			track_id
		FROM favorites
		WHERE user_id = $1
		`,
		id,
	)
	if err != nil {
		return nil, err
	}

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
