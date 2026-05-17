package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId   int
	Username string
	jwt.RegisteredClaims
}

func NewClaims(id int, username string) *Claims {
	return &Claims{
		UserId:   id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
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
