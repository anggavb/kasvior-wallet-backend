package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/apperrors"
	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/repository"
	"github.com/kasvior-wallet-backend/pkg"
)

type UserService struct {
	db                    *pgxpool.Pool
	userRepository        *repository.UserRepository
	transactionRepository *repository.TransactionRepository
}

func NewUserService(userRepository *repository.UserRepository, transactionRepository *repository.TransactionRepository, db *pgxpool.Pool) *UserService {
	return &UserService{
		db:                    db,
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
	}
}

func (us *UserService) GetProfile(ctx context.Context, userId int) (dto.UserProfileResponse, error) {
	profile, err := us.userRepository.GetProfileById(ctx, userId)
	if err != nil {
		return dto.UserProfileResponse{}, err
	}

	return dto.UserProfileResponse{
		Fullname: profile.Fullname,
		Email:    profile.Email,
		Photo:    profile.Photo,
	}, nil
}

func (us *UserService) UpdateProfile(ctx context.Context, userId int, fullname, phoneNumber *string, photo string) (dto.UserUpdateProfileResponse, error) {
	user, err := us.userRepository.UpdateProfileById(ctx, userId, fullname, phoneNumber, photo)
	if err != nil {
		return dto.UserUpdateProfileResponse{}, err
	}

	return dto.UserUpdateProfileResponse{
		Fullname:    user.Fullname,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Photo:       user.Photo,
	}, nil
}

func (us *UserService) UpdatePassword(ctx context.Context, userId int, req dto.UserUpdatePasswordRequest) error {
	user, err := us.userRepository.GetPasswordById(ctx, userId)
	if err != nil {
		return err
	}

	var hash pkg.HashConfig
	if err := hash.Compare(req.CurrentPassword, user.Password); err != nil {
		return apperrors.ErrInvalidPassword
	}

	hash.UseRecommended()
	hashedPassword := hash.GenerateHash(req.NewPassword)

	return us.userRepository.UpdatePasswordById(ctx, userId, hashedPassword)
}

func (us *UserService) UpdatePin(ctx context.Context, userId int, req dto.UserUpdatePinRequest) error {
	return us.userRepository.UpdatePinById(ctx, userId, req.Pin)
}

func (us *UserService) CheckPin(ctx context.Context, userId int, req dto.UserCheckPinRequest) (dto.UserCheckPinResponse, error) {
	user, err := us.userRepository.GetPinById(ctx, userId)
	if err != nil {
		return dto.UserCheckPinResponse{}, err
	}

	if req.TransactionId != nil {
		return dto.UserCheckPinResponse{}, us.confirmTransfer(ctx, userId, user.Pin, req)
	}

	if user.Pin == nil {
		return dto.UserCheckPinResponse{}, apperrors.ErrPinNotSet
	}

	storedPin := strings.TrimSpace(*user.Pin)
	if req.Pin != storedPin {
		return dto.UserCheckPinResponse{}, apperrors.ErrInvalidPin
	}

	return dto.UserCheckPinResponse{IsValid: true}, nil
}

func (us *UserService) confirmTransfer(ctx context.Context, userId int, storedPin *string, req dto.UserCheckPinRequest) error {
	tx, err := us.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	transfer, err := us.transactionRepository.GetPendingTransferForUpdate(ctx, tx, userId, *req.TransactionId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrTransactionNotFound
		}
		return err
	}
	if transfer.Status != "pending" {
		return apperrors.ErrTransactionFinalized
	}

	if storedPin == nil {
		if err := us.transactionRepository.UpdateTransactionStatus(ctx, tx, transfer.Id, "failed"); err != nil {
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
		return apperrors.ErrPinNotSet
	}

	if req.Pin != strings.TrimSpace(*storedPin) {
		if err := us.transactionRepository.UpdateTransactionStatus(ctx, tx, transfer.Id, "failed"); err != nil {
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
		return apperrors.ErrInvalidPin
	}

	if err := us.transactionRepository.TransferWalletBalance(ctx, tx, transfer.SenderWalletId, transfer.RecipientWalletId, transfer.Amount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if updateErr := us.transactionRepository.UpdateTransactionStatus(ctx, tx, transfer.Id, "failed"); updateErr != nil {
				return updateErr
			}
			if commitErr := tx.Commit(ctx); commitErr != nil {
				return commitErr
			}
			return apperrors.ErrInsufficientBalance
		}
		return err
	}

	if err := us.transactionRepository.UpdateTransactionStatus(ctx, tx, transfer.Id, "success"); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (us *UserService) GetDashboardInformation(ctx context.Context, userId int) (dto.UserDashboardInformationResponse, error) {
	dashboard, err := us.userRepository.GetDashboardInformationById(ctx, userId)
	if err != nil {
		return dto.UserDashboardInformationResponse{}, err
	}

	return dto.UserDashboardInformationResponse{
		Balance: dashboard.Balance,
		Income:  dashboard.Income,
		Expense: dashboard.Expense,
	}, nil
}

func (us *UserService) GetTransactionReport(ctx context.Context, userId int, reportType string) ([]dto.UserTransactionReportResponse, error) {
	endDate := truncateDate(time.Now())
	startDate := endDate.AddDate(0, 0, -6)

	reports, err := us.userRepository.GetTransactionReportById(ctx, userId, reportType, startDate, endDate)
	if err != nil {
		return nil, err
	}

	reportMap := make(map[string]dto.UserTransactionReportResponse, len(reports))
	for _, report := range reports {
		reportMap[report.Date.Format(time.DateOnly)] = dto.UserTransactionReportResponse{
			Day:     report.Date.Format("Mon"),
			Income:  report.Income,
			Expense: report.Expense,
		}
	}

	res := make([]dto.UserTransactionReportResponse, 0, 7)
	for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		report, ok := reportMap[date.Format(time.DateOnly)]
		if !ok {
			report = dto.UserTransactionReportResponse{
				Day:     date.Format("Mon"),
				Income:  0,
				Expense: 0,
			}
		}

		res = append(res, report)
	}

	return res, nil
}

func truncateDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}
