package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kasvior-wallet-backend/internal/apperrors"
	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/repository"
	"github.com/kasvior-wallet-backend/pkg"
)

type Mailer interface {
	SendMail(ctx context.Context, to, subject, body string) error
}

const passwordResetTokenTTL = 15 * time.Minute

type AuthService struct {
	authRepository *repository.AuthRepository
	mailer         Mailer
}

func NewAuthService(authRepository *repository.AuthRepository, mailer Mailer) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		mailer:         mailer,
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

	claims := pkg.NewClaims(userLogin.Id, user.Email, userLogin.IsVerified)
	token, err := claims.GenerateJWT()
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if err := as.authRepository.SaveToken(ctx, hashToken(token), userLogin.Id, claims.ExpiresAt.Time); err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		Email: user.Email,
		Token: token,
	}, nil
}

func (as *AuthService) LogoutUser(ctx context.Context, token string, expiresAt *time.Time) error {
	if expiresAt == nil {
		return errors.New("missing token expiry")
	}
	if time.Now().After(*expiresAt) {
		return apperrors.ErrTokenAlreadyExpired
	}

	return as.authRepository.DeleteToken(ctx, hashToken(token))
}

func (as *AuthService) ForgotPassword(ctx context.Context, request dto.ForgotPasswordRequest) error {
	user, err := as.authRepository.GetPasswordResetUserByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	token, err := generateResetToken()
	if err != nil {
		return err
	}

	if err := as.authRepository.SavePasswordResetToken(ctx, user.Id, hashToken(token), time.Now().Add(passwordResetTokenTTL)); err != nil {
		return err
	}

	resetPasswordURL := os.Getenv("RESET_PASSWORD_URL")
	if resetPasswordURL == "" {
		return errors.New("reset password url is required")
	}
	if as.mailer == nil {
		return errors.New("mailer is required")
	}

	resetLink := resetPasswordURL + token
	subject := "Reset Password Kasvior Wallet"
	body := fmt.Sprintf("Use this link to reset your password:\n\n%s\n\nThis link expires in 15 minutes.\n", resetLink)

	return as.mailer.SendMail(ctx, user.Email, subject, body)
}

func (as *AuthService) ResetPassword(ctx context.Context, request dto.ResetPasswordRequest) error {
	resetToken, err := as.authRepository.GetActivePasswordResetToken(ctx, hashToken(request.Token))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrInvalidPasswordResetToken
		}
		return err
	}

	var hash pkg.HashConfig
	hash.UseRecommended()

	hashedPassword := hash.GenerateHash(request.NewPassword)
	if err := as.authRepository.UpdatePasswordAndUseResetToken(ctx, resetToken, hashedPassword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrInvalidPasswordResetToken
		}
		return err
	}

	return nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func generateResetToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}
