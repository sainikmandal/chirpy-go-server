package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)

	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("password doesn't match: %w", err)

	}
	return nil
}
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// Convert string secret to ECDSA private key
	privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateJWT verifies a JWT and extracts the user ID if valid.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is ECDSA
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Parse private key to get public key
		privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(tokenSecret))
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		return &privateKey.PublicKey, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token claims")
	}

	// Parse Subject (user ID) as UUID
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header provided")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}
	return parts[1], nil
}

// MakeRefreshToken generates a random 256-bit hex-encoded string to use as a refresh token.
func MakeRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
