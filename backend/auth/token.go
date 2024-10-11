package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetBearrerOfTheToken(header http.Header) (string, error) {
	/*
		Function that extracts the bearer token from the Authorization header.
	*/
	authorizationHeader := header.Get("Authorization")

	if authorizationHeader == "" {
		return "", errors.New("no auth header provided in the request")
	}

	dataFromHeader := strings.Split(authorizationHeader, " ")
	if len(dataFromHeader) < 2 || dataFromHeader[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return dataFromHeader[1], nil
}

func MakeRefreshToken() (string, error) {
	/*
		Function that generates a random 32-byte slice and returns it as a hex-encoded string.
		This is used to generate a refresh token.
	*/
	randomSlice := make([]byte, 32)

	_, err := rand.Read(randomSlice)
	if err != nil {
		return "", err
	}

	encodedSlice := hex.EncodeToString(randomSlice)

	return encodedSlice, nil
}

func MakeAccessToken(userId int, secret string, expiresIn time.Duration) (string, error) {
	/*
		Function that generates an access token.
		It accepts the user's ID, the secret key and the expiration time.
		It returns the access token as a string.
	*/
	myClaims := jwt.RegisteredClaims{
		Issuer:    "Streamify",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   strconv.Itoa(userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ValidateAccessToken(tokenString, tokenSecret string) (int, error) {
	/*
		Function that validates the access token.
		It accepts the token string and the secret key that must be defined in .env.
	*/
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, jwt.ErrSignatureInvalid
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, jwt.ErrInvalidKeyType
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
