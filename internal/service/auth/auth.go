package auth

import (
	"context"
	"errors"
	"server/internal/service/auth/types"
	"server/pkg/token"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	pool         *pgxpool.Pool
	tokenService string
}

func New(pool *pgxpool.Pool) *AuthService {
	return &AuthService{
		pool: pool,
	}
}

func (a *AuthService) RegisteUser(ctx context.Context, user types.RequestUser) error {
	hashPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Commit(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO "user" (name, email, password_hash, is_verification) VALUES ($1, $2, $3, $4)`, user.Name, user.Email, hashPassword, false)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return nil
}

func (a *AuthService) LoginUser(ctx context.Context, params types.LoginUser) (string, string, error) {
	row := a.pool.QueryRow(ctx, `SELECT password_hash, "ID" FROM "user" WHERE email = $1`, params.Email)
	var (
		password string
		id       int
	)
	if err := row.Scan(&password, &id); err != nil {
		return "", "", err
	}
	if !comparePasswords(params.Password, password) {
		return "", "", errors.New("not valid")
	}

	accessToken, refreshToken, err := token.CreateTokens(id)
	if err != nil {
		return "", "", nil
	}
	return accessToken, refreshToken, nil
}

func (a *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	id, err := token.DecodeToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	return token.CreateTokens(id)
}
