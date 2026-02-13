package services

import (
	"errors"
	"time"

	"github.com/Aym-Aymen777/RSS-Aggregator/config"
	"github.com/Aym-Aymen777/RSS-Aggregator/models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TokenService struct {
	config *config.JWTConfig
}

func NewTokenService(config *config.JWTConfig) *TokenService {
	return &TokenService{config: config}
}

// I generate access token
func (s *TokenService) GenerateAccessToken(userID bson.ObjectID, email string) (string, error) {
	// Create claims with user ID and email
	claims := models.AccessTokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// sign the token with the secret
	tokenString, err := token.SignedString([]byte(s.config.AccessSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// II gnerate refresh token
func (s *TokenService) GenerateRefreshToken(userID bson.ObjectID) (string, error) {
	claims := models.RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.RefreshSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// III generate both tokens and return them as a struct
func (s *TokenService) GenerateTokens(userID bson.ObjectID, email string) (*models.TokenResponse, error) {
	accessToken, err := s.GenerateAccessToken(userID, email)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}
	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// IV validate access token
func (s *TokenService) ValidateAccessToken(tokenString string) (*models.AccessTokenClaims, error) {
	// parse the token
	token, err := jwt.ParseWithClaims(tokenString, &models.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(*models.AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// V validate refresh token
func (s *TokenService) ValidateRefreshToken(tokenString string) (*models.RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}
