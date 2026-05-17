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

func (ur *UserRepository) GetDashboardInformationById(ctx context.Context, userId int) (model.UserDashboardInformation, error) {
	sqlQuery := `
		SELECT
			w.balance AS balance,
			SUM(
				CASE
					WHEN t.status = 'success' AND t.type IN ('topup', 'receiver')
					THEN t.amount
					ELSE 0
				END
			) AS income,
			SUM(
				CASE
					WHEN t.status = 'success' AND t.type = 'transfer'
					THEN t.amount
					ELSE 0
				END
			) AS expense
		FROM wallets w
		LEFT JOIN transactions t ON t.user_id = w.user_id
		WHERE w.user_id = $1
		GROUP BY w.balance;
	`
	args := []any{userId}

	var dashboard model.UserDashboardInformation
	if err := ur.db.QueryRow(ctx, sqlQuery, args...).Scan(&dashboard.Balance, &dashboard.Income, &dashboard.Expense); err != nil {
		return model.UserDashboardInformation{}, err
	}

	return dashboard, nil
}
