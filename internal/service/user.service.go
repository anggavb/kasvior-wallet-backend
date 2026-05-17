package service

import (
	"context"

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
