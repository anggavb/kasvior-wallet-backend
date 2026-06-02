package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/apperrors"
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
		SELECT u.id, w.id::text AS wallet_id, u.photo, COALESCE(u.fullname, u.email) AS receiver, u.phone_number
		FROM users u
		JOIN wallets w ON w.user_id = u.id
		WHERE u.id != $1
			AND (
				COALESCE(u.fullname, u.email) ILIKE $2 || '%'
				OR u.phone_number ILIKE $2 || '%'
			)
		ORDER BY COALESCE(u.fullname, u.email) ASC
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
		if err := rows.Scan(&receiver.Id, &receiver.WalletId, &receiver.Photo, &receiver.Receiver, &receiver.PhoneNumber); err != nil {
			return nil, err
		}

		receivers = append(receivers, receiver)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return receivers, nil
}

func (tr *TransactionRepository) FindHistory(ctx context.Context, dbtx DBTX, userId int, search string, limit, offset int) ([]model.TransactionHistoryItem, int, error) {
	baseQuery := `
		WITH history AS (
			SELECT
				t.id,
				t.type::text AS transaction_type,
				'in'::text AS direction,
				t.status::text AS status,
				t.amount::float8 AS amount,
				pm.name::text AS counterparty_name,
				NULL::text AS counterparty_phone,
				pm.logo::text AS counterparty_photo,
				pm.name::text AS payment_method,
				NULL::text AS notes,
				t.created_at
			FROM transactions t
			JOIN wallets owner_wallet ON owner_wallet.id = t.wallet_id
			JOIN topup_details td ON td.transaction_id = t.id
			JOIN payment_methods pm ON pm.id = td.payment_method_id
			WHERE owner_wallet.user_id = $1
				AND t.type = 'topup'
				AND t.status IN ('success', 'failed')

			UNION ALL

			SELECT
				t.id,
				t.type::text AS transaction_type,
				'out'::text AS direction,
				t.status::text AS status,
				t.amount::float8 AS amount,
				COALESCE(recipient.fullname, recipient.email)::text AS counterparty_name,
				recipient.phone_number::text AS counterparty_phone,
				recipient.photo::text AS counterparty_photo,
				NULL::text AS payment_method,
				td.notes::text AS notes,
				t.created_at
			FROM transactions t
			JOIN wallets sender_wallet ON sender_wallet.id = t.wallet_id
			JOIN transfer_details td ON td.transaction_id = t.id
			JOIN wallets recipient_wallet ON recipient_wallet.id = td.recipient_wallet_id
			JOIN users recipient ON recipient.id = recipient_wallet.user_id
			WHERE sender_wallet.user_id = $1
				AND t.type = 'transfer'
				AND t.status IN ('success', 'failed')

			UNION ALL

			SELECT
				t.id,
				t.type::text AS transaction_type,
				'in'::text AS direction,
				t.status::text AS status,
				t.amount::float8 AS amount,
				COALESCE(sender.fullname, sender.email)::text AS counterparty_name,
				sender.phone_number::text AS counterparty_phone,
				sender.photo::text AS counterparty_photo,
				NULL::text AS payment_method,
				td.notes::text AS notes,
				t.created_at
			FROM transactions t
			JOIN wallets sender_wallet ON sender_wallet.id = t.wallet_id
			JOIN users sender ON sender.id = sender_wallet.user_id
			JOIN transfer_details td ON td.transaction_id = t.id
			JOIN wallets recipient_wallet ON recipient_wallet.id = td.recipient_wallet_id
			WHERE recipient_wallet.user_id = $1
				AND t.type = 'transfer'
				AND t.status = 'success'
		), filtered AS (
			SELECT *
			FROM history
			WHERE $2 = ''
				OR counterparty_name ILIKE '%' || $2 || '%'
				OR COALESCE(counterparty_phone, '') ILIKE '%' || $2 || '%'
				OR transaction_type ILIKE '%' || $2 || '%'
				OR COALESCE(payment_method, '') ILIKE '%' || $2 || '%'
				OR COALESCE(notes, '') ILIKE '%' || $2 || '%'
		)
	`

	var total int
	countQuery := baseQuery + `
		SELECT COUNT(*)
		FROM filtered;
	`
	if err := dbtx.QueryRow(ctx, countQuery, userId, search).Scan(&total); err != nil {
		return nil, 0, err
	}

	itemQuery := baseQuery + `
		SELECT
			id,
			transaction_type,
			direction,
			status,
			amount,
			counterparty_name,
			counterparty_phone,
			counterparty_photo,
			payment_method,
			notes,
			created_at
		FROM filtered
		ORDER BY created_at DESC, id DESC
		LIMIT $3
		OFFSET $4;
	`

	rows, err := dbtx.Query(ctx, itemQuery, userId, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := []model.TransactionHistoryItem{}
	for rows.Next() {
		var item model.TransactionHistoryItem
		if err := rows.Scan(
			&item.Id,
			&item.Type,
			&item.Direction,
			&item.Status,
			&item.Amount,
			&item.CounterpartyName,
			&item.CounterpartyPhone,
			&item.CounterpartyPhoto,
			&item.PaymentMethod,
			&item.Notes,
			&item.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, total, nil
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

func (tr *TransactionRepository) CreateTransaction(ctx context.Context, dbtx DBTX, userId int, typeTransaction, status string, amount uint) (int, error) {
	sql := `
		INSERT INTO transactions (wallet_id, amount, type, status)
		SELECT id, $2, $3, $4
		FROM wallets
		WHERE user_id = $1
		RETURNING id;
	`
	args := []any{userId, amount, typeTransaction, status}

	var transactionId int
	if err := dbtx.QueryRow(ctx, sql, args...).Scan(&transactionId); err != nil {
		return 0, err
	}

	return transactionId, nil
}

func (tr *TransactionRepository) CreateTopupTransactionDetails(ctx context.Context, dbtx DBTX, transactionId int, topup dto.TopupRequest) error {
	sql := `
		INSERT INTO topup_details (transaction_id, payment_method_id, discount, tax, sub_total)
		VALUES ($1, $2, $3, $4, $5);
	`
	args := []any{transactionId, topup.PaymentMethodId, *topup.Discount, *topup.Tax, *topup.SubTotal}

	_, err := dbtx.Exec(ctx, sql, args...)
	return err
}

func (tr *TransactionRepository) IncrementWalletBalanceByUserId(ctx context.Context, dbtx DBTX, userId int, amount uint) error {
	sql := `
		UPDATE wallets
		SET balance = balance + $2
		WHERE user_id = $1;
	`

	_, err := dbtx.Exec(ctx, sql, userId, amount)
	return err
}

func (tr *TransactionRepository) GetWalletIdByUserId(ctx context.Context, dbtx DBTX, userId int) (string, error) {
	sql := `
		SELECT id::text
		FROM wallets
		WHERE user_id = $1;
	`

	var walletId string
	if err := dbtx.QueryRow(ctx, sql, userId).Scan(&walletId); err != nil {
		return "", err
	}

	return walletId, nil
}

func (tr *TransactionRepository) WalletExists(ctx context.Context, dbtx DBTX, walletId string) (bool, error) {
	sql := `
		SELECT EXISTS (
			SELECT 1
			FROM wallets
			WHERE id = $1
		);
	`

	var exists bool
	if err := dbtx.QueryRow(ctx, sql, walletId).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (tr *TransactionRepository) CreateTransferDetail(ctx context.Context, dbtx DBTX, transactionId int, transfer dto.TransferRequest) error {
	sql := `
		INSERT INTO transfer_details (transaction_id, recipient_wallet_id, notes)
		VALUES ($1, $2, $3);
	`

	_, err := dbtx.Exec(ctx, sql, transactionId, transfer.RecipientWalletId, transfer.Notes)
	return err
}

func (tr *TransactionRepository) GetPendingTransferForUpdate(ctx context.Context, dbtx DBTX, userId, transactionId int) (model.TransferTransaction, error) {
	sql := `
		SELECT
			t.id,
			t.wallet_id::text,
			td.recipient_wallet_id::text,
			t.amount,
			t.status
		FROM transactions t
		JOIN wallets w ON w.id = t.wallet_id
		JOIN transfer_details td ON td.transaction_id = t.id
		WHERE t.id = $1
			AND w.user_id = $2
			AND t.type = 'transfer'
		FOR UPDATE OF t;
	`

	var transfer model.TransferTransaction
	if err := dbtx.QueryRow(ctx, sql, transactionId, userId).Scan(
		&transfer.Id,
		&transfer.SenderWalletId,
		&transfer.RecipientWalletId,
		&transfer.Amount,
		&transfer.Status,
	); err != nil {
		return model.TransferTransaction{}, err
	}

	return transfer, nil
}

func (tr *TransactionRepository) UpdateTransactionStatus(ctx context.Context, dbtx DBTX, transactionId int, status string) error {
	sql := `
		UPDATE transactions
		SET status = $2,
			updated_at = NOW()
		WHERE id = $1;
	`

	_, err := dbtx.Exec(ctx, sql, transactionId, status)
	return err
}

func (tr *TransactionRepository) TransferWalletBalance(ctx context.Context, dbtx DBTX, senderWalletId, recipientWalletId string, amount float64) error {
	debitSQL := `
		UPDATE wallets
		SET balance = balance - $2
		WHERE id = $1
			AND balance >= $2;
	`

	debitTag, err := dbtx.Exec(ctx, debitSQL, senderWalletId, amount)
	if err != nil {
		return err
	}
	if debitTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	creditSQL := `
		UPDATE wallets
		SET balance = balance + $2
		WHERE id = $1;
	`

	creditTag, err := dbtx.Exec(ctx, creditSQL, recipientWalletId, amount)
	if err != nil {
		return err
	}
	if creditTag.RowsAffected() == 0 {
		return apperrors.ErrInvalidRecipient
	}

	return nil
}
