package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenMalformed   = errors.New("token is malformed")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenHandle      = errors.New("couldn't handle this token")
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

var tokenBlackList = NewRedisBlacklist()

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {
	if username == "" {
		return "", errors.New("userID cannot be empty")
	}

	expireTime := time.Now().Add(24 * time.Hour)
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "user-backend",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.New("failed to sign token: " + err.Error())
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token string is empty")
	}

	isBlacklisted, err := tokenBlackList.IsTokenBlacklisted(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %v", err)
	}
	if isBlacklisted {
		return nil, ErrTokenInvalid
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, parseJWTError(err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

func InvalidateToken(tokenString string) error {
	if tokenString == "" {
		return errors.New("token string is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		expireTime := time.Unix(claims.ExpiresAt, 0)
		err := tokenBlackList.AddToken(tokenString, expireTime)
		if err != nil {
			return fmt.Errorf("failed to invalidate token: %v", err)
		}
		return nil
	}

	return ErrTokenInvalid
}

func parseJWTError(err error) error {
	if ve, ok := err.(*jwt.ValidationError); ok {
		switch {
		case ve.Errors&jwt.ValidationErrorMalformed != 0:
			return ErrTokenMalformed
		case ve.Errors&jwt.ValidationErrorExpired != 0:
			return ErrTokenExpired
		case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
			return ErrTokenNotValidYet
		default:
			return ErrTokenHandle
		}
	}
	return err
}

func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}

func GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}
