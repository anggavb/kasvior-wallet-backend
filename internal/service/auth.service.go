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

func (as *AuthService) LoginUser(ctx context.Context, user dto.AuthRequest) (dto.AuthResponse, error) {
	userLogin, err := as.authRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	var hash pkg.HashConfig
	if err := hash.Compare(user.Password, userLogin.Password); err != nil {
		return dto.AuthResponse{}, err
	}

	token, err := pkg.NewClaims(userLogin.Id, user.Email, userLogin.IsVerified).GenerateJWT()
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Email: user.Email,
		Token: token,
	}, nil
}
