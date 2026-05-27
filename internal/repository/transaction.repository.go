package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/model"
)

type DBTX interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type TransactionRepository struct{}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{}
}

func (tr *TransactionRepository) FindReceivers(ctx context.Context, dbtx DBTX, userId int, search string, limit, offset int) ([]model.Receiver, error) {
	sqlQuery := `
		SELECT id, photo, COALESCE(fullname, email) AS receiver, phone_number
		FROM users
		WHERE id != $1
			AND (
				COALESCE(fullname, email) ILIKE $2 || '%'
				OR phone_number ILIKE $2 || '%'
			)
		ORDER BY COALESCE(fullname, email) ASC
		LIMIT $3
		OFFSET $4;
	`
	args := []any{userId, search, limit, offset}

	rows, err := dbtx.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receivers := []model.Receiver{}
	for rows.Next() {
		var receiver model.Receiver
		if err := rows.Scan(&receiver.Id, &receiver.Photo, &receiver.Receiver, &receiver.PhoneNumber); err != nil {
			return nil, err
		}

		receivers = append(receivers, receiver)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return receivers, nil
}

func (tr *TransactionRepository) GetPaymentMethodById(ctx context.Context, dbtx DBTX, id int) (model.PaymentMethod, error) {
	sql := `
		SELECT id, name, logo, method, tax, created_at, updated_at
		FROM payment_methods
		WHERE id = $1;
	`

	args := []any{id}

	var paymentMethod model.PaymentMethod
	if err := dbtx.QueryRow(ctx, sql, args...).Scan(&paymentMethod.Id, &paymentMethod.Name, &paymentMethod.Logo, &paymentMethod.Method, &paymentMethod.Tax, &paymentMethod.CreatedAt, &paymentMethod.UpdatedAt); err != nil {
		return model.PaymentMethod{}, err
	}

	return paymentMethod, nil
}

func (tr *TransactionRepository) CreateTransaction(ctx context.Context, dbtx DBTX, userId int, typeTransaction string, amount uint) (int, error) {
	sql := `
		INSERT INTO transactions (wallet_id, amount, type, status)
		SELECT id, $2, $3, 'success'
		FROM wallets
		WHERE user_id = $1
		RETURNING id;
	`
	args := []any{userId, amount, typeTransaction}

	var transactionId int
	if err := dbtx.QueryRow(ctx, sql, args...).Scan(&transactionId); err != nil {
		return 0, err
	}
	log.Println(transactionId)

	return transactionId, nil
}

func (tr *TransactionRepository) CreateTopupTransactionDetails(ctx context.Context, dbtx DBTX, transactionId int, topup dto.TopupRequest) (string, error) {
	sql := `
		WITH topup AS (
			INSERT INTO topup_details (transaction_id, payment_method_id, discount, tax, sub_total)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING payment_method_id
		)
		SELECT name FROM payment_methods WHERE id = (SELECT payment_method_id FROM topup);
	`
	args := []any{transactionId, topup.PaymentMethodId, *topup.Discount, *topup.Tax, *topup.SubTotal}

	var paymentMethod string
	if err := dbtx.QueryRow(ctx, sql, args...).Scan(&paymentMethod); err != nil {
		return "", err
	}

	return paymentMethod, nil
}
