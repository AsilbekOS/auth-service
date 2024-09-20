package auth

import (
	"auth-service/internal/db"
	"auth-service/pkg/config"
	"auth-service/pkg/email"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenExpiry  = 15 * time.Minute
	refreshTokenExpiry = 7 * 24 * time.Hour
)

type Service struct {
	DB      *db.Database
	Config  *config.Config
	Emailer *email.Emailer
}

func NewService(db *db.Database, config *config.Config, emailer *email.Emailer) *Service {
	return &Service{
		DB:      db,
		Config:  config,
		Emailer: emailer,
	}
}

func (s *Service) GenerateTokens(userID string) (string, string, error) {
	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	if err := s.DB.SaveRefreshToken(userID, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(accessTokenExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(s.Config.JWTSecret))
}

func (s *Service) generateRefreshToken(userID string) (string, error) {
	refreshToken := base64.StdEncoding.EncodeToString([]byte(userID + time.Now().String()))
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *Service) RefreshToken(refreshToken string) (string, string, error) {
	userID, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	if err := s.DB.SaveRefreshToken(userID, newRefreshToken); err != nil {
		return "", "", err
	}

	if err := s.DB.DeleteRefreshToken(userID, refreshToken); err != nil {
		return "", "", err
	}

	if err := s.checkIPChange(userID); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *Service) validateRefreshToken(refreshToken string) (string, error) {
	var userID string
	err := s.DB.GetUserIDByRefreshToken(refreshToken, &userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (s *Service) checkIPChange(userID string) error {
	ipAddress, err := s.DB.GetUserIPAddress(userID)
	if err != nil {
		return err
	}

	if ipAddress != "" && ipAddress != s.Config.CurrentIPAddress {
		err := s.Emailer.SendWarningEmail(userID, ipAddress)
		if err != nil {
			return err
		}
	}

	return nil
}
