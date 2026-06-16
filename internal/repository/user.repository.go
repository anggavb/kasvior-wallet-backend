package repository

import (
	"context"
	"fmt"
	"strings"
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

func (ur *UserRepository) UpdateProfileById(ctx context.Context, userId int, fullname, phoneNumber *string, photo string) (model.User, error) {
	var sql strings.Builder
	args := []any{userId}
	counter := 2

	sql.WriteString("UPDATE users SET")
	if fullname != nil {
		fmt.Fprintf(&sql, " fullname = $%d,", counter)
		args = append(args, fullname)
		counter++
	}
	if phoneNumber != nil {
		fmt.Fprintf(&sql, " phone_number = $%d,", counter)
		args = append(args, phoneNumber)
		counter++
	}
	if photo != "" {
		fmt.Fprintf(&sql, " photo = $%d,", counter)
		args = append(args, photo)
		counter++
	}
	sql.WriteString(" updated_at = NOW()")
	sql.WriteString(" WHERE id = $1 RETURNING fullname, email, phone_number, photo;")

	var user model.User
	if err := ur.db.QueryRow(ctx, sql.String(), args...).Scan(&user.Fullname, &user.Email, &user.PhoneNumber, &user.Photo); err != nil {
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

func (ur *UserRepository) GetBalanceById(ctx context.Context, userId int) (float64, error) {
	sqlQuery := `
		SELECT balance
		FROM wallets
		WHERE user_id = $1;
	`

	var balance float64
	if err := ur.db.QueryRow(ctx, sqlQuery, userId).Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}

func (ur *UserRepository) GetIncomeById(ctx context.Context, userId int) (float64, error) {
	sqlQuery := `
		SELECT COALESCE(SUM(t.amount), 0)
		FROM transactions t
		JOIN transfer_details td ON td.transaction_id = t.id
		JOIN wallets recipient_wallet ON recipient_wallet.id = td.recipient_wallet_id
		WHERE recipient_wallet.user_id = $1
			AND t.type = 'transfer'
			AND t.status = 'success';
	`

	var income float64
	if err := ur.db.QueryRow(ctx, sqlQuery, userId).Scan(&income); err != nil {
		return 0, err
	}

	return income, nil
}

func (ur *UserRepository) GetExpenseById(ctx context.Context, userId int) (float64, error) {
	sqlQuery := `
		SELECT COALESCE(SUM(t.amount), 0)
		FROM transactions t
		JOIN wallets sender_wallet ON sender_wallet.id = t.wallet_id
		WHERE sender_wallet.user_id = $1
			AND t.type = 'transfer'
			AND t.status = 'success';
	`

	var expense float64
	if err := ur.db.QueryRow(ctx, sqlQuery, userId).Scan(&expense); err != nil {
		return 0, err
	}

	return expense, nil
}

func (ur *UserRepository) GetTransactionReportById(ctx context.Context, userId int, reportType string, startDate, endDate time.Time) ([]model.UserTransactionReport, error) {
	sqlQuery := `
		SELECT
			t.created_at AS report_date,
			SUM(
				CASE
					WHEN $2 IN ('all', 'income')
						AND t.status = 'success' AND t.type = 'transfer' AND t.wallet_id != w.id
					THEN t.amount
					ELSE 0
				END
			) AS income,
			SUM(
				CASE
					WHEN $2 IN ('all', 'expense')
						AND t.status = 'success' AND t.type = 'transfer' AND t.wallet_id = w.id
					THEN t.amount
					ELSE 0
				END
			) AS expense
		FROM transactions t
		JOIN wallets w ON w.id = t.wallet_id
		WHERE w.user_id = $1
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
