package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/apperrors"
	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/repository"
)

type TransactionService struct {
	db                    *pgxpool.Pool
	transactionRepository *repository.TransactionRepository
}

func NewTransactionService(transactionRepository *repository.TransactionRepository, db *pgxpool.Pool) *TransactionService {
	return &TransactionService{
		db:                    db,
		transactionRepository: transactionRepository,
	}
}

func (ts *TransactionService) FindReceivers(ctx context.Context, userId int, search string, page, limit int) (dto.ReceiverListResponse, error) {
	offset := (page - 1) * limit

	receivers, err := ts.transactionRepository.FindReceivers(ctx, ts.db, userId, search, limit, offset)
	if err != nil {
		return dto.ReceiverListResponse{}, err
	}

	items := make([]dto.ReceiverResponse, 0, len(receivers))
	for _, receiver := range receivers {
		items = append(items, dto.ReceiverResponse{
			Id:          receiver.Id,
			Photo:       receiver.Photo,
			Receiver:    receiver.Receiver,
			PhoneNumber: receiver.PhoneNumber,
		})
	}

	return dto.ReceiverListResponse{
		Items: items,
		Meta: dto.PaginationMetaResponse{
			Page:  page,
			Limit: limit,
		},
	}, nil
}

func (ts *TransactionService) GetPaymentMethodById(ctx context.Context, paymentMethodId int) (dto.PaymentMethodResponse, error) {
	paymentMethod, err := ts.transactionRepository.GetPaymentMethodById(ctx, ts.db, paymentMethodId)
	if err != nil {
		return dto.PaymentMethodResponse{}, err
	}

	return dto.PaymentMethodResponse{
		Id:     paymentMethod.Id,
		Name:   paymentMethod.Name,
		Logo:   paymentMethod.Logo,
		Method: paymentMethod.Method,
		Tax:    paymentMethod.Tax,
	}, nil
}

func (ts *TransactionService) CreateTransactionWithDetails(ctx context.Context, userId int, topup dto.TopupRequest) (string, error) {
	isSubtotalValid := topup.SubTotal == (int(topup.Amount) - topup.Discount + topup.Tax)
	if !isSubtotalValid {
		return "", apperrors.InvalidSubtotal
	}

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	tid, err := ts.transactionRepository.CreateTransaction(ctx, tx, userId, topup.TypeTransaction, topup.Amount)
	if err != nil {
		return "", err
	}

	paymentMethod, err := ts.transactionRepository.CreateTopupTransactionDetails(ctx, tx, tid, topup)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return paymentMethod, err
}
