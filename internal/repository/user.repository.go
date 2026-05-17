package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetProfileById(ctx context.Context, userId int) (model.User, error) {
	sqlQuery := `
		SELECT fullname, email, photo
		FROM users
		WHERE id = $1;
	`
	args := []any{userId}

	var user model.User
	if err := ur.db.QueryRow(ctx, sqlQuery, args...).Scan(&user.Fullname, &user.Email, &user.Photo); err != nil {
		return model.User{}, err
	}

	return user, nil
}
