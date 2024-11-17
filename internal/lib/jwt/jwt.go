package jwt

import (
	"errors"
	"log"
	"os"
	"sso/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	secret := fetchJWTSecret()
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateSigningKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unknown signing method")
	}
	secret := fetchJWTSecret()
	return []byte(secret), nil
}

func ParseToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, validateSigningKey)
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		res := claims["uid"].(float64)
		return int64(res), nil
	} else {
		return 0, err
	}

}

func fetchJWTSecret() string {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("couldn't load env variables %v", err)
	}
	res := os.Getenv("JWT_SECRET")
	return res
}
