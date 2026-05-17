package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/model"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (ar *AuthRepository) AddNewUser(ctx context.Context, email, hashedPassword string) (model.User, error) {
	sql := `
		WITH register AS (
			INSERT INTO users
			(email, password)
			VALUES
			($1, $2)
			RETURNING id, email, created_at
		), create_wallet AS (
			INSERT INTO wallets (user_id)
			SELECT id FROM register
		)
		SELECT id, email, created_at FROM register;
	`
	args := []any{email, hashedPassword}

	var user model.User
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Email, &user.CreatedAt); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ar *AuthRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	sql := `
		SELECT id, password, is_verified
		FROM users
		WHERE email = $1;
	`
	args := []any{email}

	var user model.User
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Password, &user.IsVerified); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ar *AuthRepository) SaveToken(ctx context.Context, tokenHash string, userId int, expiresAt time.Time) error {
	sql := `
		INSERT INTO active_tokens (token_hash, user_id, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (token_hash)
		DO UPDATE SET user_id = EXCLUDED.user_id, expires_at = EXCLUDED.expires_at;
	`
	args := []any{tokenHash, userId, expiresAt}

	_, err := ar.db.Exec(ctx, sql, args...)
	return err
}

func (ar *AuthRepository) DeleteToken(ctx context.Context, tokenHash string) error {
	sql := `
		DELETE FROM active_tokens
		WHERE token_hash = $1;
	`
	args := []any{tokenHash}

	_, err := ar.db.Exec(ctx, sql, args...)
	return err
}

func (ar *AuthRepository) IsTokenActive(ctx context.Context, tokenHash string) (bool, error) {
	sql := `
		SELECT EXISTS (
			SELECT 1
			FROM active_tokens
			WHERE token_hash = $1
				AND expires_at > NOW()
		);
	`
	args := []any{tokenHash}

	var isActive bool
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&isActive); err != nil {
		return false, err
	}

	return isActive, nil
}
