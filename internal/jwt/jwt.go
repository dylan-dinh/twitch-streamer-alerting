package jwt

import (
	"github.com/dylan-dinh/twitch-streamer-alerting/config"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"os"
	"time"
)

const jwtSigningKey = "JWT_SIGNING_KEY"

type JwtService interface {
	GenerateToken() (string, error)
}

type Jwt struct {
	logger *slog.Logger
	Conf   config.Config
}

func NewJwt(conf config.Config) Jwt {
	return Jwt{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Conf:   conf,
	}
}

// GenerateToken will generate a jwt token with basic user info
func (j Jwt) GenerateToken(ID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "tsa-backend",
		Subject:   "login+register",
		Audience:  nil,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 2190)},
		NotBefore: &jwt.NumericDate{Time: time.Now()},
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		ID:        ID,
	}

	keyBytes := []byte(j.Conf.JwtKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(keyBytes)
	if err != nil {
		return "", &GenerateError{Err: err}
	}
	return signedString, nil
}

// GenerateError custom error type
type GenerateError struct {
	Err error
}

func (e *GenerateError) Error() string {
	return e.Err.Error()
}
