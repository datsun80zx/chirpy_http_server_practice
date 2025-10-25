package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	key := []byte(tokenSecret)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			Subject:   userID.String(),
		})

	s, err := t.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("there was an error in generating jwt: %v", err)
	}
	return s, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	keyfunc := func(token *jwt.Token) (interface{}, error) {
		if token.Method == jwt.SigningMethodHS256 {
			return []byte(tokenSecret), nil
		} else {
			return nil, fmt.Errorf("invalid")
		}
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, keyfunc)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid: %v", err)
	}
	s, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid: %v", err)
	}
	return uuid.Parse(s)
}
