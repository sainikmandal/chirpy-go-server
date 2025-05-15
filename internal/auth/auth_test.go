package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "securepassword123"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hashed == password {
		t.Error("expected hashed password to be different from original")
	}
}

func TestCheckPasswordHash_Success(t *testing.T) {
	password := "securepassword123"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("error hashing password: %v", err)
	}

	err = CheckPasswordHash(hashed, password)
	if err != nil {
		t.Errorf("expected password to match, got error: %v", err)
	}
}

func TestCheckPasswordHash_Failure(t *testing.T) {
	password := "securepassword123"
	wrongPassword := "wrongpassword456"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("error hashing password: %v", err)
	}

	err = CheckPasswordHash(hashed, wrongPassword)
	if err == nil {
		t.Error("expected error for mismatched password, got nil")
	}
}

// jwt tests
func TestMakeAndValidateJWT(t *testing.T) {
	// Generate ECDSA key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	// Convert private key to PEM format
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	assert.NoError(t, err)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Convert public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	userID := uuid.New()

	token, err := MakeJWT(userID, string(privateKeyPEM), time.Hour)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedID, err := ValidateJWT(token, string(publicKeyPEM))
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedID)
}

func TestExpiredJWT(t *testing.T) {
	// Generate ECDSA key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	// Convert private key to PEM format
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	assert.NoError(t, err)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Convert public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	userID := uuid.New()

	// Token that expired 1 hour ago
	token, err := MakeJWT(userID, string(privateKeyPEM), -1*time.Hour)
	assert.NoError(t, err)

	_, err = ValidateJWT(token, string(publicKeyPEM))
	assert.Error(t, err)
}

func TestInvalidSecretJWT(t *testing.T) {
	// Generate two different ECDSA key pairs
	privateKey1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	privateKey2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	// Convert private key to PEM format
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey1)
	assert.NoError(t, err)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Convert wrong public key to PEM format
	wrongPublicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey2.PublicKey)
	assert.NoError(t, err)
	wrongPublicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: wrongPublicKeyBytes,
	})

	userID := uuid.New()

	token, err := MakeJWT(userID, string(privateKeyPEM), time.Hour)
	assert.NoError(t, err)

	_, err = ValidateJWT(token, string(wrongPublicKeyPEM))
	assert.Error(t, err)
}
