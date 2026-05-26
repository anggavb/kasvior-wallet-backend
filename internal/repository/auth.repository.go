package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
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
		SELECT email, created_at FROM register;
	`
	args := []any{email, hashedPassword}

	var user model.User
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&user.Email, &user.CreatedAt); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ar *AuthRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	sql := `
		SELECT id, email, password, is_verified, pin
		FROM users
		WHERE email = $1;
	`
	args := []any{email}

	var user model.User
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Email, &user.Password, &user.IsVerified, &user.Pin); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ar *AuthRepository) GetPasswordResetUserByEmail(ctx context.Context, email string) (model.User, error) {
	sql := `
		SELECT id, email
		FROM users
		WHERE email = $1;
	`
	args := []any{email}

	var user model.User
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Email); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ar *AuthRepository) SavePasswordResetToken(ctx context.Context, userId int, tokenHash string, expiresAt time.Time) error {
	sql := `
		INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3);
	`
	args := []any{userId, tokenHash, expiresAt}

	_, err := ar.db.Exec(ctx, sql, args...)
	return err
}

func (ar *AuthRepository) GetActivePasswordResetToken(ctx context.Context, tokenHash string) (model.PasswordResetToken, error) {
	sql := `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token_hash = $1
			AND used_at IS NULL
			AND expires_at > NOW();
	`
	args := []any{tokenHash}

	var token model.PasswordResetToken
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(
		&token.Id,
		&token.UserId,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.UsedAt,
		&token.CreatedAt,
	); err != nil {
		return model.PasswordResetToken{}, err
	}

	return token, nil
}

func (ar *AuthRepository) UpdatePasswordAndUseResetToken(ctx context.Context, resetToken model.PasswordResetToken, hashedPassword string) error {
	tx, err := ar.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	userSQL := `
		UPDATE users
		SET
			password = $2,
			updated_at = NOW()
		WHERE id = $1;
	`
	userCmd, err := tx.Exec(ctx, userSQL, resetToken.UserId, hashedPassword)
	if err != nil {
		return err
	}
	if userCmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	tokenSQL := `
		UPDATE password_reset_tokens
		SET used_at = NOW()
		WHERE id = $1
			AND used_at IS NULL
			AND expires_at > NOW();
	`
	tokenCmd, err := tx.Exec(ctx, tokenSQL, resetToken.Id)
	if err != nil {
		return err
	}
	if tokenCmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	activeTokenSQL := `
		DELETE FROM active_tokens
		WHERE user_id = $1;
	`
	if _, err := tx.Exec(ctx, activeTokenSQL, resetToken.UserId); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
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
