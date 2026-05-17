package repository

import (
	"context"
	"time"

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

func (ur *UserRepository) GetPinById(ctx context.Context, userId int) (model.User, error) {
	sqlQuery := `
		SELECT pin
		FROM users
		WHERE id = $1;
	`
	args := []any{userId}

	var user model.User
	if err := ur.db.QueryRow(ctx, sqlQuery, args...).Scan(&user.Pin); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) GetPasswordById(ctx context.Context, userId int) (model.User, error) {
	sqlQuery := `
		SELECT password
		FROM users
		WHERE id = $1;
	`
	args := []any{userId}

	var user model.User
	if err := ur.db.QueryRow(ctx, sqlQuery, args...).Scan(&user.Password); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) UpdateProfileById(ctx context.Context, userId int, fullname, phoneNumber, photo *string) (model.User, error) {
	sqlQuery := `
		UPDATE users
		SET
			fullname = COALESCE($2, fullname),
			phone_number = COALESCE($3, phone_number),
			photo = COALESCE($4, photo),
			updated_at = NOW()
		WHERE id = $1
		RETURNING fullname, email, phone_number, photo;
	`
	args := []any{userId, fullname, phoneNumber, photo}

	var user model.User
	if err := ur.db.QueryRow(ctx, sqlQuery, args...).Scan(&user.Fullname, &user.Email, &user.PhoneNumber, &user.Photo); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) UpdatePasswordById(ctx context.Context, userId int, hashedPassword string) error {
	sqlQuery := `
		UPDATE users
		SET
			password = $2,
			updated_at = NOW()
		WHERE id = $1;
	`
	args := []any{userId, hashedPassword}

	_, err := ur.db.Exec(ctx, sqlQuery, args...)
	return err
}

func (ur *UserRepository) UpdatePinById(ctx context.Context, userId int, pin string) error {
	sqlQuery := `
		UPDATE users
		SET
			pin = $2,
			updated_at = NOW()
		WHERE id = $1;
	`
	args := []any{userId, pin}

	_, err := ur.db.Exec(ctx, sqlQuery, args...)
	return err
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

func (ur *UserRepository) GetTransactionReportById(ctx context.Context, userId int, reportType string, startDate, endDate time.Time) ([]model.UserTransactionReport, error) {
	sqlQuery := `
		SELECT
			t.created_at AS report_date,
			SUM(
				CASE
					WHEN $2 IN ('all', 'income')
						AND t.type IN ('topup', 'receiver')
					THEN t.amount
					ELSE 0
				END
			) AS income,
			SUM(
				CASE
					WHEN $2 IN ('all', 'expense')
						AND t.type = 'transfer'
					THEN t.amount
					ELSE 0
				END
			) AS expense
		FROM transactions t
		WHERE t.user_id = $1
			AND t.status = 'success'
			AND t.created_at >= $3
			AND t.created_at < $4
		GROUP BY t.created_at
		ORDER BY t.created_at ASC;
	`
	args := []any{userId, reportType, startDate, endDate.AddDate(0, 0, 1)}

	rows, err := ur.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := []model.UserTransactionReport{}
	for rows.Next() {
		var report model.UserTransactionReport
		if err := rows.Scan(&report.Date, &report.Income, &report.Expense); err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}
