package repository

import (
	"context"

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
