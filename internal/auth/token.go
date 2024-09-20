package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type TokenPayload struct {
	UserID string `json:"userID"`
	IP     string `json:"ip"`
	Exp    int64  `json:"exp"`
}

func (p TokenPayload) Valid() error {
	if time.Unix(p.Exp, 0).Before(time.Now()) {
		return errors.New("token expired")
	}
	return nil
}

func CreateJWTToken(userID, ip, secret string, duration time.Duration) (string, error) {
	exp := time.Now().Add(duration).Unix()

	claims := TokenPayload{
		UserID: userID,
		IP:     ip,
		Exp:    exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func VerifyJWTToken(tokenString, secret string) (*TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func HashRefreshToken(token string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedToken), nil
}

func VerifyRefreshToken(hashedToken, token string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
	if err != nil {
		return false, err
	}
	return true, nil
}
