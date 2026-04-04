package refreshtoken_repository

import (
	"context"
	"time"

	"gaudiot.com/fonli/core/database"

	"github.com/jackc/pgx/v5"
)

type pgxRefreshTokenRepository struct {
	db *database.DB
}

func NewPgxRefreshTokenRepository(db *database.DB) RefreshTokenRepository {
	return &pgxRefreshTokenRepository{db: db}
}

func (r *pgxRefreshTokenRepository) CreateRefreshToken(token, userID string, expiresAt time.Time) (*RefreshToken, error) {
	ctx := context.Background()
	query := `
		INSERT INTO refresh_tokens (token, user_id, expires_at, is_valid)
		VALUES ($1, $2, $3, $4)
		RETURNING token, user_id, expires_at, is_valid
	`
	row := r.db.Pool.QueryRow(ctx, query, token, userID, expiresAt, true)

	var rt RefreshToken
	err := row.Scan(&rt.Token, &rt.UserID, &rt.ExpiresAt, &rt.IsValid)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *pgxRefreshTokenRepository) GetByToken(token string) (*RefreshToken, error) {
	ctx := context.Background()
	query := `
		SELECT token, user_id, expires_at, is_valid
		FROM refresh_tokens
		WHERE token = $1
	`
	row := r.db.Pool.QueryRow(ctx, query, token)
	var rt RefreshToken
	err := row.Scan(&rt.Token, &rt.UserID, &rt.ExpiresAt, &rt.IsValid)
	if err == pgx.ErrNoRows {
		return nil, ErrInvalidRefreshToken
	} else if err != nil {
		return nil, err
	}
	if !rt.IsValid {
		return nil, ErrInvalidRefreshToken
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, ErrExpiredRefreshToken
	}
	return &rt, nil
}

func (r *pgxRefreshTokenRepository) InvalidateUserRefreshTokens(userID string) error {
	ctx := context.Background()
	query := `
		UPDATE refresh_tokens
		SET is_valid = FALSE
		WHERE user_id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, userID)
	return err
}

func (r *pgxRefreshTokenRepository) DeleteByToken(token string) error {
	ctx := context.Background()
	query := `
		DELETE FROM refresh_tokens
		WHERE token = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, token)
	return err
}

func (r *pgxRefreshTokenRepository) DeleteAllByUserID(userID string) error {
	ctx := context.Background()
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, userID)
	return err
}
