package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId     int
	Email      string
	IsVerified bool
	jwt.RegisteredClaims
}

func NewClaims(id int, email string, isVerified bool) *Claims {
	return &Claims{
		UserId:     id,
		Email:      email,
		IsVerified: isVerified,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Hour)),
		},
	}
}

func (c *Claims) GenerateJWT() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("missing jwt secret")
	}
	uToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return uToken.SignedString([]byte(jwtSecret))
}

func (c *Claims) VerifyJWT(token string) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("missing jwt secret")
	}

	jwtToken, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return jwt.ErrTokenExpired
	}

	iss, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		return err
	}

	if iss != os.Getenv("JWT_ISSUER") {
		return jwt.ErrTokenInvalidIssuer
	}

	return nil
}
