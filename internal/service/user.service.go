package service

import (
	"context"
	"time"

	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
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
