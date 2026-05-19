package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/model"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (tr *TransactionRepository) FindReceivers(ctx context.Context, userId int, search string, limit, offset int) ([]model.Receiver, error) {
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

	rows, err := tr.db.Query(ctx, sqlQuery, args...)
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

func (tr *TransactionRepository) AddTopupTransaction(ctx context.Context, userId, paymentMethodId, amount, discount, tax int, typeTransaction string) {
}
