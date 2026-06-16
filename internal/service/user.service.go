package service

import (
	"context"
	"errors"
	"log"
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
	authCache             *repository.AuthCacheRepository
	dashboardCache        *repository.DashboardCacheRepository
}

func NewUserService(userRepository *repository.UserRepository, transactionRepository *repository.TransactionRepository, authCache *repository.AuthCacheRepository, dashboardCache *repository.DashboardCacheRepository, db *pgxpool.Pool) *UserService {
	return &UserService{
		db:                    db,
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
		authCache:             authCache,
		dashboardCache:        dashboardCache,
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

	if err := us.userRepository.UpdatePasswordById(ctx, userId, hashedPassword); err != nil {
		return err
	}

	return us.authCache.InvalidateUserTokens(ctx, userId)
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

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	us.invalidateSuccessfulTransferDashboard(ctx, userId, transfer.RecipientUserId)
	return nil
}

func (us *UserService) GetBalance(ctx context.Context, userId int) (dto.UserBalanceResponse, error) {
	if balance, ok, err := us.dashboardCache.GetBalance(ctx, userId); err == nil && ok {
		return dto.UserBalanceResponse{Balance: balance}, nil
	} else if err != nil {
		log.Println("Error reading balance cache: ", err.Error())
	}

	balance, err := us.userRepository.GetBalanceById(ctx, userId)
	if err != nil {
		return dto.UserBalanceResponse{}, err
	}

	if err := us.dashboardCache.SetBalance(ctx, userId, balance); err != nil {
		log.Println("Error setting balance cache: ", err.Error())
	}

	return dto.UserBalanceResponse{Balance: balance}, nil
}

func (us *UserService) GetIncome(ctx context.Context, userId int) (dto.UserIncomeResponse, error) {
	if income, ok, err := us.dashboardCache.GetIncome(ctx, userId); err == nil && ok {
		return dto.UserIncomeResponse{Income: income}, nil
	} else if err != nil {
		log.Println("Error reading income cache: ", err.Error())
	}

	income, err := us.userRepository.GetIncomeById(ctx, userId)
	if err != nil {
		return dto.UserIncomeResponse{}, err
	}

	if err := us.dashboardCache.SetIncome(ctx, userId, income); err != nil {
		log.Println("Error setting income cache: ", err.Error())
	}

	return dto.UserIncomeResponse{Income: income}, nil
}

func (us *UserService) GetExpense(ctx context.Context, userId int) (dto.UserExpenseResponse, error) {
	if expense, ok, err := us.dashboardCache.GetExpense(ctx, userId); err == nil && ok {
		return dto.UserExpenseResponse{Expense: expense}, nil
	} else if err != nil {
		log.Println("Error reading expense cache: ", err.Error())
	}

	expense, err := us.userRepository.GetExpenseById(ctx, userId)
	if err != nil {
		return dto.UserExpenseResponse{}, err
	}

	if err := us.dashboardCache.SetExpense(ctx, userId, expense); err != nil {
		log.Println("Error setting expense cache: ", err.Error())
	}

	return dto.UserExpenseResponse{Expense: expense}, nil
}

func (us *UserService) invalidateSuccessfulTransferDashboard(ctx context.Context, senderUserId, recipientUserId int) {
	if err := us.dashboardCache.InvalidateBalance(ctx, senderUserId); err != nil {
		log.Println("Error invalidating sender balance cache: ", err.Error())
	}
	if err := us.dashboardCache.InvalidateExpense(ctx, senderUserId); err != nil {
		log.Println("Error invalidating sender expense cache: ", err.Error())
	}
	if err := us.dashboardCache.InvalidateBalance(ctx, recipientUserId); err != nil {
		log.Println("Error invalidating recipient balance cache: ", err.Error())
	}
	if err := us.dashboardCache.InvalidateIncome(ctx, recipientUserId); err != nil {
		log.Println("Error invalidating recipient income cache: ", err.Error())
	}
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
