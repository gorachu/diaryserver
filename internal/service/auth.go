package service

import (
	"diaryserver/internal/storage/sqlite"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	storage       *sqlite.Storage
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewAuthService(storage *sqlite.Storage, accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		storage:       storage,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (s *AuthService) ValidateUser(username, password string) (*sqlite.UserInfo, error) {
	user, err := s.storage.GetUser(username)
	if err != nil {
		return nil, err
	}

	if err := VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) GenerateAccessToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "access",
		"exp":     time.Now().Add(s.accessTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.accessSecret))
}

func (s *AuthService) GenerateRefreshToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(s.refreshTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecret))
}
func (s *AuthService) GenerateTokens(userID int64) (string, string, error) {
	accessToken, err := s.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	if claims["type"] != "refresh" {
		return "", "", errors.New("wrong token type")
	}

	userID := int64(claims["user_id"].(float64))

	newAccessToken, err := s.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) ValidateAccessToken(accessToken string) (int64, error) {
	token2 := accessToken
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v, ожидается HS256", token.Header["alg"])
		}

		if token.Header["alg"] != "HS256" {
			return nil, fmt.Errorf("неподдерживаемый алгоритм: %v, ожидается HS256", token.Header["alg"])
		}
		return []byte(s.accessSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid access token " + token2)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	if claims["type"] != "access" {
		return 0, errors.New("wrong token type")
	}

	userID := int64(claims["user_id"].(float64))
	return userID, nil
}

func (s *AuthService) ValidateRefreshToken(refreshToken string) (int64, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	if claims["type"] != "refresh" {
		return 0, errors.New("wrong token type")
	}

	userID := int64(claims["user_id"].(float64))
	return userID, nil
}
