package token

import (
	"auth-service/config"
	pb "auth-service/generated/auth_service"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func GenerateAccessJWT(user *pb.LoginResponse) (string, error) {
	cfg := config.Load()
	claims := &Claims{
		UserId:   user.UserId,
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(20 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(cfg.ACCESS_TOKEN))
}

func GenerateRefreshJWT(user *pb.LoginResponse) (string, error) {
	cfg := config.Load()
	claims := &Claims{
		UserId:   user.UserId,
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(14 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(cfg.REFRESH_TOKEN))
}

func ExtractClaim(tokenString string) (*Claims, error) {
	cfg := config.Load()
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.REFRESH_TOKEN), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
