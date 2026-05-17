package service

import (
	"context"

	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/repository"
	"github.com/kasvior-wallet-backend/pkg"
)

type AuthService struct {
	authRepository *repository.AuthRepository
}

func NewAuthService(authRepository *repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepository: authRepository,
	}
}

func (as *AuthService) RegisterUser(ctx context.Context, user dto.AuthRequest) (dto.AuthResponse, error) {
	// hashing password
	var hash pkg.HashConfig
	hash.UseRecommended()

	hashedPassword := hash.GenerateHash(user.Password)
	newUser, err := as.authRepository.AddNewUser(ctx, user.Email, hashedPassword)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Id:        newUser.Id,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
	}, nil
}
