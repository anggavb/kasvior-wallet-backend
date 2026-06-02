package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
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
			WalletId:    receiver.WalletId,
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

func (ts *TransactionService) FindHistory(ctx context.Context, userId int, search string, page, limit int) (dto.TransactionHistoryResponse, error) {
	offset := (page - 1) * limit

	history, total, err := ts.transactionRepository.FindHistory(ctx, ts.db, userId, search, limit, offset)
	if err != nil {
		return dto.TransactionHistoryResponse{}, err
	}

	items := make([]dto.TransactionHistoryItemResponse, 0, len(history))
	for _, item := range history {
		items = append(items, dto.TransactionHistoryItemResponse{
			Id:                item.Id,
			Type:              item.Type,
			Direction:         item.Direction,
			Status:            item.Status,
			Amount:            item.Amount,
			CounterpartyName:  item.CounterpartyName,
			CounterpartyPhone: item.CounterpartyPhone,
			CounterpartyPhoto: item.CounterpartyPhoto,
			PaymentMethod:     item.PaymentMethod,
			Notes:             item.Notes,
			CreatedAt:         item.CreatedAt,
		})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return dto.TransactionHistoryResponse{
		Items: items,
		Meta: dto.TransactionHistoryMetaResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
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

func (ts *TransactionService) CreateTransactionWithDetails(ctx context.Context, userId int, topup dto.TopupRequest) error {
	isSubtotalValid := *topup.SubTotal == (int(topup.Amount) - *topup.Discount + *topup.Tax)
	if !isSubtotalValid {
		return apperrors.InvalidSubtotal
	}

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tid, err := ts.transactionRepository.CreateTransaction(ctx, tx, userId, topup.TypeTransaction, "success", topup.Amount)
	if err != nil {
		return err
	}

	if err := ts.transactionRepository.CreateTopupTransactionDetails(ctx, tx, tid, topup); err != nil {
		return err
	}

	if err := ts.transactionRepository.IncrementWalletBalanceByUserId(ctx, tx, userId, topup.Amount); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (ts *TransactionService) CreatePendingTransfer(ctx context.Context, userId int, transfer dto.TransferRequest) (dto.TransactionCreatedResponse, error) {
	transfer.RecipientWalletId = strings.ToLower(transfer.RecipientWalletId)

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return dto.TransactionCreatedResponse{}, err
	}
	defer tx.Rollback(ctx)

	senderWalletId, err := ts.transactionRepository.GetWalletIdByUserId(ctx, tx, userId)
	if err != nil {
		return dto.TransactionCreatedResponse{}, err
	}

	if senderWalletId == transfer.RecipientWalletId {
		return dto.TransactionCreatedResponse{}, apperrors.ErrSelfTransfer
	}

	exists, err := ts.transactionRepository.WalletExists(ctx, tx, transfer.RecipientWalletId)
	if err != nil {
		return dto.TransactionCreatedResponse{}, err
	}
	if !exists {
		return dto.TransactionCreatedResponse{}, apperrors.ErrInvalidRecipient
	}

	transactionId, err := ts.transactionRepository.CreateTransaction(ctx, tx, userId, "transfer", "pending", transfer.Amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.TransactionCreatedResponse{}, apperrors.ErrInvalidRecipient
		}
		return dto.TransactionCreatedResponse{}, err
	}

	if err := ts.transactionRepository.CreateTransferDetail(ctx, tx, transactionId, transfer); err != nil {
		return dto.TransactionCreatedResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.TransactionCreatedResponse{}, err
	}

	return dto.TransactionCreatedResponse{TransactionId: transactionId}, nil
}
