package service

import (
	"context"

	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/repository"
)

type TransactionService struct {
	transactionRepository *repository.TransactionRepository
}

func NewTransactionService(transactionRepository *repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
	}
}

func (ts *TransactionService) FindReceivers(ctx context.Context, userId int, search string, page, limit int) (dto.ReceiverListResponse, error) {
	offset := (page - 1) * limit

	receivers, err := ts.transactionRepository.FindReceivers(ctx, userId, search, limit, offset)
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
